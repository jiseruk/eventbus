package test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	_ "github.com/wenance/wequeue-management_api/app/validation"
)

func TestCreateSubscription(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	router := server.GetRouter()
	//For lambda creation
	ioutil.WriteFile("/tmp/function.zip", []byte("data loca"), 0644)
	mockDAO := &SubscriptionsDaoMock{}
	service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}

	t.Run("It should create a subscriber to a topic", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		mockLambda := &LambdaAPIMock{}
		mockSNS := &SNSAPIMock{}
		mockSQS := &SQSAPIMock{}

		subscriber := &model.Subscriber{Endpoint: "http://subscriber/endp",
			Name: "subs", Topic: "topic",
			ResourceID:      "arn:subs",
			DeadLetterQueue: "queue:subs",
			CreatedAt:       model.Clock.Now(),
		}
		//The endpoint should be an active http endpoint
		client.HTTPClient = NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), subscriber.Endpoint)
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			}, nil
		})

		lambdaArgs := getLambdaMock("http://subscriber/endp", "subs", "topic", "queue:subs")
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
		mockDAO.On("CreateSubscription", "subs", "topic", "http://subscriber/endp", "arn:subs", "queueUrl").
			Return(subscriber, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp"}`)

		assert.Equal(t, 201, rec.Code)
		assert.JSONEq(t, `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp"}`,
			rec.Body.String())

		topicServiceMock.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
		mockSQS.AssertExpectations(t)
		mockLambda.AssertExpectations(t)
	})

	t.Run("it should fail creating the subscriber if the endpoint is not a valid url", func(t *testing.T) {
		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"/endp"}`)
		assert.Equal(t, 400, rec.Code)
	})

	for _, rtFn := range []RoundTripFunc{
		func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), "http://subscriber/endp")
			return nil, errors.New("No response")
		},
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.URL.String(), "http://subscriber/endp")
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`Error`)),
			}, nil
		},
	} {
		t.Run("it should fail creating the subscriber if the endpoint is down", func(t *testing.T) {
			topicServiceMock := &TopicServiceMock{}
			service.TopicsService = topicServiceMock
			client.HTTPClient = NewTestClient(rtFn)
			//Topic should exist
			topicServiceMock.On("GetTopic", "topic").
				Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
			//Subscriber with the provided name shold not exist in DB
			mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()

			rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp"}`)
			assert.Equal(t, 400, rec.Code)

			topicServiceMock.AssertExpectations(t)
		})
	}

	t.Run("it should fail creating the subscriber if the topic doesn't exist", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		topicServiceMock.On("GetTopic", "topic").Return(nil, nil).Once()
		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp"}`)

		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, `{"message":"The topic topic doesn't exist","code":"validation_error","status":400}`,
			rec.Body.String())
	})

	for _, r := range []struct {
		body string
		err  string
	}{
		{body: `{"topic": "topic", "endpoint": "lala"}`, err: "name is not present"},
		{body: `{"topic": "topic", "endpoint": "lala", "name":null}`, err: "name is null"},
		{body: `{"name": 10, "endpoint": "http://www.ole.com.ar", "topic": "topic"}`, err: "name is numeric"},
		{body: `{"name": "", "endpoint": "http://www.ole.com.ar", "topic": "topic"}`, err: "name is empty"},
		{body: `{"name": "subscriber", "topic": "topic"}`},
		{body: `{"name": "subscriber", "endpoint": 8}`},
		{body: `{"invalid": "topic", "invalid2": "lala"}`},
		{body: `{"topic": "topic"}`},
		{body: `{"name": "subscriber", "topic":"", "endpoint": "http://www.ole.com.ar"}`, err: "topic is empty"},
		{body: `{}`},
		{body: ``},
	} {
		t.Run("it should fail create the subscriber if the json fields are invalid ["+r.err+"]", func(t *testing.T) {

			res := executeMockedRequest(router, "POST", "/subscribers", r.body)
			assert.Contains(t, res.Body.String(), `"code":"json_error"`)
			fmt.Printf(res.Body.String())
			assert.Equal(t, 400, res.Code, res.Body.String())
		})
	}
}

func TestConsumeDeadLetterQueueMessages(t *testing.T) {
	router := server.GetRouter()
	subscriber := &model.Subscriber{Endpoint: "http://subscriber/endp", Name: "subs", Topic: "topic",
		ResourceID: "arn:subs", DeadLetterQueue: "queue:subs"}

	t.Run("It should get messages from dead letter queue", func(t *testing.T) {
		mockSQS := &SQSAPIMock{}
		mockDAO := &SubscriptionsDaoMock{}

		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
		client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}
		mockDAO.On("GetSubscription", "subs").
			Return(subscriber, nil).Once()
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()

		mockSQS.On("ReceiveMessage", &sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(10), QueueUrl: aws.String("queue:subs")}).
			Return(&sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{Body: aws.String(`{"msg":"lala"}`), MessageId: aws.String("1"), ReceiptHandle: aws.String("x")},
				},
			}, nil).Once()

		res := executeMockedRequest(router, "GET", "/messages?subscriber=subs&max_messages=10", "")
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, `{"topic":"topic","messages":[{"payload":{"msg":"lala"},"message_id":"1","delete_token":"x"}]}`, res.Body.String())
		mockSQS.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
		topicServiceMock.AssertExpectations(t)

	})

	t.Run("it should fail consuming dead letter queue messages if the subscriber doesn't exist", func(t *testing.T) {
		mockDAO := &SubscriptionsDaoMock{}
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
		mockDAO.On("GetSubscription", "subs").
			Return(nil, nil).Once()

		res := executeMockedRequest(router, "GET", "/messages?subscriber=subs&max_messages=10", "")
		assert.Equal(t, 404, res.Code)
		assert.Equal(t, `{"message":"The subscriber subs doesn't exist","code":"database_error","status":404}`, res.Body.String())
		mockDAO.AssertExpectations(t)
	})

	t.Run("it should fail consuming dead letter queue messages if the topic doesn't exist", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		mockDAO := &SubscriptionsDaoMock{}
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
		mockDAO.On("GetSubscription", "subs").
			Return(subscriber, nil).Once()
		topicServiceMock.On("GetTopic", subscriber.Topic).Return(nil, nil).Once()

		res := executeMockedRequest(router, "GET", "/messages?subscriber=subs&max_messages=10", "")
		assert.Equal(t, 404, res.Code)
		assert.Equal(t, `{"message":"The topic topic doesn't exist","code":"database_error","status":404}`, res.Body.String())
		mockDAO.AssertExpectations(t)
		topicServiceMock.AssertExpectations(t)
	})

	for _, r := range []struct {
		query string
		err   string
	}{
		{query: `subscriber=&max_messages=`, err: "subscriber/max_messages is empty"},
		{query: `subscriber=lala&max_messages=11`, err: "max_messages > 10"},
		{query: `subscriber=lala&max_messages=0`, err: "max_messages < 1"},
		{query: ``, err: "subscriber/max_messages are empty"},
	} {
		t.Run("It should fail getting messages from dead letter queue if query params are invalid ["+r.err+"]", func(t *testing.T) {
			res := executeMockedRequest(router, "GET", "/messages?"+r.query, "")
			assert.Equal(t, 400, res.Code)
		})
	}
}
func TestDeleteMessages(t *testing.T) {
	router := server.GetRouter()

	t.Run("It should delete specific messages from dead letter queue", func(t *testing.T) {
		mockSQS := &SQSAPIMock{}
		mockDAO := &SubscriptionsDaoMock{}
		subscriber := &model.Subscriber{Endpoint: "http://subscriber/endp", Name: "subs", Topic: "topic",
			ResourceID: "arn:subs", DeadLetterQueue: "queue:subs"}

		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
		client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}
		mockDAO.On("GetSubscription", "subs").
			Return(subscriber, nil).Once()
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()

		mockSQS.On("DeleteMessageBatch", &sqs.DeleteMessageBatchInput{
			Entries: []*sqs.DeleteMessageBatchRequestEntry{
				{Id: aws.String("1"), ReceiptHandle: aws.String("x")},
			}, QueueUrl: aws.String("queue:subs")}).
			Return(&sqs.DeleteMessageBatchOutput{
				Failed: make([]*sqs.BatchResultErrorEntry, 0),
			}, nil).Once()

		res := executeMockedRequest(router, "DELETE", "/messages", `{"subscriber":"subs", "messages": [{"message_id":"1", "delete_token":"x"}]}`)
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, `{"failed":[],"topic":"topic"}`, res.Body.String())
		mockSQS.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
		topicServiceMock.AssertExpectations(t)

	})
}
