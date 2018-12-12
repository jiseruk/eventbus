package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/model"
)

var snsEndpoint = config.Get("engines.AWS.sns.endpoint")
var lambdaEndpoint = config.Get("engines.AWS.lambda.endpoint")
var kinesisEndpoint = config.Get("engines.AWS.kinesis.endpoint")
var sqsEndpoint = config.Get("engines.AWS.sqs.endpoint")

const AWS_RESOURCE_PREFIX = "wequeue-"

type PolicyDocument struct {
	Version   string
	Id        string
	Statement []StatementEntry
}
type Condition struct {
	ArnLike map[string]string
}

type StatementEntry struct {
	Sid       string
	Effect    string
	Action    []string
	Resource  string
	Principal map[string]string
	Condition Condition
}

type AWSEngine struct {
	SNSClient    snsiface.SNSAPI
	LambdaClient lambdaiface.LambdaAPI
	SQSClient    sqsiface.SQSAPI
}

type DeadLetterQueueInput struct {
	QueueArn  *string
	QueueName *string
}

type DLQSNSNotification struct {
	Records []map[string]map[string]interface{}
}

type SNSObject struct {
	Sns SNSNotification
}

type SNSNotification struct {
	Message          string
	MessageId        string
	TopicArn         string
	Type             string
	SignatureVersion string
	Signature        string
	SigningCertURL   string
	UnsubscribeURL   string
	Timestamp        string
}

func GetClients() (snsiface.SNSAPI, lambdaiface.LambdaAPI, kinesisiface.KinesisAPI, sqsiface.SQSAPI) {
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
	myCustomResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		switch service {
		case endpoints.SnsServiceID:
			if snsEndpoint != "" {
				return endpoints.ResolvedEndpoint{
					URL: snsEndpoint,
				}, nil
			}
		case endpoints.SqsServiceID:
			if sqsEndpoint != "" {
				return endpoints.ResolvedEndpoint{
					URL: sqsEndpoint,
				}, nil
			}
		case endpoints.LambdaServiceID:
			if lambdaEndpoint != "" {
				return endpoints.ResolvedEndpoint{
					URL: lambdaEndpoint,
				}, nil
			}
		case endpoints.KinesisServiceID:
			if kinesisEndpoint != "" {
				return endpoints.ResolvedEndpoint{
					URL: kinesisEndpoint,
				}, nil
			}
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	snsClient := sns.New(sess, aws.NewConfig().
		WithEndpointResolver(endpoints.ResolverFunc(myCustomResolver)),
	)

	lambdaClient := lambda.New(sess, aws.NewConfig().
		WithEndpointResolver(endpoints.ResolverFunc(myCustomResolver)),
	)

	kinesisClient := kinesis.New(sess, aws.NewConfig().
		WithEndpointResolver(endpoints.ResolverFunc(myCustomResolver)),
	)

	sqsClient := sqs.New(sess, aws.NewConfig().
		WithEndpointResolver(endpoints.ResolverFunc(myCustomResolver)),
	)
	return snsClient, lambdaClient, kinesisClient, sqsClient
}

//CreateTopic creates a topic in AWS SNS.
//It returns the topic Arn and any error encountered
func (azn *AWSEngine) CreateTopic(name string) (*CreateTopicOutput, error) {
	var input = &sns.CreateTopicInput{Name: aws.String(GetAWSResourcePrefix() + name)}
	snsoutput, err := azn.SNSClient.CreateTopic(input)
	if err != nil {
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

func (azn AWSEngine) Publish(topicResourceID string, message *model.PublishMessage) (*model.PublishMessage, error) {
	bytesMessage, _ := json.Marshal(&message)
	strMessage := string(bytesMessage)
	publishInput := &sns.PublishInput{Message: &strMessage, TopicArn: &topicResourceID}
	_, err := azn.SNSClient.Publish(publishInput)
	if err != nil {
		return nil, err
	}
	return message, nil
}

//CreateSubscriber creates a sns subscriber, that basically is a Lambda Function which receives the push notification and
//makes the HTTP POST to the subscriber's endpoint
func (azn AWSEngine) CreatePushSubscriber(topic model.Topic, subscriber string, endpoint string) (*SubscriberOutput, error) {
	qoutput, err := azn.SQSClient.CreateQueue(&sqs.CreateQueueInput{QueueName: aws.String(GetAWSResourcePrefix() + "dead-letter-" + subscriber)})
	if err != nil {
		return nil, err
	}
	qattrs, err := azn.SQSClient.GetQueueAttributes(
		&sqs.GetQueueAttributesInput{QueueUrl: qoutput.QueueUrl, AttributeNames: []*string{aws.String("QueueArn")}})
	if err != nil {
		//TODO: Borrar Queue
		return nil, err
	}

	lambdaConf, err := createLambdaSubscriber(azn.LambdaClient, topic.Name, subscriber,
		endpoint, "subscriber.handler", "python2.7", "/tmp/function.zip",
		&DeadLetterQueueInput{QueueName: aws.String(GetAWSResourcePrefix() + "dead-letter-" + subscriber), QueueArn: qattrs.Attributes["QueueArn"]})

	if err != nil {
		return nil, err
	}

	output, err := azn.SNSClient.Subscribe(&sns.SubscribeInput{TopicArn: &topic.ResourceID,
		Protocol: aws.String("lambda"),
		Endpoint: lambdaConf.FunctionArn,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber for topic ["+topic.ResourceID+"]")
	}

	if config.GetBool("engines.AWS.lambda.createPolicy") {
		_, err = azn.LambdaClient.AddPermission(&lambda.AddPermissionInput{
			Action:       aws.String("lambda:InvokeFunction"),
			FunctionName: aws.String(GetAWSResourcePrefix() + "lambda-" + subscriber),
			Principal:    aws.String("sns.amazonaws.com"),
			SourceArn:    &topic.ResourceID,
			StatementId:  aws.String(GetAWSResourcePrefix() + "lambda-policy-" + subscriber),
		})

		if err != nil {
			return nil, errors.Wrap(err, "Error creating subscriber's policy")
		}
	}
	return &SubscriberOutput{SubscriptionID: *output.SubscriptionArn, DeadLetterQueue: *qoutput.QueueUrl}, nil

}

func (azn AWSEngine) CreatePullSubscriber(topic model.Topic, subscriber string, visibilityTimeout int) (*SubscriberOutput, error) {
	qoutput, err := azn.SQSClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(GetAWSResourcePrefix() + "pull-queue-" + subscriber),
		Attributes: map[string]*string{
			"VisibilityTimeout": aws.String(strconv.Itoa(visibilityTimeout)),
			//"ReceiveMessageWaitTimeSeconds": aws.String("1"),
		},
	})
	if err != nil {
		return nil, err
	}
	qattrs, err := azn.SQSClient.GetQueueAttributes(
		&sqs.GetQueueAttributesInput{QueueUrl: qoutput.QueueUrl, AttributeNames: []*string{aws.String("QueueArn")}})
	if err != nil {
		//TODO: Borrar Queue
		return nil, err
	}
	output, err := azn.SNSClient.Subscribe(&sns.SubscribeInput{TopicArn: &topic.ResourceID,
		Protocol: aws.String("sqs"),
		Endpoint: qattrs.Attributes["QueueArn"]},
	)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber for topic ["+topic.ResourceID+"]")
	}

	_, err = azn.SQSClient.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		QueueUrl: qoutput.QueueUrl,
		Attributes: map[string]*string{
			"Policy": getPolicy(topic.ResourceID, *qattrs.Attributes["QueueArn"], GetAWSResourcePrefix()+"pull-queue-policy-"+subscriber),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating subscriber's policy")
	}

	return &SubscriberOutput{SubscriptionID: *output.SubscriptionArn, PullingQueue: *qoutput.QueueUrl}, nil

}

func (azn AWSEngine) ReceiveMessages(resourceID string, maxMessages int64) ([]model.Message, error) {
	output, err := azn.SQSClient.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: &resourceID, MaxNumberOfMessages: &maxMessages})
	if err != nil {
		return nil, err
	}
	messages := make([]model.Message, len(output.Messages))
	for i, msg := range output.Messages {
		buff := bytes.NewBufferString(*msg.Body)
		decoder := json.NewDecoder(buff)
		decoder.DisallowUnknownFields()
		var snsnotif SNSNotification
		err := decoder.Decode(&snsnotif)
		//Si es una dead-letter-queue de un push subscriber
		if err != nil {
			var dlqPayload DLQSNSNotification
			err := json.Unmarshal([]byte(*msg.Body), &dlqPayload)
			if err != nil {
				fmt.Printf("Error unmarshalling data %s", *msg.Body)
				return nil, err
			}
			err = mapstructure.Decode(dlqPayload.Records[0]["Sns"], &snsnotif)
			if err != nil {
				fmt.Print(err.Error())
				return nil, err
			}
		}

		var publishedMessage model.PublishMessage
		err = json.Unmarshal([]byte(snsnotif.Message), &publishedMessage)
		if err != nil {
			fmt.Printf("Error unmarshalling payload %s", snsnotif.Message)
			return nil, err
		}
		messages[i] = model.Message{
			Message:     publishedMessage,
			MessageID:   *msg.MessageId,
			DeleteToken: msg.ReceiptHandle,
		}
	}
	return messages, nil
}

func createLambdaSubscriber(client lambdaiface.LambdaAPI, topic string, subscriber string,
	endpoint string, handler string, runtime string, lambdaZipPath string,
	deadLetterQueueInfo *DeadLetterQueueInput) (*lambda.FunctionConfiguration, error) {

	contents, err := ioutil.ReadFile(lambdaZipPath)
	if err != nil {
		return nil, err
	}

	createCode := &lambda.FunctionCode{
		ZipFile: contents,
	}

	environment := lambda.Environment{Variables: make(map[string]*string)}
	environment.Variables["subscriber_url"] = &endpoint
	environment.Variables["topic"] = &topic
	environment.Variables["environment"] = config.GetCurrentEnvironment()

	createArgs := &lambda.CreateFunctionInput{
		Code:         createCode,
		FunctionName: aws.String(GetAWSResourcePrefix() + "lambda-" + subscriber),
		Handler:      &handler,
		Role:         aws.String(config.Get("engines.AWS.lambda.executionRole")),
		Runtime:      &runtime,
		Environment:  &environment,
		Tags:         map[string]*string{"project": aws.String("wequeue")},
	}
	if deadLetterQueueInfo != nil {
		createArgs.DeadLetterConfig = &lambda.DeadLetterConfig{TargetArn: deadLetterQueueInfo.QueueArn}
		createArgs.Environment.Variables["queue_name"] = deadLetterQueueInfo.QueueName
	}

	result, err := client.CreateFunction(createArgs)
	return result, err

}

func (azn AWSEngine) DeleteMessages(messages []model.Message, queueUrl string) ([]model.Message, error) {
	messagesToDelete := make([]*sqs.DeleteMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		messagesToDelete[i] = &sqs.DeleteMessageBatchRequestEntry{Id: &message.MessageID, ReceiptHandle: message.DeleteToken}
	}
	output, err := azn.SQSClient.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{Entries: messagesToDelete, QueueUrl: &queueUrl})
	if err != nil {
		return nil, err
	}
	failedDeleteMessages := make([]model.Message, 0)
	for _, errorMsg := range output.Failed {
		msg := model.Message{MessageID: *errorMsg.Id}
		msg.DeleteError.Code = errorMsg.Code
		msg.DeleteError.Message = errorMsg.Message
		failedDeleteMessages = append(failedDeleteMessages, model.Message{MessageID: *errorMsg.Id})
	}
	return failedDeleteMessages, nil
}

func getPolicy(snsTopicArn string, sqsArn string, sqsName string) *string {
	policy := PolicyDocument{
		Version: "2012-10-17",
		Id:      sqsArn + "/SQSDefaultPolicy",
		Statement: []StatementEntry{
			StatementEntry{
				Sid:    sqsName,
				Effect: "Allow",
				Action: []string{
					"sqs:SendMessage",
				},
				Resource: sqsArn,
				Principal: map[string]string{
					"Service": "sns.amazonaws.com",
				},
				Condition: Condition{
					ArnLike: map[string]string{"AWS:SourceArn": snsTopicArn},
				},
			},
		},
	}
	b, err := json.Marshal(&policy)
	if err != nil {
		fmt.Println("Error marshaling policy", err)
		return nil
	}
	return aws.String(string(b))
}

func GetAWSResourcePrefix() string {
	return AWS_RESOURCE_PREFIX + *config.GetCurrentEnvironment() + "-"
}

/*
func (azn AWSEngine) InvokeLambda(name string, payload string) (*lambda.InvokeOutput, error) {
	output, err := azn.LambdaClient.Invoke(&lambda.InvokeInput{
		Payload:      []byte(payload),
		FunctionName: &name,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("%#v", output)
	return output, nil
}
*/
