package test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/server"
	"github.com/wenance/wequeue-management_api/app/service"
)

func TestPublishMessage(t *testing.T) {
	model.Clock = clockwork.NewFakeClock()
	mockSNS := &SNSAPIMock{}
	//mockDAO := &PublishersDaoMock{}
	topicServiceMock := &TopicServiceMock{}
	router := server.GetRouter()
	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
	service.PublishersService = service.PublisherServiceImpl{}
	service.TopicsService = topicServiceMock

	for _, payload := range []struct {
		payload string
	}{
		{payload: `{"message":"hola"}`},
		{payload: `["message",{"hola":"chau"}]`},
	} {
		t.Run("It should publish a message in a topic", func(t *testing.T) {
			var topic = "topic"
			var resource = "arn:topic"
			topicMock := getTopicMock(topic, "AWS", resource)

			topicServiceMock.On("GetTopic", topic).Return(topicMock, nil).Once()
			mockSNS.On("Publish", &sns.PublishInput{TopicArn: &resource, Message: &payload.payload}).Return(&sns.PublishOutput{MessageId: aws.String("1")}, nil)

			rec := executeMockedRequest(router, "POST", "/messages", fmt.Sprintf(`{"topic": "topic", "payload":%v}`, payload.payload))

			assert.JSONEq(t,
				fmt.Sprintf(`{"message_id": "1", "payload": %s, "topic":"topic"}`, payload.payload),
				rec.Body.String())
			assert.Equal(t, 201, rec.Code)
			mockSNS.AssertExpectations(t)
			topicServiceMock.AssertExpectations(t)
			//mockDAO.AssertExpectations(t)
		})
	}

	t.Run("It should fail publishing a message if payload is not a json", func(t *testing.T) {

		rec := executeMockedRequest(router, "POST", "/messages", `{"topic": "topic", "payload":"message"}`)

		assert.JSONEq(t,
			`{"message": "The payload should be a json", "status": 400, "code": "json_error"}`,
			rec.Body.String())
		assert.Equal(t, 400, rec.Code)

	})

	t.Run("It should fail publishing a message if the topic doesn't exist", func(t *testing.T) {
		topicServiceMock := &TopicServiceMock{}
		service.TopicsService = topicServiceMock

		var topic = "topic"
		var message = `{"message": "hola"}`

		topicServiceMock.On("GetTopic", topic).Return(nil, nil).Once()
		rec := executeMockedRequest(router, "POST", "/messages", fmt.Sprintf(`{"topic": "topic", "payload":%v}`, message))

		assert.JSONEq(t,
			`{"message": "The topic topic doesn't exist", "status": 400, "code": "topic_not_exists"}`,
			rec.Body.String())
		assert.Equal(t, 400, rec.Code)
		topicServiceMock.AssertExpectations(t)

	})

}
