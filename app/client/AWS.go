package client

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/pkg/errors"
	"io/ioutil"
)

const SNSEndpoint = "http://localstack:4575"

type AWSEngine struct {
	SNSClient snsiface.SNSAPI
	LambdaClient lambdaiface.LambdaAPI
}

func GetClients() (snsiface.SNSAPI, lambdaiface.LambdaAPI) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		//Credentials: credentials.NewSharedCredentials("", "default"),
		Credentials: credentials.NewStaticCredentials("foo", "bar", ""),
		//Endpoint:    aws.String(SNSEndpoint),
	})
	if err != nil {

	}
	svc := sns.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint(SNSEndpoint))
	lambda := lambda.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint("http://localstack:4574"))
	return svc, lambda
}

//CreateTopic lala
func (azn *AWSEngine) CreateTopic(name string) (*CreateTopicOutput, error) {
	// Make svc.AddPermission request
	var input = &sns.CreateTopicInput{Name: &name}
	snsoutput, err := azn.SNSClient.CreateTopic(input)
	if err != nil {
		fmt.Printf("Error: %#v", err)
		return nil, err
	}
	output := &CreateTopicOutput{Resource: *snsoutput.TopicArn}
	return output, nil
}

func (azn AWSEngine) GetName() string {
	return "AWS"
}

func (azn AWSEngine) Publish(topicResourceID string, message interface{}) (*PublishOutput, error){
	strMessage := message.(string)
	publishInput := &sns.PublishInput{Message: &strMessage, TopicArn: &topicResourceID}
	output, err := azn.SNSClient.Publish(publishInput)
	if err != nil {
		return nil, err
	}
	return &PublishOutput{MessageID: *output.MessageId}, nil
}

func (azn AWSEngine) CreateSubscriber(topicResourceID string, subscriber string, endpoint string) (*SubscriberOutput, error) {
	lambdaConf, err := azn.createLambdaSubscriber(topicResourceID, subscriber, endpoint)
	if err != nil {
		return nil, err
	}

	output, err := azn.SNSClient.Subscribe(&sns.SubscribeInput{TopicArn:&topicResourceID,
		Protocol:aws.String("lambda"),
		Endpoint: lambdaConf.FunctionArn},
	)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber")
	}
	return &SubscriberOutput{*output.SubscriptionArn}, nil

}

func (azn AWSEngine)createLambdaSubscriber(topic string, subscriber string, endpoint string) (*lambda.FunctionConfiguration, error){
	//zipFileName := "lambda_subscriber"
	zipFileName := "function"
	contents, err := ioutil.ReadFile("/go/src/github.com/wenance/wequeue-management_api/app/lambda/" + zipFileName + ".zip")
	if err != nil {
		return nil, err
	}
	//contents = []byte("lalala")
	createCode := &lambda.FunctionCode{
		//S3Bucket:        aws.String("EventBusSubscribers"),
		//S3Key:           aws.String(zipFileName),
		//S3ObjectVersion: aws.String("1"),
		ZipFile:         contents,
	}

	environment:=  lambda.Environment{Variables: make(map[string]*string)}
	environment.Variables["subscriber_url"] = &endpoint
	environment.Variables["topic"] = &topic

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String("lambda_subscriber"),
		Handler:      aws.String("index.handler"),
		Role:         aws.String("arn:role:dummy"),
		Runtime:      aws.String("nodejs8.10"),
		Environment:  &environment,
	}
	fmt.Printf("Function input: %#v", *createArgs)
	result, err := azn.LambdaClient.CreateFunction(createArgs)

	if err != nil {
		fmt.Println("Cannot create function: " + err.Error())
		return nil, err
	} else {
		fmt.Println(result)
	}
	return result, nil
}