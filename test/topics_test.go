package test

import (
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
	topicMock := getTopicMock(topic, "AWS", resource, "owner", "descr")
	topicItem, _ := dynamodbattribute.MarshalMap(topicMock)

	t.Run("It should create the topic in AWS and the DB entity", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.TopicsService = service.TopicServiceImpl{Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}}}
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
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS", "description":"descr", "owner":"owner"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t,
			fmt.Sprintf(`{"name": "topic", "engine": "AWS", "created_at": "%s", "security_token":"uuid", "owner":"owner", "description":"descr"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")), rec.Body.String())
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
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS", "description": "descr", "owner":"owner"}`))
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
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS", "description":"descr", "owner":"owner"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "database_error", "message": "Dynamodb error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		//mockDAO.AssertExpectations(t)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("it should fail create the topic if a dynamo.PutItem() error happend", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		//service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		service.TopicsService = service.TopicServiceImpl{
			Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}},
		}

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
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS", "description":"descr", "owner":"owner"}`))
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
		service.ADMIN_TOKEN_HASH = ""
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: aws.String(client.GetAWSResourcePrefix() + topic)}).
			Return(nil, errors.New("engine error")).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS", "description":"descr", "owner":"owner"}`))
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
	model.Clock = clockwork.NewFakeClock()
	mockDynamo := &DynamoDBAPIMock{}
	service.TopicsService = service.TopicServiceImpl{
		Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}},
	}
	//service.TopicsService = service.TopicServiceImpl{Dao: mockDAO}
	router := server.GetRouter()
	topic := getTopicMock("topic", "AWS", "arn:topic", "owner", "descr")
	topicItem, _ := dynamodbattribute.MarshalMap(topic)
	//	topicJSON, _ := json.Marshal(&topic)

	t.Run("It should return the topic", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic", "")
		assert.JSONEq(t,
			fmt.Sprintf(`{"name":"topic", "engine":"AWS", "created_at":"%s", "owner": "owner", "description":"descr"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")),
			res.Body.String())
		assert.Equal(t, 200, res.Code)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("It should return the topic with the security token if admin token header is provided", func(t *testing.T) {
		//Para no develar el verdadero Admin Token
		service.ADMIN_TOKEN_HASH = "d6d88aa97c88342259b82f64771658adcadcde3b18913b9c64e76129713c7605"
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")
		assert.JSONEq(t, fmt.Sprintf(`{"name":"topic", "engine":"AWS", "created_at":"%s", "security_token":"uuid", "owner":"owner", "description":"descr"}`,
			model.Clock.Now().Format("2006-01-02T15:04:05Z")), res.Body.String())
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

	t.Run("It should return not found error if the topic doesn't exist", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic", "")
		assert.JSONEq(t, res.Body.String(), `{"status":404,"code":"not_found_error","message":"The topic topic doesn't exist"}`)
		assert.Equal(t, 404, res.Code)
		mockDynamo.AssertExpectations(t)
	})
}

func TestListTopics(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockDynamo := &DynamoDBAPIMock{}
	service.TopicsService = service.TopicServiceImpl{
		Dao: &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}},
	}
	topic := &model.Topic{Name: "topic", Engine: "AWS", CreatedAt: model.Clock.Now()}
	topic2 := &model.Topic{Name: "topic2", Engine: "AWS", CreatedAt: model.Clock.Now()}
	topicItem, _ := dynamodbattribute.MarshalMap(topic)
	topicItem2, _ := dynamodbattribute.MarshalMap(topic2)
	router := server.GetRouter()

	t.Run("it should return the topics list", func(t *testing.T) {
		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.TableName == "Topics"
		})).Return(
			&dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{topicItem, topicItem2},
			}, nil,
		).Once()

		res := executeMockedRequest(router, "GET", "/topics", "")

		assert.Equal(t, 200, res.Code)
		assert.JSONEq(t,
			fmt.Sprintf(`{"topics": [{"name":"topic","engine":"AWS","created_at":"%s"}, 
						{"name":"topic2","engine":"AWS","created_at":"%s"}]}`,
				model.Clock.Now().Format("2006-01-02T15:04:05Z"), model.Clock.Now().Format("2006-01-02T15:04:05Z")),
			res.Body.String())
		mockDynamo.AssertExpectations(t)

	})

	t.Run("it fail returning the topics list if a dynamodb error happend", func(t *testing.T) {
		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.TableName == "Topics"
		})).Return(
			&dynamodb.ScanOutput{
				Items: nil,
			}, errors.New("Dynamodb error"),
		).Once()

		res := executeMockedRequest(router, "GET", "/topics", "")

		assert.Equal(t, 500, res.Code)
		assert.JSONEq(t, `{"status": 500, "code":"database_error", "message":"Dynamodb error"}`,
			res.Body.String())
		mockDynamo.AssertExpectations(t)

	})
}

func TestDeleteTopic(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockSNS := &SNSAPIMock{}
	mockDynamo := &DynamoDBAPIMock{}
	router := server.GetRouter()
	service.TopicsService = service.TopicServiceImpl{
		Dao:     &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}},
		SubsDao: &model.SubscriberDaoDynamoImpl{DynamoClient: mockDynamo},
	}
	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
	topic := getTopicMock("topic", "AWS", "arn:topic", "owner", "descr")
	subscriber := getSubscriberMock("subs", "topic", "push", "arn:subs")
	topicItem, _ := dynamodbattribute.MarshalMap(topic)
	subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)
	//Para no develar el verdadero Admin Token
	service.ADMIN_TOKEN_HASH = "d6d88aa97c88342259b82f64771658adcadcde3b18913b9c64e76129713c7605"

	t.Run("It should delete the topic and subscribers when admin token is provided", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		mockSNS.On("DeleteTopic", &sns.DeleteTopicInput{TopicArn: aws.String("arn:topic")}).
			Return(&sns.DeleteTopicOutput{}, nil).Once()

		mockDynamo.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.DeleteItemOutput{}, nil).Once()

		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.ExpressionAttributeValues[":t"].S == topic.Name &&
				*input.TableName == "Subscribers"
		})).Return(&dynamodb.ScanOutput{Items: []map[string]*dynamodb.AttributeValue{subscriberItem}}, nil).Once()

		mockDynamo.On("BatchWriteItem", mock.MatchedBy(func(input *dynamodb.BatchWriteItemInput) bool {
			return *input.RequestItems["Subscribers"][0].DeleteRequest.Key["name"].S == subscriber.Name
		})).Return(&dynamodb.BatchWriteItemOutput{}, nil).Once()

		res := executeMockedRequest(router, "DELETE", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")

		assert.Equal(t, http.StatusNoContent, res.Code)

		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)

	})

	t.Run("It should return Unauthorized error if no admin token is provided or is invalid", func(t *testing.T) {
		res := executeMockedRequest(router, "DELETE", "/topics/topic", "")
		assert.Equal(t, http.StatusUnauthorized, res.Code)

		res = executeMockedRequest(router, "DELETE", "/topics/topic", "X-Admin-Token:invalid")
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("It should return Not Found if the topic doesn't exist", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		res := executeMockedRequest(router, "DELETE", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")
		assert.Equal(t, http.StatusNotFound, res.Code)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("It should return an Error if a dynamodb.GetItem() error happends", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, errors.New("Dynamodb error")).Once()

		res := executeMockedRequest(router, "DELETE", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		mockDynamo.AssertExpectations(t)
	})

	t.Run("It should return an Error if a dynamodb.DeleteItem() error happends", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		mockSNS.On("DeleteTopic", &sns.DeleteTopicInput{TopicArn: aws.String("arn:topic")}).
			Return(&sns.DeleteTopicOutput{}, nil).Once()

		mockDynamo.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.DeleteItemOutput{}, errors.New("Dynamodb delete error")).Once()

		res := executeMockedRequest(router, "DELETE", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
	})

	t.Run("It should return an Error if a sns.DeleteTopic() error happends", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		mockSNS.On("DeleteTopic", &sns.DeleteTopicInput{TopicArn: aws.String("arn:topic")}).
			Return(&sns.DeleteTopicOutput{}, errors.New("Sns delete error")).Once()

		res := executeMockedRequest(router, "DELETE", "/topics/topic", "", "X-Admin-Token:PutoElQueDesencripta")
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		mockDynamo.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
	})
}

func TestListTopicSubscriptions(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockDynamo := &DynamoDBAPIMock{}
	router := server.GetRouter()
	service.TopicsService = service.TopicServiceImpl{
		Dao:     &model.TopicsDaoDynamoImpl{DynamoClient: mockDynamo, UUID: &UUIDMock{}},
		SubsDao: &model.SubscriberDaoDynamoImpl{DynamoClient: mockDynamo},
	}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: &model.SubscriberDaoDynamoImpl{DynamoClient: mockDynamo},
	}
	topic := getTopicMock("topic", "AWS", "arn:topic", "owner", "descr")
	topicItem, _ := dynamodbattribute.MarshalMap(topic)
	subscriber := getSubscriberMock("subs_push", "topic", "push", "arn:subs")
	subscriber2 := getSubscriberMock("subs_pull", "topic", "pull", "arn:subs")
	//topicItem, _ := dynamodbattribute.MarshalMap(topic)
	subscriberItem, _ := dynamodbattribute.MarshalMap(subscriber)
	subscriberItem2, _ := dynamodbattribute.MarshalMap(subscriber2)

	t.Run("it should return the topic subscribers", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: topicItem}, nil).Once()

		mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
			return *input.FilterExpression == "topic = :t" && *input.TableName == "Subscribers" &&
				*input.ExpressionAttributeValues[":t"].S == topic.Name
		})).Return(&dynamodb.ScanOutput{
			Items: []map[string]*dynamodb.AttributeValue{subscriberItem, subscriberItem2},
		}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic/subscribers", "")
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, `{"subscribers":[{"name":"subs_push","resource_id":"arn:subs","topic":"topic","type":"push"},{"name":"subs_pull","resource_id":"arn:subs","topic":"topic","type":"pull"}],"topic":"topic"}`, res.Body.String())
		mockDynamo.AssertExpectations(t)
	})

	t.Run("it should return bad request error if topic doesn't exist", func(t *testing.T) {
		mockDynamo.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.Key["name"].S == topic.Name && *input.TableName == "Topics"
		})).Return(&dynamodb.GetItemOutput{Item: nil}, nil).Once()

		res := executeMockedRequest(router, "GET", "/topics/topic/subscribers", "")
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"status": 404, "code":"not_found_error", "message":"The topic topic doesn't exist"}`,
			res.Body.String())
		mockDynamo.AssertExpectations(t)
	})
}
