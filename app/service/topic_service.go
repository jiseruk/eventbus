package service

import (
	"fmt"
	"github.com/wenance/wequeue-management_api/app"
	"net/http"

	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
)

type TopicService interface {
	CreateTopic(name string, engine client.EngineService) (*model.Topic, *app.APIError)
	GetTopic(name string) (*model.Topic, *app.APIError)
}

type TopicServiceImpl struct {
	Db model.TopicsDao
}

var TopicsService TopicService

func (t TopicServiceImpl) CreateTopic(name string, engine client.EngineService) (*model.Topic, *app.APIError) {
	topic, err := t.Db.GetTopic(name)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	if topic != nil {
		return nil, app.NewAPIError(http.StatusBadRequest, "database_error", fmt.Sprintf("Topic with name %s already exists", name))
	}

	output, err := engine.CreateTopic(name)
	if err != nil {
		return nil, &app.APIError{Status: http.StatusInternalServerError, Code: "engine_error", Message:err.Error()}
	}
	if topic, err := t.Db.CreateTopic(name, engine.GetName(), output.Resource); err != nil {
		//TODO: Delete Topic in Engine
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_create_topic_error", err.Error())
	} else {
		return topic, nil
	}
}

func (t TopicServiceImpl) GetTopic(name string)  (*model.Topic, *app.APIError){
	topic, err := t.Db.GetTopic(name)
	if err != nil {
		return nil, app.NewAPIError(http.StatusInternalServerError, "database_error", err.Error())
	}
	return topic, nil
}