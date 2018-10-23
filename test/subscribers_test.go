package test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	"io/ioutil"
	"testing"
)

func TestCreateSubscription(t *testing.T){
	model.Clock = clockwork.NewFakeClock()
	router := server.GetRouter()
	//For lambda creation
	ioutil.TempFile("/tmp", "function.zip")

	t.Run("It should create the subscription in aws and DB", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		mockDAO := &SubscriptionsDaoMock{}
		service.TopicsService = topicServiceMock
		service.SubscriptionsService = service.SubscriptionServiceImpl{Db: mockDAO}
		mockLambda := &LambdaAPIMock{}
		mockSNS := &SNSAPIMock{}
		subscriber := &model.Subscriber{Endpoint:"/endp", Name:"subs", Topic:"topic", ResourceID:"arn:subs", CreatedAt: model.Clock.Now()}
		client.EnginesMap["AWS"] = &client.AWSEngine{LambdaClient:mockLambda, SNSClient:mockSNS }
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name:"topic", Engine:"AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()
		//The lambda function is created in AWS
		mockLambda.On("CreateFunction", mock.AnythingOfType("*lambda.CreateFunctionInput")).
			Return(&lambda.FunctionConfiguration{FunctionArn:aws.String("func:arn")}, nil).Once()
		//The lambda is subscriber to the topic
		mockSNS.On("Subscribe",
			&sns.SubscribeInput{Endpoint:aws.String("func:arn"),
								TopicArn:aws.String("arn:topic"),
								Protocol:aws.String("lambda")}).
			Return(&sns.SubscribeOutput{SubscriptionArn: aws.String("arn:subs")}, nil).Once()
		//Finnaly, The subscriber is created in the database
		mockDAO.On("CreateSubscription", "subs", "arn:topic", "/endp", "arn:subs").
			Return(subscriber, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscriptions", `{"topic": "topic", "name":"subs", "endpoint":"/endp"}`)

		assert.Equal(t, 201, rec.Code)
		assert.JSONEq(t,
			fmt.Sprintf(`{"topic": "topic", "name":"subs", "endpoint":"/endp", "resource_id":"arn:subs", "created_at": "%s"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")) ,
			rec.Body.String())
		topicServiceMock.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
		mockLambda.AssertExpectations(t)
	})
}
