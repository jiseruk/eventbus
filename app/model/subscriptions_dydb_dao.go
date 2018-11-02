package model

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var subscribersTable = "Subscribers"

type SubscriberDaoDynamoImpl struct {
	DynamoClient dynamodbiface.DynamoDBAPI
}

func (s *SubscriberDaoDynamoImpl) CreateSubscription(name string, topic string, endpoint string,
	resource string, pullResource string) (*Subscriber, error) {
	subscription := Subscriber{Name: name, Topic: topic, Endpoint: endpoint, ResourceID: resource, PullResourceID: pullResource}
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
	fmt.Printf("GET ITEM: %v", output)
	//value := output.Item
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