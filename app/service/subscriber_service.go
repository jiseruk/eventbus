package service

import (
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
)

type SubscriptionService interface {
	CreateSubscription(name string, endpoint string, topic string) (*model.Subscriber, *app.APIError)
	ConsumeMessages(subscriber string, maxCount int64) (*model.Messages, *app.APIError)
}

type SubscriptionServiceImpl struct {
	Dao model.SubscriptionsDao
}

var SubscriptionsService SubscriptionService

func (s SubscriptionServiceImpl) CreateSubscription(name string, endpoint string, topic string) (*model.Subscriber, *app.APIError) {
	topicObj, apierr := TopicsService.GetTopic(topic)
	if apierr != nil {
		return nil, apierr
	}

	subscription, err := s.Dao.GetSubscription(name)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if subscription != nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Subscription with name %s already exists", name))
	}

	engine := client.GetEngineService(topicObj.Engine)
	output, err := engine.CreateSubscriber(*topicObj, name, endpoint)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "engine_error", err.Error())
	}

	if subscription, err := s.Dao.CreateSubscription(name, topic, endpoint, output.SubscriptionID, output.PullResourceID); err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_create_subscriber_error", err.Error())
	} else {
		return subscription, nil
	}

}

func (s SubscriptionServiceImpl) ConsumeMessages(subscriber string, maxCount int64) (*model.Messages, *app.APIError) {
	subscriberObj, err := s.Dao.GetSubscription(subscriber)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	topic, _ := TopicsService.GetTopic(subscriberObj.Topic)
	engine := client.GetEngineService(topic.Engine)
	messages, _ := engine.ReceiveMessages(subscriberObj.PullResourceID, maxCount)

	return messages, nil
}
