package server

import (
	"log"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"github.com/wenance/wequeue-management_api/app/utils"
)

func Init() {
	r := GetRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sns, lambda, kinesis, sqs := client.GetClients()

	utils.RecursiveZip("/go/src/github.com/wenance/wequeue-management_api/lambda/", "/tmp/function.zip")

	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: sns, LambdaClient: lambda, SQSClient: sqs}
	client.EnginesMap["AWSStream"] = &client.AWSStreamEngine{LambdaClient: lambda, KinesisClient: kinesis}
	db, err := model.NewDB()
	if err != nil {
		log.Panic(err)
	}
	dynamo := model.GetClient()
	service.TopicsService = service.TopicServiceImpl{Db: db}
	//service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: &model.SubscriberDaoImpl{Db: *db}}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: &model.SubscriberDaoDynamoImpl{DynamoClient: dynamo},
	}
	service.PublishersService = service.PublisherServiceImpl{}

	r.Run(":8080")
}
