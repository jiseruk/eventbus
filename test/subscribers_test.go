package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jonboulle/clockwork"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wenance/wequeue-management_api/app/client"
	apierrors "github.com/wenance/wequeue-management_api/app/errors"
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

	mockDynamo := &DynamoDBAPIMock{}
	//service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: &model.SubscriberDaoDynamoImpl{DynamoClient: mockDynamo},
	}

	t.Run("It should create a push subscriber to a topic", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		mockLambda := &LambdaAPIMock{}
		mockSNS := &SNSAPIMock{}
		mockSQS := &SQSAPIMock{}

		subscriber := &model.Subscriber{Endpoint: aws.String("http://subscriber/endp"),
			Name:            "subs",
			Topic:           "topic",
			ResourceID:      "arn:subs",
			Type:            "push",
			DeadLetterQueue: "queueUrl",
			CreatedAt:       model.Clock.Now(),
		}
		subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)
		//The endpoint should be an active http endpoint
		client.HTTPClient = NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), *subscriber.Endpoint)
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
		//Subscriber with the provided name should not exist in DB
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == subscriber.Name && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()
		//Endpoint should be unique
		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.ExpressionAttributeValues[":e"].S == *subscriber.Endpoint &&
				*input.TableName == "Subscribers"
		})).Return(&dynamodb.ScanOutput{Items: nil}, nil).Once()

		//mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()
		//The lambda function is created in AWS
		mockLambda.On("CreateFunction", lambdaArgs).
			Return(&lambda.FunctionConfiguration{FunctionArn: aws.String("func:arn")}, nil).Once()
		//The lambda is subscriber to the topic
		mockSNS.On("Subscribe",
			&sns.SubscribeInput{Endpoint: aws.String("func:arn"),
				TopicArn: aws.String("arn:topic"),
				Protocol: aws.String("lambda")}).
			Return(&sns.SubscribeOutput{SubscriptionArn: aws.String("arn:subs")}, nil).Once()

		mockSQS.On("CreateQueue", &sqs.CreateQueueInput{QueueName: aws.String(client.GetAWSResourcePrefix() + "dead-letter-subs")}).
			Return(&sqs.CreateQueueOutput{QueueUrl: aws.String("queueUrl")}, nil).Once()
		mockSQS.On("GetQueueAttributes", &sqs.GetQueueAttributesInput{QueueUrl: aws.String("queueUrl"), AttributeNames: []*string{aws.String("QueueArn")}}).
			Return(&sqs.GetQueueAttributesOutput{Attributes: map[string]*string{"QueueArn": aws.String("queue:subs")}}, nil).Once()

		mockLambda.On("AddPermission", mock.MatchedBy(func(input *lambda.AddPermissionInput) bool {
			return *input.FunctionName == client.GetAWSResourcePrefix()+"lambda-subs" &&
				*input.StatementId == client.GetAWSResourcePrefix()+"lambda-policy-subs"
		})).Return(&lambda.AddPermissionOutput{}, nil).Once()
		//Finnaly, The subscriber is created in the database
		mockDynamo.On("PutItem", &dynamodb.PutItemInput{
			Item:      subscriberItem,
			TableName: aws.String("Subscribers"),
		}).Return(&dynamodb.PutItemOutput{}, nil).Once()

		/* mockDAO.On("CreateSubscription", "subs", "topic", "push", "arn:subs", aws.String("http://subscriber/endp"), "queueUrl", "").
		Return(subscriber, nil).Once() */

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`)

		assert.Equal(t, 201, rec.Code)
		assert.JSONEq(t, `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`,
			rec.Body.String())

		topicServiceMock.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
		mockSQS.AssertExpectations(t)
		mockLambda.AssertExpectations(t)
	})

	t.Run("It should create a pull subscriber to a topic", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		mockSNS := &SNSAPIMock{}
		mockSQS := &SQSAPIMock{}

		subscriber := &model.Subscriber{Name: "subs",
			Topic:             "topic",
			Type:              "pull",
			ResourceID:        "arn:subs",
			PullingQueue:      "queueUrl",
			VisibilityTimeout: aws.Int(10),
			CreatedAt:         model.Clock.Now(),
		}
		subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)

		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS, SQSClient: mockSQS}

		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == subscriber.Name && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()
		//mockDAO.On("GetSubscription", "subs").Return(nil, nil).Once()

		mockSQS.On("CreateQueue", &sqs.CreateQueueInput{QueueName: aws.String(client.GetAWSResourcePrefix() + "pull-queue-subs"),
			Attributes: map[string]*string{"VisibilityTimeout": aws.String("10")},
		}).Return(&sqs.CreateQueueOutput{QueueUrl: aws.String("queueUrl")}, nil).Once()
		mockSQS.On("GetQueueAttributes", &sqs.GetQueueAttributesInput{QueueUrl: aws.String("queueUrl"), AttributeNames: []*string{aws.String("QueueArn")}}).
			Return(&sqs.GetQueueAttributesOutput{Attributes: map[string]*string{"QueueArn": aws.String("queue:subs")}}, nil).Once()
		//The sqs queue is subscribed to the topic
		mockSNS.On("Subscribe",
			&sns.SubscribeInput{Endpoint: aws.String("queue:subs"),
				TopicArn: aws.String("arn:topic"),
				Protocol: aws.String("sqs")}).
			Return(&sns.SubscribeOutput{SubscriptionArn: aws.String("arn:subs")}, nil).Once()

		mockSQS.On("SetQueueAttributes", mock.MatchedBy(func(input *sqs.SetQueueAttributesInput) bool {
			return true
		})).Return(&sqs.SetQueueAttributesOutput{}, nil).Once()
		//Finnaly, The subscriber is created in the database
		mockDynamo.On("PutItem", &dynamodb.PutItemInput{
			Item:      subscriberItem,
			TableName: aws.String("Subscribers"),
		}).Return(&dynamodb.PutItemOutput{}, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "type":"pull", "visibility_timeout":10}`)

		assert.Equal(t, 201, rec.Code)
		assert.JSONEq(t, `{"topic": "topic", "name":"subs", "type":"pull", "visibility_timeout":10}`,
			rec.Body.String())

		topicServiceMock.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
		mockSQS.AssertExpectations(t)
	})

	t.Run("it should fail creating the push subscriber if the endpoint is not a valid url", func(t *testing.T) {
		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"/endp", "type":"push"}`)
		assert.Equal(t, 400, rec.Code)
	})

	t.Run("it should fail creating the push subscriber if the endpoint is already in use by another subscriber", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		subscriber := &model.Subscriber{Name: "subs",
			Topic:      "topic",
			Type:       "push",
			ResourceID: "arn:subs",
			Endpoint:   aws.String("http://endpoint/"),
			CreatedAt:  model.Clock.Now(),
		}

		subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == subscriber.Name && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.ExpressionAttributeValues[":e"].S == *subscriber.Endpoint &&
				*input.TableName == "Subscribers"
		})).Return(&dynamodb.ScanOutput{Items: []map[string]*dynamodb.AttributeValue{subscriberItem}}, nil).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://endpoint/", "type":"push"}`)
		assert.Equal(t, 400, rec.Code)
		topicServiceMock.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)

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
		t.Run("it should fail creating the push subscriber if the endpoint is down", func(t *testing.T) {
			topicServiceMock := &TopicServiceMock{}
			service.TopicsService = topicServiceMock
			client.HTTPClient = NewTestClient(rtFn)
			//Topic should exist
			topicServiceMock.On("GetTopic", "topic").
				Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
			//Subscriber with the provided name shold not exist in DB
			mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
				return *input.Key["name"].S == "subs" && *input.TableName == "Subscribers"
			})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

			mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
				return *input.ExpressionAttributeValues[":e"].S == "http://subscriber/endp" &&
					*input.TableName == "Subscribers"
			})).Return(&dynamodb.ScanOutput{Items: nil}, nil).Once()

			rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`)
			assert.Equal(t, 400, rec.Code)

			topicServiceMock.AssertExpectations(t)
			mockDynamo.AssertExpectations(t)
		})
	}

	t.Run("it should fail creating the subscriber if the topic doesn't exist", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		topicServiceMock.On("GetTopic", "topic").Return(nil, nil).Once()
		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`)

		assert.Equal(t, 400, rec.Code)
		assert.Equal(t, `{"message":"The topic topic doesn't exist","code":"validation_error","status":400}`,
			rec.Body.String())
	})

	t.Run("it should fail creating the subscriber when a getTopic() error happends", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		topicServiceMock.On("GetTopic", "topic").
			Return(nil, apierrors.NewAPIError(500, "database_error", "Error getting topic")).Once()
		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`)

		assert.Equal(t, 500, rec.Code)
		assert.Equal(t, `{"message":"Error getting topic","code":"database_error","status":500}`,
			rec.Body.String())
	})

	t.Run("it should fail creating the subscriber when a dao.getSubscriber() error happends", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()

		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == "subs" && *input.TableName == "Subscribers"
		})).Return(nil, errors.New("Dynamodb error")).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://subscriber/endp", "type":"push"}`)

		assert.Equal(t, 500, rec.Code)
		assert.Equal(t, `{"message":"Dynamodb error","code":"database_error","status":500}`,
			rec.Body.String())
	})
	for _, r := range []struct {
		body string
		err  string
	}{
		{body: `{"topic": "topic", "endpoint": "lala", "type":"push"}`, err: "name is not present"},
		{body: `{"topic": "topic", "endpoint": "lala", "name":null, "type":"push"}`, err: "name is null"},
		{body: `{"name": 10, "endpoint": "http://www.ole.com.ar", "topic": "topic", "type":"push"}`, err: "name is numeric"},
		{body: `{"name": "", "endpoint": "http://www.ole.com.ar", "topic": "topic", "type":"push"}`, err: "name is empty"},
		{body: `{"name": "subs", "endpoint": "http://www.ole.com.ar", "topic": "topic", "type":"pull"}`, err: "endpoint is invalid for type pull"},
		{body: `{"name": "subscriber", "topic": "topic", "type":"push"}`, err: "endpoint is mandatory"},
		{body: `{"name": "subscriber", "endpoint": 8, "type":"push"}`},
		{body: `{"invalid": "topic", "invalid2": "lala"}`},
		{body: `{"name": "subscriber", "visibility_timeout":-1, "type":"pull"}`, err: "visibility_timeout out of range"},
		{body: `{"name": "subscriber", "visibility_timeout":50000, "type":"pull"}`, err: "visibility_timeout out of range"},
		{body: `{"topic": "topic"}`},
		{body: `{"name": "subscriber", "topic":"", "endpoint": "http://www.ole.com.ar", "type":"push"}`, err: "topic is empty"},
		{body: `{}`},
		{body: ``},
	} {
		t.Run("it should fail create the subscriber if the json fields are invalid ["+r.err+"]", func(t *testing.T) {

			res := executeMockedRequest(router, "POST", "/subscribers", r.body)
			assert.Contains(t, res.Body.String(), `"code":"validation_error"`)
			assert.Equal(t, 400, res.Code, res.Body.String())
		})
	}

	t.Run("it should fail creating the pull subscriber if an sqs.CreateQueue() error happends", func(t *testing.T) {
		mockSQS := &SQSAPIMock{}
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}

		subscriber := &model.Subscriber{Name: "subs",
			Topic:             "topic",
			Type:              "pull",
			ResourceID:        "arn:subs",
			PullingQueue:      "queueUrl",
			VisibilityTimeout: aws.Int(10),
			CreatedAt:         model.Clock.Now(),
		}
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == subscriber.Name && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockSQS.On("CreateQueue", &sqs.CreateQueueInput{QueueName: aws.String(client.GetAWSResourcePrefix() + "pull-queue-subs"),
			Attributes: map[string]*string{"VisibilityTimeout": aws.String("10")},
		}).Return(nil, errors.New("Create Queue error")).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "visibility_timeout":10, "type":"pull"}`)
		assert.Equal(t, 500, rec.Code)
		topicServiceMock.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		mockSQS.AssertExpectations(t)

	})

	t.Run("it should fail creating the push subscriber if an sqs.CreateQueue() error happends", func(t *testing.T) {
		mockSQS := &SQSAPIMock{}
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock
		client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}
		//The endpoint should be an active http endpoint
		client.HTTPClient = NewTestClient(func(req *http.Request) (*http.Response, error) {
			// Test request parameters
			assert.Equal(t, req.URL.String(), "http://endpoint")
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			}, nil
		})

		subscriber := getSubscriberMock("subs", "topic", "push", "arn:subs")
		//Topic should exist
		topicServiceMock.On("GetTopic", "topic").
			Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()
		//Subscriber with the provided name shold not exist in DB
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == subscriber.Name && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.ExpressionAttributeValues[":e"].S == "http://endpoint" &&
				*input.TableName == "Subscribers"
		})).Return(&dynamodb.ScanOutput{Items: nil}, nil).Once()

		mockSQS.On("CreateQueue", &sqs.CreateQueueInput{
			QueueName: aws.String(client.GetAWSResourcePrefix() + "dead-letter-subs"),
		}).Return(nil, errors.New("Create Queue error")).Once()

		rec := executeMockedRequest(router, "POST", "/subscribers", `{"topic": "topic", "name":"subs", "endpoint":"http://endpoint", "type":"push"}`)
		assert.Equal(t, 500, rec.Code)
		topicServiceMock.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		mockSQS.AssertExpectations(t)

	})
}

func TestConsumeQueueMessages(t *testing.T) {
	router := server.GetRouter()
	model.Clock = clockwork.NewFakeClock()

	for _, test := range []struct {
		subscriber *model.Subscriber
		queueType  string
		queueURL   string
		body       *string
	}{
		{
			subscriber: &model.Subscriber{Endpoint: aws.String("http://subscriber/endp"),
				Name: "subs", Topic: "topic", ResourceID: "arn:subs",
				DeadLetterQueue: "dlq_queue", Type: "push",
			},
			queueType: "dead letter queue",
			queueURL:  "dlq_queue",
			body:      aws.String(fmt.Sprintf(`{"Records": [{"Sns": {"Timestamp": "2018-11-22T14:02:21.284430Z", "Message": "{\"topic\":\"topic\",\"payload\":{\"hola\":\"lala\"},\"timestamp\":%d}", "Type": "Notification", "TopicArn": "arn:aws:sns:us-east-1:123456789012:test_topic", "Subject": null}}]}`, model.Clock.Now().UnixNano())),
		},
		{
			subscriber: &model.Subscriber{Name: "subs", Topic: "topic",
				ResourceID: "arn:subs", PullingQueue: "queue", Type: "pull",
			},
			queueType: "queue",
			queueURL:  "queue",
			body: aws.String(fmt.Sprintf(`{"Message":"{\"payload\":%s,\"timestamp\":%d,\"topic\":\"topic\"}","MessageId":"1","Type":"Notification","TopicArn":"arn:topic","SignatureVersion":"1"}`,
				`{\"hola\":\"lala\"}`, model.Clock.Now().UnixNano())),
		},
	} {
		t.Run("It should get messages from "+test.queueType, func(t *testing.T) {
			mockSQS := &SQSAPIMock{}
			mockDAO := &SubscriptionsDaoMock{}

			topicServiceMock := &TopicServiceMock{}
			service.TopicsService = topicServiceMock
			service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
			client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}
			mockDAO.On("GetSubscription", "subs").
				Return(test.subscriber, nil).Once()
			//Topic should exist
			topicServiceMock.On("GetTopic", "topic").
				Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()

			mockSQS.On("ReceiveMessage", &sqs.ReceiveMessageInput{
				MaxNumberOfMessages: aws.Int64(10), QueueUrl: &test.queueURL, WaitTimeSeconds: aws.Int64(0)}).
				Return(&sqs.ReceiveMessageOutput{
					Messages: []*sqs.Message{
						{Body: test.body,
							MessageId:     aws.String("1"),
							ReceiptHandle: aws.String("x")},
					},
				}, nil).Once()

			res := executeMockedRequest(router, "GET", "/messages?subscriber=subs&max_messages=10", "")
			assert.Equal(t, 200, res.Code)
			assert.JSONEq(t,
				fmt.Sprintf(`{"messages":[{"message":{"payload":{"hola":"lala"},"timestamp":%d,"topic":"topic"},
					"message_id":"1","delete_token":"x"}]}`,
					model.Clock.Now().UnixNano()),
				res.Body.String())

			mockSQS.AssertExpectations(t)
			mockDAO.AssertExpectations(t)
			topicServiceMock.AssertExpectations(t)

		})
	}
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
		subscriber := &model.Subscriber{Endpoint: aws.String("http://subscriber/endp"), Name: "subs", Topic: "topic",
			ResourceID: "arn:subs", DeadLetterQueue: "queue:subs"}
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

	for _, test := range []struct {
		subscriber *model.Subscriber
		queueType  string
		queueURL   string
	}{
		{subscriber: &model.Subscriber{Endpoint: aws.String("http://subscriber/endp"),
			Name: "subs", Topic: "topic", ResourceID: "arn:subs",
			DeadLetterQueue: "dlq_queue", Type: "push",
		},
			queueType: "dead letter queue",
			queueURL:  "dlq_queue"},

		{subscriber: &model.Subscriber{Name: "subs", Topic: "topic",
			ResourceID: "arn:subs", PullingQueue: "queue", Type: "pull",
		},
			queueType: "queue",
			queueURL:  "queue"},
	} {

		t.Run("It should delete specific messages from "+test.queueType, func(t *testing.T) {
			mockSQS := &SQSAPIMock{}
			mockDAO := &SubscriptionsDaoMock{}

			topicServiceMock := &TopicServiceMock{}
			service.TopicsService = topicServiceMock
			service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}
			client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS}
			mockDAO.On("GetSubscription", "subs").
				Return(test.subscriber, nil).Once()
			//Topic should exist
			topicServiceMock.On("GetTopic", "topic").
				Return(&model.Topic{ResourceID: "arn:topic", Name: "topic", Engine: "AWS"}, nil).Once()

			mockSQS.On("DeleteMessageBatch", &sqs.DeleteMessageBatchInput{
				Entries: []*sqs.DeleteMessageBatchRequestEntry{
					{Id: aws.String("1"), ReceiptHandle: aws.String("x")},
				}, QueueUrl: &test.queueURL}).
				Return(&sqs.DeleteMessageBatchOutput{
					Failed: make([]*sqs.BatchResultErrorEntry, 0),
				}, nil).Once()

			res := executeMockedRequest(router, "DELETE", "/messages", `{"subscriber":"subs", "messages": [{"message_id":"1", "delete_token":"x"}]}`)
			assert.Equal(t, 200, res.Code)
			assert.Equal(t, `{"failed":[]}`, res.Body.String())
			mockSQS.AssertExpectations(t)
			mockDAO.AssertExpectations(t)
			topicServiceMock.AssertExpectations(t)

		})

	}

	t.Run("It should fail deleting messages if a dao.GetSubscription() error happends", func(t *testing.T) {
		mockDAO := &SubscriptionsDaoMock{}
		mockDAO.On("GetSubscription", "subs").
			Return(nil, errors.New("Database error")).Once()
		service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: mockDAO}

		res := executeMockedRequest(router, "DELETE", "/messages", `{"subscriber":"subs", "messages": [{"message_id":"1", "delete_token":"x"}]}`)
		assert.Equal(t, 500, res.Code)

	})

}

func TestGetSubscription(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	router := server.GetRouter()
	subscriber := getSubscriberMock("subs", "topic", "push", "arn:subs")
	subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)
	mockDynamo := &DynamoDBAPIMock{}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: &model.SubscriberDaoDynamoImpl{
			DynamoClient: mockDynamo,
		},
	}
	t.Run("it should return the subscriber information", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == "subs" && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: subscriberItem}, nil).Once()

		res := executeMockedRequest(router, "GET", "/subscribers/subs", "")
		assert.Equal(t, 200, res.Code)
		assert.JSONEq(t, `{"name":"subs","type":"push","topic":"topic"}`, res.Body.String())
	})

	t.Run("it should return not found error if the subscriber doesn't exist", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == "subs" && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		res := executeMockedRequest(router, "GET", "/subscribers/subs", "")
		assert.Equal(t, 404, res.Code)
		assert.JSONEq(t, `{"message":"The subscriber subs doesn't exist","code":"not_found_error","status":404}`, res.Body.String())
	})

	t.Run("it should return an error if a dynamodb.GetItem() error happend", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == "subs" && *input.TableName == "Subscribers"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, errors.New("Dynamodb error")).Once()

		res := executeMockedRequest(router, "GET", "/subscribers/subs", "")
		assert.Equal(t, 500, res.Code)
		assert.JSONEq(t, `{"message":"Dynamodb error","code":"database_error","status":500}`, res.Body.String())
	})
}

func TestDeleteSubscriber(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	router := server.GetRouter()
	mockDynamo := &DynamoDBAPIMock{}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: &model.SubscriberDaoDynamoImpl{
			DynamoClient: mockDynamo,
		},
	}
	mockTopicService := &TopicServiceMock{}
	service.TopicsService = mockTopicService

	topicMock := getTopicMock("topic", "AWS", "arn:topic", "owner", "descr")

	for _, subs := range []model.Subscriber{
		model.Subscriber{Name: "subs", Type: "pull", ResourceID: "arn:subs", PullingQueue: "pullingQueueUrl", Topic: "topic"},
		model.Subscriber{Name: "subs", Type: "push", ResourceID: "arn:subs", DeadLetterQueue: "dlqQueueUrl", Topic: "topic"},
	} {
		t.Run("It should delete the "+subs.Type+" subscriber", func(t *testing.T) {
			mockSNS := &SNSAPIMock{}
			mockSQS := &SQSAPIMock{}
			mockLambda := &LambdaAPIMock{}
			client.EnginesMap["AWS"] = &client.AWSEngine{SQSClient: mockSQS, SNSClient: mockSNS, LambdaClient: mockLambda}

			subscriberItem, _ := dynamodbattribute.MarshalMap(subs)

			mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
				return *input.Key["name"].S == subs.Name && *input.TableName == "Subscribers"
			})).Return(&dynamodb.GetItemOutput{Item: subscriberItem}, nil).Once()

			mockTopicService.On("GetTopic", subs.Topic).Return(topicMock, nil).Once()

			mockSNS.On("Unsubscribe", &sns.UnsubscribeInput{SubscriptionArn: &subs.ResourceID}).
				Return(&sns.UnsubscribeOutput{}, nil).Once()

			mockSQS.On("DeleteQueue", &sqs.DeleteQueueInput{QueueUrl: aws.String(subs.GetQueueURL())}).
				Return(&sqs.DeleteQueueOutput{}, nil).Once()

			if subs.Type == model.SUBSCRIBER_PUSH {
				mockLambda.On("DeleteFunction", &lambda.DeleteFunctionInput{FunctionName: aws.String("wequeue-dev-lambda-" + subs.Name)}).
					Return(&lambda.DeleteFunctionOutput{}, nil).Once()
			}
			mockDynamo.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
				return *input.Key["name"].S == subs.Name
			})).Return(&dynamodb.DeleteItemOutput{}, nil).Once()

			res := executeMockedRequest(router, "DELETE", "/subscribers/subs", "")
			assert.Equal(t, 204, res.Code)
			mockTopicService.AssertExpectations(t)
			mockDynamo.AssertExpectations(t)
			mockSNS.AssertExpectations(t)
			mockSQS.AssertExpectations(t)
			mockLambda.AssertExpectations(t)
		})
	}
}

func TestMessage(t *testing.T) {

	msg := `{"Records":[{"EventSource":"aws:sns","EventVersion":"1.0",
	"EventSubscriptionArn":"arn:aws:sns:us-east-1:719849485599:wequeue-dev-sf_opp_stage_notif:f6b1a5b7-f2b1-451a-8840-197d93e547b1",
	"Sns":{"Type":"Notification","MessageId":"ab9c60de-9fe5-556e-877b-344f6a7eb51f",
	"TopicArn":"arn:aws:sns:us-east-1:719849485599:wequeue-dev-sf_opp_stage_notif",
	"Subject":null,"Message":"{\"topic\":\"sf_opp_stage_notif\",\"payload\":{\"idAccount\":\"0014100001SLgdLAAT\",\"idMambu\":\"254886\",\"idOpportunity\":\"0064100000aGLMzAAO\",\"stageName\":\"Entregado\"},\"timestamp\":1545073499329314304}",
	"Timestamp":"2018-12-17T19:04:59.340Z","SignatureVersion":"1",
	"Signature":"U4KIb2LC1CjB7UCN+ivfSKmXjAMs0AGvY1pr/LGT/iWIWtDzb6IEUfwtGM70z0vJTGbHR5LjPdV3I0E7H3bX7queGJQbteaAo/Q4mrD3z3AJEjBC65ggFUTbUC5WBbijbP50pLbFT01QYTR5TxPxDL/cs0DbKxYyJ+ZDBt+aZcP4TPOsqer0zkt3oHBYrEeuaZ7e3aU6ddjpt9x/X6kOhl0BnPD389lQnkdnUjB1/zyaKVTfxx5BHuK2JGd/UK1dxXvk3bNenexWTe0u59c9EI14BZy6Y0kXQ0iKU3IORT5Nn5MOzUtfhCUhBPPDB3PIjiRdxyywLsny/rSgnz6clg==",
	"SigningCertUrl":"https://sns.us-east-1.amazonaws.com/SimpleNotificationService-ac565b8b1a6c5d002d285f9598aa1d9b.pem",
	"UnsubscribeUrl":"https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:719849485599:wequeue-dev-sf_opp_stage_notif:f6b1a5b7-f2b1-451a-8840-197d93e547b1",
	"MessageAttributes":{}}}]}`
	buff := bytes.NewBufferString(msg)
	decoder := json.NewDecoder(buff)
	decoder.DisallowUnknownFields()
	var snsnotif client.SNSNotification
	err := decoder.Decode(&snsnotif)
	//Si es una dead-letter-queue de un push subscriber
	if err != nil {
		var dlqPayload client.DLQSNSNotification
		err := json.Unmarshal([]byte(msg), &dlqPayload)
		if err != nil {
			fmt.Printf("Error unmarshalling data %s, error: %s", msg, err.Error())
			t.Fail()
		}
		err = mapstructure.Decode(dlqPayload.Records[0]["Sns"], &snsnotif)
		if err != nil {
			fmt.Print(err.Error())
			t.Fail()
		}
	} else {
		t.Fail()
	}

	var publishedMessage model.PublishMessage
	err = json.Unmarshal([]byte(snsnotif.Message), &publishedMessage)
	if err != nil {
		fmt.Printf("Error unmarshalling payload %s", snsnotif.Message)
		t.Fail()
	}
}
