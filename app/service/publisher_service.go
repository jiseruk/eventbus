package service

import (
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
)

type PublisherService interface {
	Publish(message model.PublishMessage) (*model.PublishMessage, *app.APIError)
}

type PublisherServiceImpl struct {
}

var PublishersService PublisherService

func (PublisherServiceImpl) Publish(message model.PublishMessage) (*model.PublishMessage, *app.APIError) {
	topicObj, apierr := TopicsService.GetTopic(message.Topic)
	if apierr != nil {
		return nil, apierr
	}
	if topicObj == nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "topic_not_exists", fmt.Sprintf("The topic %s doesn't exist", message.Topic))
	}
	engine := client.GetEngineService(topicObj.Engine)
	timestamp := model.Clock.Now().UnixNano()
	message.Timestamp = &timestamp
	_, err := engine.Publish(topicObj.ResourceID, &message)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "publish_error", err.Error())
	}
	return &message, nil
}
