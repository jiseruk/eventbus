package service

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/errors"
	"github.com/wenance/wequeue-management_api/app/model"
)

var ADMIN_TOKEN_HASH = "d571a372f7d39bcdcbb0c0e686c3b5118537d9e6266a999a83b578a9cf4ebdac"

type TopicService interface {
	CreateTopic(name string, owner string, description string, engine client.EngineService) (*model.Topic, *errors.APIError)
	GetTopic(name string, adminToken ...string) (*model.Topic, *errors.APIError)
	ListTopics() ([]model.Topic, *errors.APIError)
	DeleteTopic(name string, adminToken ...string) *errors.APIError
}

type TopicServiceImpl struct {
	Dao     model.TopicsDao
	SubsDao model.SubscriptionsDao
}

var TopicsService TopicService

func (t TopicServiceImpl) CreateTopic(name string, owner string, description string, engine client.EngineService) (*model.Topic, *errors.APIError) {
	topic, err := t.Dao.GetTopic(name)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if topic != nil {
		return nil, errors.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Topic with name %s already exists", name))
	}

	output, err := engine.CreateTopic(name)
	if err != nil {
		return nil, &errors.APIError{Status: http.StatusInternalServerError, Code: "engine_error", Message: err.Error()}
	}
	if topic, err := t.Dao.CreateTopic(name, engine.GetName(), owner, description, output.Resource); err != nil {
		//TODO: Delete Topic in Engine

		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_create_topic_error", err.Error())
	} else {
		topic.ResourceID = ""
		return topic, nil
	}
}

func (t TopicServiceImpl) GetTopic(name string, adminToken ...string) (*model.Topic, *errors.APIError) {

	topic, err := t.Dao.GetTopic(name)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if topic == nil {
		return nil, errors.NewAPIError(http.StatusNotFound, "not_found_error", "The topic "+name+" doesn't exist")
	}

	for _, token := range adminToken {
		if !IsValidAdminToken(token) {
			topic.SecurityToken = ""
			break
		}
	}
	return topic, nil
}

func (t TopicServiceImpl) ListTopics() ([]model.Topic, *errors.APIError) {
	topics, err := t.Dao.ListTopics()
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	return topics, nil
}

func (t TopicServiceImpl) DeleteTopic(name string, adminToken ...string) *errors.APIError {
	for _, token := range adminToken {
		if !IsValidAdminToken(token) {
			return errors.NewAPIError(http.StatusUnauthorized, "unauthorized_error", "The X-Admin-Token header should be provided")
		}
	}
	topic, apierr := t.GetTopic(name)
	if apierr != nil {
		return apierr
	}
	if topic == nil {
		return errors.NewAPIError(http.StatusNotFound, "not_found_error", "The topic "+name+" doesn't exist")
	}

	engine := client.GetEngineService(topic.Engine)
	err := engine.DeleteTopic(topic.ResourceID)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "engine_error", topic.ResourceID+err.Error())
	}

	err = t.Dao.DeleteTopic(name)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	//The sns subscribers are automatically deleted when the sns topic is deleted.
	err = t.SubsDao.DeleteTopicSubscriptions(name)
	if err != nil {
		return errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())

	}
	return nil
}

func IsValidAdminToken(token string) bool {
	if token == "" {
		return false
	}
	h := sha256.New()
	h.Write([]byte(token))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash == ADMIN_TOKEN_HASH
}
