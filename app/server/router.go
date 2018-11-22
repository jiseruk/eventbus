package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/controller"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	SetNewrelic(router)
	topics := controller.TopicController{}
	subscribers := controller.SubscriptionController{}
	publishers := controller.PublisherController{}
	health := controller.HealthController{}

	router.GET("/ping", health.Ping)
	router.POST("/topics", topics.Create)
	router.GET("/topics/:topic", topics.Get)
	router.POST("/subscribers", subscribers.Create)
	router.GET("/messages", subscribers.Consume)
	router.DELETE("/messages", subscribers.DeleteMessages)
	router.POST("/messages", publishers.Publish)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "page_not_found",
			"message": "Page not found",
			"status":  404})
	})

	return router

}
