package client

type CreateTopicOutput struct {
	Resource string
}

type EngineService interface {
	CreateTopic(name string) (*CreateTopicOutput, error)
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
