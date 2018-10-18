package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/controller"
	"github.com/wenance/wequeue-management_api/app/model"
	"net/http"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	topics := controller.TopicController{}
	subscribers := controller.SubscriptionController{}
	publishers := controller.PublisherController{}

	router.POST("/topics", topics.Create)
	router.POST("/subscriptions", subscribers.Create)
	router.POST("/messages", publishers.Publish)

	router.POST("/test_subscriber", func(c *gin.Context) {
		var json model.PublishMessage
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, &json)

	})

	return router

}