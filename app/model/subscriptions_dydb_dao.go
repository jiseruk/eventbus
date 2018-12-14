package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/wenance/wequeue-management_api/app/config"
)

var subscribersTable = config.Get("databases.dynamodb.tables.subscribers")

type SubscriberDaoDynamoImpl struct {
	DynamoClient dynamodbiface.DynamoDBAPI
}

func (s *SubscriberDaoDynamoImpl) CreateSubscription(name string, topic string, Type string, resource string,
	endpoint *string, deadLetterQueue string, pullingQueue string, visibilityTimeout *int) (*Subscriber, error) {

	subscription := Subscriber{Name: name, Topic: topic, Endpoint: endpoint,
		ResourceID: resource, DeadLetterQueue: deadLetterQueue,
		PullingQueue: pullingQueue, Type: Type,
		VisibilityTimeout: visibilityTimeout,
	}
	subscription.CreatedAt = Clock.Now()
	subscription.UpdatedAt = Clock.Now()

	item, err := dynamodbattribute.MarshalMap(subscription)
	if err != nil {
		return nil, err
	}
	//TODO: Manage errors
	_, err = s.DynamoClient.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: &subscribersTable,
	})
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (s *SubscriberDaoDynamoImpl) GetSubscription(name string) (*Subscriber, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: &subscribersTable,
	}
	output, err := s.DynamoClient.GetItem(input)
	if err != nil || output.Item == nil {
		return nil, err
	}

	var subscriber Subscriber
	err = dynamodbattribute.UnmarshalMap(output.Item, &subscriber)
	if err != nil {
		return nil, err
	}
	return &subscriber, nil
}

func (s *SubscriberDaoDynamoImpl) GetSubscriptionByEndpoint(endpoint string) (*Subscriber, error) {
	return nil, nil
}

func (s *SubscriberDaoDynamoImpl) GetSubscriptionsByTopic(topic string) ([]Subscriber, error) {
	input := &dynamodb.ScanInput{
		ProjectionExpression:     aws.String("#name, #type, endpoint, visibility_timeout, created_at"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name"), "#type": aws.String("type")},
		FilterExpression:         aws.String("topic = :t"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: &topic,
			},
		},
		TableName: &subscribersTable,
	}
	output, err := s.DynamoClient.Scan(input)
	var subscribers []Subscriber
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &subscribers)
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (s *SubscriberDaoDynamoImpl) DeleteTopicSubscriptions(topic string) error {
	subscribers, _ := s.GetSubscriptionsByTopic(topic)
	if len(subscribers) == 0 {
		return nil
	}

	ids := make([]*dynamodb.WriteRequest, len(subscribers))

	for i, s := range subscribers {
		ids[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"name": {
						S: aws.String(s.Name),
					},
				},
			},
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			subscribersTable: ids,
		},
	}
	_, err := s.DynamoClient.BatchWriteItem(input)
	return err
}

func (s *SubscriberDaoDynamoImpl) DeleteSubscription(name string) error {
	_, err := s.DynamoClient.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: &subscribersTable,
	})
	if err != nil {
		if apierr, ok := err.(awserr.Error); ok {
			if apierr.Code() != dynamodb.ErrCodeResourceNotFoundException {
				return err
			}
		}

	}
	return nil
}
