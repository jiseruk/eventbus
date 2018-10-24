package client

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
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
	KinesisClient kinesisiface.KinesisAPI
}

func GetClients() (snsiface.SNSAPI, lambdaiface.LambdaAPI, kinesisiface.KinesisAPI) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		//Credentials: credentials.NewSharedCredentials("", "default"),
		Credentials: credentials.NewStaticCredentials("foo", "bar", ""),
		//Endpoint:    aws.String(SNSEndpoint),
	})
	if err != nil {
		panic("FATAL: Connot connect to AWS")
	}
	snsClient := sns.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint(SNSEndpoint))
	lambdaClient := lambda.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint("http://localstack:4574"))
	kinesisClient := kinesis.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody).WithEndpoint("http://localstack:4574"))
	return snsClient, lambdaClient, kinesisClient
}

//CreateTopic lala
func (azn *AWSEngine) CreateTopic(name string) (*CreateTopicOutput, error) {
	var input = &sns.CreateTopicInput{Name: &name}
	snsoutput, err := azn.SNSClient.CreateTopic(input)
	if err != nil {
		fmt.Printf("Error: %#v", err)
		return nil, err
	}
	output := &CreateTopicOutput{Resource: *snsoutput.TopicArn}
	return output, nil
}

func (azn *AWSEngine) DeleteTopic(resource string) error {
	_, err := azn.SNSClient.DeleteTopic(&sns.DeleteTopicInput{TopicArn: &resource})
	return err
}

func (azn AWSEngine) GetName() string {
	return "AWS"
}

func (azn AWSEngine) Publish(topicResourceID string, message interface{}) (*PublishOutput, error){
	bytesMessage, _ := json.Marshal(message)
	strMessage := string(bytesMessage)
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

	contents, err := ioutil.ReadFile("/tmp/function.zip")
	if err != nil {
		return nil, err
	}

	createCode := &lambda.FunctionCode{
		ZipFile:         contents,
	}

	environment:= lambda.Environment{Variables: make(map[string]*string)}
	environment.Variables["subscriber_url"] = &endpoint
	environment.Variables["topic"] = &topic

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String("lambda_subscriber_" + subscriber),
		//Handler:      aws.String("index.handler"),
		Handler:      aws.String("subscriber.handler"),
		Role:         aws.String("arn:role:dummy"),
		//Runtime:      aws.String("nodejs8.10"),
		Runtime:      aws.String("python2.7"),
		Environment:  &environment,
	}
	fmt.Printf("Function input: %#v", *createArgs)
	result, err := azn.LambdaClient.CreateFunction(createArgs)
	//azn.LambdaClient.Invoke(lambda.InvokeInput{})
	if err != nil {
		fmt.Println("Cannot create function: " + err.Error())
		return nil, err
	} else {
		fmt.Println(result)
	}
	return result, nil
}

func (azn AWSEngine) InvokeLambda(name string, payload string) (*lambda.InvokeOutput, error){
	output, err := azn.LambdaClient.Invoke(&lambda.InvokeInput{
		Payload: []byte(payload),
		FunctionName: &name,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", output)
	return output, nil
}

func (azn *AWSEngine) CreateStream(name string) (*CreateTopicOutput, error) {
	koutput, err := azn.KinesisClient.CreateStream(&kinesis.CreateStreamInput{ShardCount: aws.Int64(1), StreamName: &name})
	if err != nil {
		fmt.Printf("Error: %#v", err)
		return nil, err
	}

	output := &CreateTopicOutput{Resource: koutput.String()}
	return output, nil

}

func (azn AWSEngine) CreateKinesisSubscriber(topicResourceID string, subscriber string, endpoint string) (*SubscriberOutput, error) {
	lambdaConf, err := azn.createLambdaSubscriber(topicResourceID, subscriber, endpoint)
	if err != nil {
		return nil, err
	}

	output, err := azn.LambdaClient.CreateEventSourceMapping(&lambda.CreateEventSourceMappingInput{FunctionName:lambdaConf.FunctionName,
		EventSourceArn:&topicResourceID, BatchSize: aws.Int64(1)})

	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber")
	}
	return &SubscriberOutput{*output.UUID}, nil

}

func (azn AWSEngine) PublishKinesis(topicResourceID string, message interface{}) (*PublishOutput, error){
	bytesMessage, _ := json.Marshal(message)
	publishInput := &kinesis.PutRecordInput{StreamName: &topicResourceID, Data:bytesMessage}
	output, err := azn.KinesisClient.PutRecord(publishInput)
	if err != nil {
		return nil, err
	}
	return &PublishOutput{MessageID: *output.SequenceNumber}, nil
}
