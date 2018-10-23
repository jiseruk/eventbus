package service

import (
	"fmt"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"net/http"
)

type SubscriptionService interface {
	CreateSubscription(name string, endpoint string, topic string) (*model.Subscriber, *app.APIError)
}

type SubscriptionServiceImpl struct {
	Db model.SubscriptionsDao
}

var SubscriptionsService SubscriptionService

func (s SubscriptionServiceImpl) CreateSubscription(name string, endpoint string, topic string) (*model.Subscriber, *app.APIError) {
	topicObj, apierr := TopicsService.GetTopic(topic)
	if apierr != nil {
		return nil, apierr
	}
	subscription, err := s.Db.GetSubscription(name)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if subscription != nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Subscription with name %s already exists", name))
	}
	engine := client.GetEngineService(topicObj.Engine)
	output, err := engine.CreateSubscriber(topicObj.ResourceID, name, endpoint)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "engine_error", err.Error())
	}

	if subscription, err := s.Db.CreateSubscription(name, topicObj.ResourceID, endpoint, output.SubscriptionID); err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_create_subscriber_error", err.Error())
	} else {
		return subscription, nil
	}

	return nil, nil
}
