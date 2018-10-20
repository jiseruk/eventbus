package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
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
)


func TestPublishMessage(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockSNS := &SNSAPIMock{}
	//mockDAO := &PublishersDaoMock{}
	topicServiceMock := &TopicServiceMock{}
	router := server.GetRouter()

	t.Run("It should publish a message in a topic", func(t *testing.T) {
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.PublishersService = service.PublisherServiceImpl{}
		service.TopicsService = topicServiceMock
		var topic= "topic"
		var resource= "arn:topic"
		var message= "message"
		topicMock := getTopicMock(topic, "AWS", resource)

		topicServiceMock.On("GetTopic", topic).Return(topicMock, nil).Once()
		mockSNS.On("Publish", sns.PublishInput{TopicArn: &resource, Message: &message}).Return(&sns.PublishOutput{MessageId: aws.String("1")})
		//mockDAO.On("GetTopic", topic).Return(nil, nil).Once()
		//mockDAO.On("CreateTopic", topic, "AWS", resource).Return(topicMock, nil).Once()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/messages", strings.NewReader(`{"topic": "topic", "payload": "message"}`))
		router.ServeHTTP(rec, req)

		assert.JSONEq(t, `{"message_id": "1"}`, rec.Body.String())
		assert.Equal(t, 201, rec.Code)
		mockSNS.AssertExpectations(t)
		//mockDAO.AssertExpectations(t)
	})
}