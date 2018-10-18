package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"net/http"
)

type SubscriptionController struct{
}

func (t SubscriptionController) Create(c *gin.Context) {
	var json model.Subscriber
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	topic, err := service.TopicsService.GetTopic(json.Topic)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	engine := client.GetEngineService(topic.Engine)
	subscriber, err := service.SubscriptionsService.CreateSubscription(json.Name, json.Endpoint, json.Topic, engine)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &subscriber)
}
