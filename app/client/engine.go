package client

import "github.com/aws/aws-sdk-go/service/lambda"

type CreateTopicOutput struct {
	Resource string
}
type PublishOutput struct{
	MessageID string
}
type SubscriberOutput struct {
	SubscriptionID string
}

type EngineService interface {
	CreateTopic(name string) (*CreateTopicOutput, error)
	DeleteTopic(resource string) error
	Publish(topicResourceID string, message interface{}) (*PublishOutput, error)
	CreateSubscriber(topicResourceID string, subscriber string, endpoint string) (*SubscriberOutput, error)
	GetName() string
	InvokeLambda(name string, payload string) (*lambda.InvokeOutput, error)
}

var EnginesMap = make(map[string]EngineService)


func GetEngineService(name string) EngineService {
	engine, ok := EnginesMap[name]
	if ok {
		return engine
	}
	return nil
}
