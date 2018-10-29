package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	sqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
)

func TestCreateSubscription(t *testing.T) {
	os.Setenv("GO_ENVIRONMENT", config.LOCAL)
	model.Clock = clockwork.NewFakeClock()
	router := server.GetRouter()
	//For lambda creation
	ioutil.WriteFile("/tmp/function.zip", []byte("data loca"), 0644)

	t.Run("It should create the subscription in aws and DB", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		mockDAO := &SubscriptionsDaoMock{}
		service.TopicsService = topicServiceMock
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
		mockLambda := &LambdaAPIMock{}
		mockSNS := &SNSAPIMock{}
		mockSQS := &SQSAPIMock{}
		subscriber := &model.Subscriber{Endpoint: "/endp", Name: "subs", Topic: "topic",
			ResourceID: "arn:subs", PullResourceID: "queue:subs", CreatedAt: model.Clock.Now()}
		lambdaArgs := getLambdaMock("/endp", "subs", "topic", "queue:subs")
		client.EnginesMap["AWS"] = &client.AWSEngine{LambdaClient: mockLambda, SNSClient: mockSNS, SQSClient: mockSQS}
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()
		//The lambda function is created in AWS
		mockLambda.On("CreateFunction", lambdaArgs).
			Return(&lambda.FunctionConfiguration{FunctionArn: aws.String("func:arn")}, nil).Once()
		//The lambda is subscriber to the topic
		mockSNS.On("Subscribe",
			&sns.SubscribeInput{Endpoint: aws.String("func:arn"),
				TopicArn: aws.String("arn:topic"),
				Protocol: aws.String("lambda")}).
			Return(&sns.SubscribeOutput{SubscriptionArn: aws.String("arn:subs")}, nil).Once()

		mockSQS.On("CreateQueue", &sqs.CreateQueueInput{QueueName: aws.String("dlq_lambda_subs")}).
			Return(&sqs.CreateQueueOutput{QueueUrl: aws.String("queueUrl")}, nil).Once()
		mockSQS.On("GetQueueAttributes", &sqs.GetQueueAttributesInput{QueueUrl: aws.String("queueUrl"), AttributeNames: []*string{aws.String("QueueArn")}}).
			Return(&sqs.GetQueueAttributesOutput{Attributes: map[string]*string{"QueueArn": aws.String("queue:subs")}}, nil).Once()
		//Finnaly, The subscriber is created in the database
		mockDAO.On("CreateSubscription", "subs", "topic", "/endp", "arn:subs", "queueUrl").
			Return(subscriber, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscriptions", `{"topic": "topic", "name":"subs", "endpoint":"/endp"}`)

		assert.Equal(t, 201, rec.Code)
		assert.JSONEq(t,
			fmt.Sprintf(`{"topic": "topic", "name":"subs", "endpoint":"/endp", "resource_id":"arn:subs", "created_at": "%s"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")),
			rec.Body.String())
		topicServiceMock.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
		mockSQS.AssertExpectations(t)
		mockLambda.AssertExpectations(t)
	})

	//Endpoint should be unique
	//Endpoint shold be a valid url

}
