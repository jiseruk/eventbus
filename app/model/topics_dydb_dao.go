package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/wenance/wequeue-management_api/app/config"
)

var topicsTable = config.Get("databases.dynamodb.tables.topics")

type TopicsDaoDynamoImpl struct {
	DynamoClient dynamodbiface.DynamoDBAPI
	UUID         UUID
}

func (t *TopicsDaoDynamoImpl) CreateTopic(name string, engine string, owner string, description string, resourceID string) (*Topic, error) {
	topic := Topic{Name: name, Engine: engine, ResourceID: resourceID, Owner: owner, Description: description}
	topic.CreatedAt = Clock.Now()
	topic.UpdatedAt = Clock.Now()
	topic.SecurityToken = t.UUID.GetUUID()
	//topic.ID = uuid.New()
	//topic.ID = 1
	item, err := dynamodbattribute.MarshalMap(topic)
	if err != nil {
		return nil, err
	}
	//TODO: Manage errors
	_, err = t.DynamoClient.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: &topicsTable,
	})
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (t *TopicsDaoDynamoImpl) GetTopic(name string) (*Topic, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: &topicsTable,
	}

	output, err := t.DynamoClient.GetItem(input)
	if err != nil || output.Item == nil {
		return nil, err
	}
	var topic Topic
	err = dynamodbattribute.UnmarshalMap(output.Item, &topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (t *TopicsDaoDynamoImpl) DeleteTopic(name string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: &topicsTable,
	}

	_, err := t.DynamoClient.DeleteItem(input)
	return err

}

func (t *TopicsDaoDynamoImpl) ListTopics() ([]Topic, error) {
	input := &dynamodb.ScanInput{
		TableName:                &topicsTable,
		ProjectionExpression:     aws.String("#name, engine, #owner, description, created_at"),
		ExpressionAttributeNames: map[string]*string{"#name": aws.String("name"), "#owner": aws.String("owner")},
	}

	output, err := t.DynamoClient.Scan(input)
	if err != nil {
		return nil, err
	}

	var topics []Topic
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &topics)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
