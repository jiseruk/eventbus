package service

import (
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/errors"
	"github.com/wenance/wequeue-management_api/app/model"
)

type SubscriptionService interface {
	CreateSubscription(name string, endpoint *string, topic string, Type string, visibilityTimeout *int) (*model.Subscriber, *errors.APIError)
	ConsumeMessages(subscriber string, maxCount int64) (*model.Messages, *errors.APIError)
	DeleteMessages(subscriber string, messages []model.Message) (*model.Messages, *errors.APIError)
	DeleteSubscription(name string) *errors.APIError
}

type SubscriptionServiceImpl struct {
	Dao        model.SubscriptionsDao
	HTTPClient *http.Client
}

var SubscriptionsService SubscriptionService

func (s SubscriptionServiceImpl) CreateSubscription(name string, endpoint *string, topic string, Type string, visibilityTimeout *int) (*model.Subscriber, *errors.APIError) {
	topicObj, apierr := TopicsService.GetTopic(topic)
	if apierr != nil {
		return nil, apierr
	}
	if topicObj == nil {
		return nil, errors.NewAPIError(http.StatusBadRequest, "validation_error", fmt.Sprintf("The topic %s doesn't exist", topic))
	}
	subscription, err := s.Dao.GetSubscription(name)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}

	if subscription != nil {
		return nil, errors.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Subscription with name %s already exists", name))
	}

	if ok, err := client.CheckEndpoint(endpoint); !ok {
		if err != nil {
			return nil, errors.NewAPIError(http.StatusBadRequest, "endpoint_error", fmt.Sprintf("The endpoint %s should return 2xx to a POST HTTP Call, but return error: %v", *endpoint, err.Error()))
		}
		return nil, errors.NewAPIError(http.StatusBadRequest, "endpoint_error", fmt.Sprintf("The endpoint %s should return 2xx to a POST HTTP Call", *endpoint))
	}

	engine := client.GetEngineService(topicObj.Engine)
	var output *client.SubscriberOutput
	if Type == model.SUBSCRIBER_PUSH {
		output, err = engine.CreatePushSubscriber(*topicObj, name, *endpoint)
	} else {
		output, err = engine.CreatePullSubscriber(*topicObj, name, *visibilityTimeout)

	}
	if err != nil {
		defer s.DeleteSubscription(name)
		return nil, errors.NewAPIError(http.StatusInternalServerError, "engine_error", err.Error())
	}

	if subscription, err := s.Dao.CreateSubscription(name, topic, Type, output.SubscriptionID, endpoint,
		output.DeadLetterQueue, output.PullingQueue, visibilityTimeout); err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_create_subscriber_error", err.Error())
	} else {
		return subscription, nil
	}

}

func (s SubscriptionServiceImpl) ConsumeMessages(subscriber string, maxCount int64) (*model.Messages, *errors.APIError) {
	subscriberObj, err := s.Dao.GetSubscription(subscriber)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if subscriberObj == nil {
		return nil, errors.NewAPIError(http.StatusNotFound, "database_error", "The subscriber "+subscriber+" doesn't exist")
	}

	topic, apierr := TopicsService.GetTopic(subscriberObj.Topic)
	if apierr != nil {
		return nil, apierr
	}
	if topic == nil {
		return nil, errors.NewAPIError(http.StatusNotFound, "database_error", "The topic "+subscriberObj.Topic+" doesn't exist")
	}

	engine := client.GetEngineService(topic.Engine)
	var messages []model.Message

	messages, err = engine.ReceiveMessages(subscriberObj.GetQueueURL(), maxCount)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "engine_error", err.Error())
	}
	return &model.Messages{Messages: messages}, nil
}

func (s SubscriptionServiceImpl) DeleteMessages(subscriber string, messages []model.Message) (*model.Messages, *errors.APIError) {
	//TODO: manage errors
	subscriberObj, err := s.Dao.GetSubscription(subscriber)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	topic, _ := TopicsService.GetTopic(subscriberObj.Topic)
	engine := client.GetEngineService(topic.Engine)
	result, _ := engine.DeleteMessages(messages, subscriberObj.GetQueueURL())

	return &model.Messages{Messages: result}, nil
}

func (s SubscriptionServiceImpl) DeleteSubscription(name string) *errors.APIError {
	return nil
}
