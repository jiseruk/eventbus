package service

import (
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
)

type PublisherService interface {
	Publish(topic string, message interface{}) (*string, *app.APIError)
}

type PublisherServiceImpl struct {
}

var PublishersService PublisherService

func (PublisherServiceImpl) Publish(topic string, message interface{}) (*string, *app.APIError) {
	topicObj, apierr := TopicsService.GetTopic(topic)
	if apierr != nil {
		return nil, apierr
	}
	if topicObj == nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "topic_not_exists", fmt.Sprintf("The topic %s doesn't exist", topic))
	}
	engine := client.GetEngineService(topicObj.Engine)
	output, err := engine.Publish(topicObj.ResourceID, message)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "publish_error", err.Error())
	}
	return &output.MessageID, nil
}
