package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	_ "github.com/wenance/wequeue-management_api/app/validation"
)

func TestCreateTopic(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockSNS := &SNSAPIMock{}
	//mockDAO := &TopicsDaoMock{}
	mockDynamo := &DynamoDBAPIMock{}
	router := server.GetRouter()
	var topic = "topic"
	var resource = "arn:topic"
	topicMock := getTopicMock(topic, "AWS", resource)
	topicItem, _ := dynamodbattribute.MarshalMap(topicMock)

	t.Run("It should create the topic in AWS and the DB entity", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}
		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: aws.String(client.GetAWSResourcePrefix() + topic)}).
			Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()

		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockDynamo.On("PutItem", &dynamodb.PutItemInput{
			Item:      topicItem,
			TableName: aws.String("Topics"),
		}).Return(&dynamodb.PutItemOutput{}, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t, fmt.Sprintf(`{"name": "topic", "engine": "AWS", "resource_id":"arn:topic", "created_at": "%s"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")), rec.Body.String())
		assert.Equal(t, 201, rec.Code)
		mockSNS.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		//mockDAO.AssertExpectations(t)
	})

	t.Run("it should fail create the topic if it already exists", func(t *testing.T) {
		//service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}
		var topic = "topic"
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{},
		}, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t, `{"status": 400, "code": "database_error", "message": "Topic with name topic already exists"}`, rec.Body.String())
		assert.Equal(t, 400, rec.Code)
		mockDynamo.AssertExpectations(t)

	})

	t.Run("it should fail create the topic if a dynamodb.GetItem() error happends", func(t *testing.T) {
		//service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}
		mockDynamo.On("GetItem", mock.Anything).Return(nil, errors.New("Dynamodb error")).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "database_error", "message": "Dynamodb error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		//mockDAO.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("it should fail create the topic if a dynamo.PutItem() error happend", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		//service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}

		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockDynamo.On("PutItem", &dynamodb.PutItemInput{
			Item:      topicItem,
			TableName: aws.String("Topics"),
		}).Return(nil, errors.New("Dynamodb error")).Once()

		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: aws.String(client.GetAWSResourcePrefix() + topic)}).
			Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "database_create_topic_error", "message": "Dynamodb error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		//mockDAO.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
	})

	t.Run("It should fail creating the topic if an engine.CreateTopic() error happends", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		//service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}

		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: aws.String(client.GetAWSResourcePrefix() + topic)}).
			Return(nil, errors.New("engine error")).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "engine_error", "message": "engine error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)

		mockSNS.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
	})

	for _, r := range []struct {
		body string
	}{
		{body: `{"invalid": "topic", "engine": "AWS"}`},
		{body: `{"name": "topic", "invalid": "AWS"}`},
		{body: `{"name": 5, "engine": "AWS"}`},
		{body: `{"name": "topic", "engine": 8}`},
		{body: `{"invalid": "topic", "invalid2": "AWS"}`},
		{body: `{"engine": "AWS"}`},
		{body: `{"name": "topic"}`},
		{body: `{"name": "topic", "engine": "invalid"}`},
		{body: `{}`},
		{body: ``},
	} {
		t.Run("it should fail create the topic if the json fields are invalid", func(t *testing.T) {

			res := executeMockedRequest(router, "POST", "/topics", r.body)
			assert.Contains(t, res.Body.String(), `"code":"validation_error"`)
			assert.Equal(t, 400, res.Code, res.Body.String())
		})
	}

}

func TestGetTopic(t *testing.T) {
	//mockDAO := &TopicsDaoMock{}
	mockDynamo := &DynamoDBAPIMock{}
	service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo}}
	//service.TopicsService = service.TopicServiceImpl{Dao: mockDAO}
	router := server.GetRouter()
	topic := &model.Topic{Name: "topic", Engine: "AWS"}
	topicItem, _ := dynamodbattribute.MarshalMap(topic)
	topicJSON, _ := json.Marshal(&topic)

	t.Run("It should return the topic", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic", "")
		assert.JSONEq(t, res.Body.String(), string(topicJSON))
		assert.Equal(t, 200, res.Code)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("It should fail returning the topic when a database error happend", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, errors.New("Dynamo DB error")).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic", "")
		assert.JSONEq(t, res.Body.String(), `{"status":500,"code":"database_error","message":"Dynamo DB error"}`)
		assert.Equal(t, 500, res.Code)
		mockDynamo.AssertExpectations(t)
	})

	for _, path := range []string{"/topics/", "/topics"} {

		t.Run("It should fail getting topic if the topic param is not present", func(t *testing.T) {

			res := executeMockedRequest(router, "GET", path, "")
			assert.JSONEq(t, res.Body.String(), `{"message":"Page not found", "code":"page_not_found", "status":404}`)
			assert.Equal(t, 404, res.Code)

		})
	}

}
