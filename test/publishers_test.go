package test

import (
	"errors"
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
	_ "github.com/wenance/wequeue-management_api/app/validation"
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
		{payload: `[{"message":"hola"}]`},
		{payload: `{"message":{"hola":"chau"}}`},
	} {
		t.Run("It should publish a message in a topic", func(t *testing.T) {
			var topic = "topic"
			var resource = "arn:topic"
			topicMock := getTopicMock(topic, "AWS", resource)

			topicServiceMock.On("GetTopic", topic).Return(topicMock, nil).Once()
			mockSNS.On("Publish", &sns.PublishInput{TopicArn: &resource,
				Message: aws.String(fmt.Sprintf(
					`{"topic":"topic","payload":%s,"timestamp":%d}`,
					payload.payload,
					model.Clock.Now().UnixNano()),
				),
			}).Return(&sns.PublishOutput{MessageId: aws.String("1")}, nil)

			rec := executeMockedRequest(router, "POST", "/messages", fmt.Sprintf(`{"topic": "topic", "payload":%v}`, payload.payload))

			assert.JSONEq(t,
				fmt.Sprintf(`{"timestamp": %d, "payload": %s, "topic":"topic"}`, model.Clock.Now().UnixNano(), payload.payload),
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
			`{"message": "payload: it should be a valid json object.", "status": 400, "code": "json_error"}`,
			rec.Body.String())
		assert.Equal(t, 400, rec.Code)

	})

	for _, body := range []string{
		``,
		`lala`,
		`{{}`,
	} {
		t.Run("It should fail publishing a message if the json is not valid", func(t *testing.T) {

			rec := executeMockedRequest(router, "POST", "/messages", body)
			assert.JSONEq(t,
				`{"message": "The request body is not a valid json", "status": 400, "code": "json_error"}`,
				rec.Body.String())
			assert.Equal(t, 400, rec.Code)

		})
	}

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

	t.Run("It should fail publishing a message if a sns.Publish() error happends", func(t *testing.T) {
		mockSNS := &SNSAPIMock{}
		//mockDAO := &PublishersDaoMock{}
		topicServiceMock := &TopicServiceMock{}
		router := server.GetRouter()
		client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: mockSNS}
		service.PublishersService = service.PublisherServiceImpl{}
		service.TopicsService = topicServiceMock

		topicMock := getTopicMock("topic", "AWS", "arn:topic")
		var message = `{"message":"hola"}`

		topicServiceMock.On("GetTopic", topicMock.Name).Return(topicMock, nil).Once()

		mockSNS.On("Publish", &sns.PublishInput{TopicArn: &topicMock.ResourceID,
			Message: aws.String(fmt.Sprintf(
				`{"topic":"topic","payload":%s,"timestamp":%d}`,
				message,
				model.Clock.Now().UnixNano()),
			),
		}).Return(nil, errors.New("Publish error"))

		rec := executeMockedRequest(router, "POST", "/messages", fmt.Sprintf(`{"topic": "topic", "payload":%v}`, message))

		assert.JSONEq(t,
			`{"message": "Publish error", "status": 500, "code": "engine_error"}`,
			rec.Body.String())
		assert.Equal(t, 500, rec.Code)
		topicServiceMock.AssertExpectations(t)

	})

}
