package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/wenance/wequeue-management_api/app/config"
)

/*
type DynamoDB struct {
	//Mock with https://github.com/Selvatico/go-mocket/wiki/Documentation
	dynamoClient dynamodbiface.DynamoDBAPI
}
*/
var dynamoEndpoint = config.Get("engines.AWS.dynamodb.endpoint")

func GetClient() dynamodbiface.DynamoDBAPI {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: config.GetObject("aws_credentials").(*credentials.Credentials),
	})
	if err != nil {
		panic("FATAL: Connot connect to AWS")
	}
	dynamoClient := dynamodb.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint(dynamoEndpoint))
	if *config.GetCurrentEnvironment() == config.LOCAL {
		dynamoClient.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String("Subscribers"),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("name"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("name"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
	}
	return dynamoClient
}