package test

import (
	"github.com/wenance/wequeue-management_api/app/client"
	"testing"

	"github.com/wenance/wequeue-management_api/app/service"

	"github.com/aws/aws-sdk-go/service/sns"
)

func TestCreateTopic(t *testing.T) {
	t.Run("It should create the topic in AWS and the DB entity", func(t *testing.T) {
		mockSNS := &SNSAPIMock{}
		mockDB := &TopicMock{}
		//client.Engine = &client.AWSEngine{SNSClient: mockSNS}
		topicService := service.TopicServiceImpl{mockDB}
		var topic = "topic"
		var resource = "arn:topic"
		mockSNS.On("CreateTopic", &sns.CreateTopicInput{Name: &topic}).Return(&sns.CreateTopicOutput{TopicArn: &resource}, nil).Once()
		mockDB.On("GetTopic", topic).Return(nil, nil).Once()
		mockDB.On("CreateTopic", topic, "AWS", resource).Return(nil).Once()

		topicService.CreateTopic(topic, &client.AWSEngine{SNSClient: mockSNS})
		mockSNS.AssertExpectations(t)
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