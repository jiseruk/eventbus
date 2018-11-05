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
	DeleteMessages(subscriber string, messages []model.Message) (*model.Messages, *app.APIError)
}

type SubscriptionServiceImpl struct {
	Dao        model.SubscriptionsDao
	HTTPClient *http.Client
}

var SubscriptionsService SubscriptionService

func (s SubscriptionServiceImpl) CreateSubscription(name string, endpoint string, topic string) (*model.Subscriber, *app.APIError) {
	topicObj, apierr := TopicsService.GetTopic(topic)
	if apierr != nil {
		return nil, apierr
	}
	if topicObj == nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "validation_error", fmt.Sprintf("The topic %s doesn't exist", topic))
	}
	subscription, err := s.Dao.GetSubscription(name)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}

	if subscription != nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Subscription with name %s already exists", name))
	}

	if ok, err := client.CheckEndpoint(endpoint); !ok || err != nil {
		if err != nil {
			return nil, app.NewAPIError(http.StatusBadRequest, "endpoint_error", fmt.Sprintf("The endpoint %s should return 2xx to a POST HTTP Call, but return error: %v", endpoint, err.Error()))
		}
		return nil, app.NewAPIError(http.StatusBadRequest, "endpoint_error", fmt.Sprintf("The endpoint %s should return 2xx to a POST HTTP Call", endpoint))
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
	messages, _ := engine.ReceiveMessages(subscriberObj.DeadLetterQueue, maxCount)

	return &model.Messages{Messages: messages, Topic: topic.Name}, nil
}

func (s SubscriptionServiceImpl) DeleteMessages(subscriber string, messages []model.Message) (*model.Messages, *app.APIError) {
	//TODO: manage errors
	subscriberObj, err := s.Dao.GetSubscription(subscriber)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	topic, _ := TopicsService.GetTopic(subscriberObj.Topic)
	engine := client.GetEngineService(topic.Engine)
	result, _ := engine.DeleteMessages(messages, subscriberObj.DeadLetterQueue)

	return &model.Messages{Messages: result, Topic: topic.Name}, nil
}
