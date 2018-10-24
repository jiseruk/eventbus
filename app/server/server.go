package server

import (
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"github.com/wenance/wequeue-management_api/app/utils"
	"log"
)

func Init() {
	r := GetRouter()
	sns, lambda, kinesis := client.GetClients()

	utils.RecursiveZip("/go/src/github.com/wenance/wequeue-management_api/lambda/", "/tmp/function.zip")

	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient:sns, LambdaClient:lambda, KinesisClient:kinesis}
	db, err := model.NewDB()
	if err != nil {
		log.Panic(err)
	}
	service.TopicsService = service.TopicServiceImpl{db}
	service.SubscriptionsService = service.SubscriptionServiceImpl{Db:db}
	service.PublishersService = service.PublisherServiceImpl{}

	r.Run(":8080")
}