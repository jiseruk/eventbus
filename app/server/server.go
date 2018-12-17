package server

import (
	"fmt"

	"github.com/jonboulle/clockwork"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/config"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"github.com/wenance/wequeue-management_api/app/utils"
)

func Init() {
	model.Clock = clockwork.NewRealClock()

	r := GetRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sns, lambda, kinesis, sqs := client.GetClients()

	utils.RecursiveZip(config.Get("engines.AWS.lambda.zipDir"), "/tmp/function.zip")

	client.EnginesMap["AWS"] = &client.AWSEngine{SNSClient: sns, LambdaClient: lambda, SQSClient: sqs}
	client.EnginesMap["AWSStream"] = &client.AWSStreamEngine{LambdaClient: lambda, KinesisClient: kinesis}
	/*db, err := model.NewDB()
	if err != nil {
		log.Panic(err)
	}*/
	dynamo := model.GetClient()
	subscribersDao := &model.SubscriberDaoDynamoImpl{
		DynamoClient: dynamo,
	}
	//service.TopicsService = service.TopicServiceImpl{Db: db}
	service.TopicsService = service.TopicServiceImpl{
		Dao: &model.TopicsDaoDynamoImpl{
			DynamoClient: dynamo,
			UUID:         model.UUIDImpl{},
		},
		SubsDao: subscribersDao,
	}
	//service.SubscriptionsService = service.SubscriptionServiceImpl{Dao: &model.SubscriberDaoImpl{Db: *db}}
	service.SubscriptionsService = service.SubscriptionServiceImpl{
		Dao: subscribersDao,
	}
	service.PublishersService = service.PublisherServiceImpl{}
	r.Run(":8080")
}
