package client

type CreateTopicOutput struct {
	Resource string
}

type EngineService interface {
	CreateTopic(name string) (*CreateTopicOutput, error)
	GetName() string
}

func GetEngineServiceImpl(name string) EngineService {
	if name == "AWS" {
		return &AWSEngine{SNSClient: GetSNSClient()}
	}
	return nil
}
