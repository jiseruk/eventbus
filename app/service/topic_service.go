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
}

type TopicServiceImpl struct {
	Dao model.TopicsDao
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
	for i, token := range adminToken {
		if token == "" || i > 0 {
			topic.SecurityToken = ""
			break
		}

		h := sha256.New()
		h.Write([]byte(token))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		if hash != ADMIN_TOKEN_HASH {
			topic.SecurityToken = ""
			break
		}
	}
	topic.ResourceID = ""
	return topic, nil
}

func (t TopicServiceImpl) ListTopics() ([]model.Topic, *errors.APIError) {
	topics, err := t.Dao.ListTopics()
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	return topics, nil
}
