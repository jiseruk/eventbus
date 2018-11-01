package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/pkg/errors"
	"github.com/wenance/wequeue-management_api/app/model"
)

//AWSStreamEngine is a AWS WeQueue implementation, on top of SNS, SQS and Lambda
type AWSStreamEngine struct {
	LambdaClient  lambdaiface.LambdaAPI
	KinesisClient kinesisiface.KinesisAPI
}

func (azn *AWSStreamEngine) CreateTopic(name string) (*CreateTopicOutput, error) {
	_, err := azn.KinesisClient.CreateStream(&kinesis.CreateStreamInput{ShardCount: aws.Int64(1), StreamName: &name})
	if err != nil {
		fmt.Printf("Error: %#v", err)
		return nil, err
	}
	output, err := azn.KinesisClient.DescribeStream(&kinesis.DescribeStreamInput{StreamName: &name})
	if err != nil {
		return nil, err
	}
	return &CreateTopicOutput{Resource: *output.StreamDescription.StreamARN}, nil

}

func (azn AWSStreamEngine) CreateSubscriber(topic model.Topic, subscriber string, endpoint string) (*SubscriberOutput, error) {
	lambdaConf, err := createLambdaSubscriber(azn.LambdaClient, topic.Name, subscriber, endpoint,
		"consumer.handler", "python2.7", "/tmp/function.zip", nil)
	if err != nil {
		return nil, err
	}

	output, err := azn.LambdaClient.CreateEventSourceMapping(&lambda.CreateEventSourceMappingInput{FunctionName: lambdaConf.FunctionName,
		EventSourceArn: &topic.ResourceID, BatchSize: aws.Int64(1)})

	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber")
	}
	return &SubscriberOutput{SubscriptionID: *output.UUID, PullResourceID: topic.Name}, nil

}

func (azn AWSStreamEngine) Publish(topicResourceID string, message interface{}) (*PublishOutput, error) {
	bytesMessage, _ := json.Marshal(message)
	topic := strings.Split(topicResourceID, "/")[1]
	publishInput := &kinesis.PutRecordInput{StreamName: &topic,
		Data:         bytesMessage,
		PartitionKey: aws.String("*"),
	}

	output, err := azn.KinesisClient.PutRecord(publishInput)
	if err != nil {
		return nil, err
	}
	return &PublishOutput{MessageID: *output.SequenceNumber}, nil
}

func (azn *AWSStreamEngine) DeleteTopic(resource string) error {
	_, err := azn.KinesisClient.DeleteStream(&kinesis.DeleteStreamInput{StreamName: &resource})
	return err
}

func (azn AWSStreamEngine) ReceiveMessages(resourceID string, maxMessages int64) (*model.Messages, error) {
	shards, err := azn.KinesisClient.GetShardIterator(&kinesis.GetShardIteratorInput{StreamName: &resourceID,
		ShardId: aws.String("1")})

	if err != nil {
		return nil, err
	}
	output, err := azn.KinesisClient.GetRecords(&kinesis.GetRecordsInput{ShardIterator: shards.ShardIterator, Limit: &maxMessages})
	if err != nil {
		return nil, err
	}
	messages := &model.Messages{Messages: make([]model.Message, len(output.Records))}
	for i, record := range output.Records {
		messages.Messages[i] = model.Message{Payload: record.Data, MessageID: *record.SequenceNumber}

	}

	return messages, nil
}

func (azn AWSStreamEngine) DeleteMessages(messages []model.Message, queueUrl string) ([]*model.Message, error) {
	return nil, errors.New("You can't delete messages in this stream topic")
}

func (azn AWSStreamEngine) GetName() string {
	return "AWSStream"
}
