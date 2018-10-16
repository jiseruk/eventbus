package server

import (
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"log"
)

func Init() {
	r := GetRouter()
	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient:client.GetSNSClient()}
	db, err := model.NewDB()
	if err != nil {
		log.Panic(err)
	}
	service.TopicsService = service.TopicServiceImpl{db}


	r.Run(":8080")
}