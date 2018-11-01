package test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	_ "github.com/wenance/wequeue-management_api/app/validation"
)

func TestCreateTopic(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockSNS := &SNSAPIMock{}
	mockDAO := &TopicsDaoMock{}
	router := server.GetRouter()

	t.Run("It should create the topic in AWS and the DB entity", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		var topic = "topic"
		var resource = "arn:topic"
		topicMock := getTopicMock(topic, "AWS", resource)

		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: &topic}).Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()
		mockDAO.On("GetTopic", topic).Return(nil, nil).Once()
		mockDAO.On("CreateTopic", topic, "AWS", resource).Return(topicMock, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t, fmt.Sprintf(`{"name": "topic", "engine": "AWS", "resource_id":"arn:topic", "created_at": "%s"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")), rec.Body.String())
		assert.Equal(t, 201, rec.Code)
		mockSNS.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
	})

	t.Run("it should fail create the topic if it already exists", func(t *testing.T) {
		service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		var topic = "topic"
		mockDAO.On("GetTopic", topic).Return(&model.Topic{Name: topic}, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t, `{"status": 400, "code": "database_error", "message": "Topic with name topic already exists"}`, rec.Body.String())
		assert.Equal(t, 400, rec.Code)
		mockDAO.AssertExpectations(t)

	})

	t.Run("it should fail create the topic if a dao.GetTopic() error happend", func(t *testing.T) {
		service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		var topic = "topic"
		mockDAO.On("GetTopic", topic).Return(nil, errors.New("database error")).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "database_error", "message": "database error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		mockDAO.AssertExpectations(t)

	})

	t.Run("it should fail create the topic if a dao.CreateTopic() error happend", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		var topic = "topic"
		var resource = "arn:topic"

		mockDAO.On("GetTopic", topic).Return(nil, nil).Once()
		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: &topic}).Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()
		mockDAO.On("CreateTopic", topic, "AWS", resource).Return(nil, errors.New("database error")).Once()
		mockSNS.On("DeleteTopic", &sns.DeleteTopicInput{TopicArn: &resource}).Return(&sns.DeleteTopicOutput{}, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "database_create_topic_error", "message": "database error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		mockDAO.AssertExpectations(t)
		mockSNS.AssertExpectations(t)
	})

	t.Run("It should fail creating the topic if an engine.CreateTopic() error happends", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.TopicsService = service.TopicServiceImpl{Db: mockDAO}
		var topic = "topic"

		mockDAO.On("GetTopic", topic).Return(nil, nil).Once()
		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: &topic}).Return(nil, errors.New("engine error")).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 500, "code": "engine_error", "message": "engine error"}`, rec.Body.String())
		assert.Equal(t, 500, rec.Code)

		mockSNS.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
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
			assert.Contains(t, res.Body.String(), `"code":"json_error"`)
			assert.Equal(t, 400, res.Code, res.Body.String())
		})
	}

}
