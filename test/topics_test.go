package test

import (
	"fmt"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
)


func TestCreateTopic(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	t.Run("It should create the topic in AWS and the DB entity", func(t *testing.T) {
		mockSNS := &SNSAPIMock{}
		mockDB := &TopicMock{}
		client.EnginesMap["AWS"] = &client.AWSEngine{mockSNS}
		service.TopicsService = service.TopicServiceImpl{mockDB}
		var topic = "topic"
		var resource = "arn:topic"
		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: &topic}).Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()
		mockDB.On("GetTopic", topic).Return(nil, nil).Once()
		mockDB.On("CreateTopic", topic, "AWS", resource).Return(nil).Once()
		router := server.GetRouter()
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, fmt.Sprintf(`{"name": "topic", "engine": "AWS", "resource_id":"arn:topic", "created_at": "%s"}`, model.Clock.Now().Format("2006-01-02T15:04:05Z")), rec.Body.String())

		mockSNS.AssertExpectations(t)
		mockDB.AssertExpectations(t)
	})

	t.Run("it should fail create the topic if it already exists", func(t *testing.T) {
		mockDB := &TopicMock{}
		service.TopicsService = service.TopicServiceImpl{mockDB}
		var topic = "topic"
		mockDB.On("GetTopic", topic).Return(&model.Topic{Name: topic}, nil).Once()
		router := server.GetRouter()
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/topics", strings.NewReader(`{"name": "topic", "engine": "AWS"}`))
		router.ServeHTTP(rec, req)
		assert.JSONEq(t, `{"status": 400, "code": "database_error", "message": "Topic with name topic already exists"}`, rec.Body.String())
		assert.Equal(t, 400, rec.Code)
		mockDB.AssertExpectations(t)

	})


}

/*
func TestCreateRealTopic(t *testing.T) {
	EngineService := client.GetEngineServiceImpl("AWS")
	output, _ := EngineService.CreateTopic("sarasa")
	fmt.Printf("%#v", output)
}
*/