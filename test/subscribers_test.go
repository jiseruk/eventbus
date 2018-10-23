package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	"testing"
)

func TestCreateSubscription(t *testing.T){
	router := server.GetRouter()
	t.Run("It should create the subscription", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		mockDAO := &SubscriptionsDaoMock{}
		service.TopicsService = topicServiceMock
		service.SubscriptionsService = service.SubscriptionServiceImpl{Db: mockDAO}
		mockLambda := &LambdaAPIMock{}
		mockSNS := &SNSAPIMock{}
		subscriber := &model.Subscriber{Endpoint:"/endp", Name:"subs", Topic:"topic", ResourceID:"arn:subs"}
		client.EnginesMap["AWS"] = &client.AWSEngine{LambdaClient:mockLambda }

		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name:"topic", Engine:"AWS"}, nil).Once()
		mockLambda.On("CreateFunction", mock.AnythingOfType("lambda.CreateFunctionInput")).
			Return(&lambda.FunctionConfiguration{FunctionArn:aws.String("func:arn")}).Once()
		mockSNS.On("Subscribe",
			&sns.SubscribeInput{Endpoint:aws.String("func:arn"),
								TopicArn:aws.String("arn:topic"),
								Protocol:aws.String("lambda")}).
			Return(&sns.SubscribeOutput{SubscriptionArn: aws.String("arn:subs")}).Once()
		mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()
		mockDAO.On("CreateSubscription", "subs", "topic:arn", "/endp", "arn:subs").
			Return(subscriber, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscriptions", `{"topic": "topic", "name":"subs", "endpoint":"/endp"}`)

		assert.Equal(t, 201, rec.Code)
	})
}
