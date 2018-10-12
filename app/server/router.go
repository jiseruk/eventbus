package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/controller"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"log"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	db, err := model.NewDB()
	if err != nil {
		log.Panic(err)
	}

	topics := controller.TopicController{TopicService: service.TopicServiceImpl{Db:db}}

	router.POST("/topics", topics.Create)


	return router

}