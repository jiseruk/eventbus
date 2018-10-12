package service

import (
	"errors"
	"fmt"

	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
)

type TopicService interface {
	CreateTopic(name string, engine client.EngineService) error
}

type TopicServiceImpl struct {
	Db model.Topics
}

func (t TopicServiceImpl) CreateTopic(name string, engine client.EngineService) error {
	topic, err := t.Db.GetTopic(name)
	if err != nil {

	}
	if topic != nil {
		return errors.New(fmt.Sprintf("Topic with name %s already exists", name))
	}
	//engineService := client.GetEngineServiceImpl(engine)

	output, err := engine.CreateTopic(name)
	if err != nil {
		return err
	}
	if err = t.Db.CreateTopic(name, engine.GetName(), output.Resource); err != nil {

	}
	return nil
}
