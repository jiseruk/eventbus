package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/wenance/wequeue-management_api/app/config"
)

var dynamoEndpoint = config.Get("engines.AWS.dynamodb.endpoint")

func GetClient() dynamodbiface.DynamoDBAPI {
	var sess *session.Session
	if config.GetObject("aws_credentials") == nil {

		sess = session.Must(
			session.NewSession(&aws.Config{
				Region: aws.String("us-east-1"),
			}),
		)
	} else {
		sess = session.Must(
			session.NewSession(&aws.Config{
				Region: aws.String("us-east-1"),
				//Credentials: credentials.NewSharedCredentials("", "default"),
				Credentials: config.GetObject("aws_credentials").(*credentials.Credentials),
			}),
		)
	}
	/*sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: config.GetObject("aws_credentials").(*credentials.Credentials),
	})
	if err != nil {
		panic("FATAL: Connot connect to AWS")
	}*/
	dynamoClient := dynamodb.New(sess, aws.NewConfig().
		WithLogLevel(aws.LogDebugWithHTTPBody).
		WithEndpoint(dynamoEndpoint).
		WithDisableSSL(true),
	)
	if *config.GetCurrentEnvironment() == config.LOCAL || *config.GetCurrentEnvironment() == config.DEVELOP {
		dynamoClient.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(subscribersTable),
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

		dynamoClient.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(topicsTable),
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
