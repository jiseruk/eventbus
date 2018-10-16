package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/controller"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	topics := controller.TopicController{}

	router.POST("/topics", topics.Create)


	return router

}