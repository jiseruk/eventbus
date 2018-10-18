package test

import "github.com/wenance/wequeue-management_api/app/model"

func getTopicMock(name string, engine string, resource string) *model.Topic{
	topic := model.Topic{ResourceID:resource, Name:name, Engine:engine}
	topic.CreatedAt = model.Clock.Now()
	topic.UpdatedAt = model.Clock.Now()
	return &topic
}
