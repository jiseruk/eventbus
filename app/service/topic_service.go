package service

import (
	"fmt"
	"net/http"

	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/errors"
	"github.com/wenance/wequeue-management_api/app/model"
)

type TopicService interface {
	CreateTopic(name string, engine client.EngineService) (*model.Topic, *errors.APIError)
	GetTopic(name string) (*model.Topic, *errors.APIError)
}

type TopicServiceImpl struct {
	Dao model.TopicsDao
}

var TopicsService TopicService

func (t TopicServiceImpl) CreateTopic(name string, engine client.EngineService) (*model.Topic, *errors.APIError) {
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
	if topic, err := t.Dao.CreateTopic(name, engine.GetName(), output.Resource); err != nil {
		//TODO: Delete Topic in Engine

		delErr := engine.DeleteTopic(output.Resource)
		var multipleErrors string
		if delErr != nil {
			multipleErrors = fmt.Sprintf("%s | %s", err.Error(), delErr.Error())
		} else {
			multipleErrors = err.Error()
		}
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_create_topic_error", multipleErrors)
	} else {
		return topic, nil
	}
}

func (t TopicServiceImpl) GetTopic(name string) (*model.Topic, *errors.APIError) {
	topic, err := t.Dao.GetTopic(name)
	if err != nil {
		return nil, errors.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	return topic, nil
}
