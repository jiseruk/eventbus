package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/controller"
	"github.com/wenance/wequeue-management_api/app/errors"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	router.POST("/test_subscriber", func(c *gin.Context) {
		var message map[string]interface{}
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, errors.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
			return
		}
		log.Printf("Message received: %#v", message)

		/*payload := message.Payload.(map[string]interface{})
		log.Printf("Message received, payload: %#v", payload)
		if payload != nil && payload["fail"] == true {
			c.JSON(http.StatusInternalServerError, &message)
			return
		}*/
		c.JSON(http.StatusOK, &message)

	})
	/*
		router.POST("/lambda", func(c *gin.Context) {
			var json map[string]string
			if err := c.ShouldBindJSON(&json); err != nil {
				c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
				return
			}
			engine := client.GetEngineService("AWS").(client.AWSEngine)
			output, err := engine.InvokeLambda(json["name"], json["payload"])
			if err != nil {
				c.JSON(500, err.Error())
				return
			}
			c.JSON(200, &output)

		})
	*/
	return router

}
