package client

import "github.com/wenance/wequeue-management_api/app/model"

type CreateTopicOutput struct {
	Resource string
	Shards   []*string
}
type PublishOutput struct {
	MessageID string
}
type SubscriberOutput struct {
	SubscriptionID string
	PullResourceID string
}

type EngineService interface {
	CreateTopic(name string) (*CreateTopicOutput, error)
	DeleteTopic(resource string) error
	Publish(topicResourceID string, message interface{}) (*PublishOutput, error)
	CreateSubscriber(topic model.Topic, subscriber string, endpoint string) (*SubscriberOutput, error)
	ReceiveMessages(resourceID string, maxMessages int64) ([]model.Message, error)
	DeleteMessages(messages []model.Message, queueUrl string) ([]model.Message, error)
	GetName() string
}

var EnginesMap = make(map[string]EngineService)

func GetEngineService(name string) EngineService {
	engine, ok := EnginesMap[name]
	if ok {
		return engine
	}
	return nil
}
