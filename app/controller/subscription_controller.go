package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
)

type SubscriptionController struct {
}

func (t SubscriptionController) Create(c *gin.Context) {
	var json model.Subscriber
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	subscriber, err := service.SubscriptionsService.CreateSubscription(json.Name, json.Endpoint, json.Topic)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &subscriber)
}

func (t SubscriptionController) Consume(c *gin.Context) {
	var consumeReq model.ConsumerRequest
	if err := c.ShouldBindQuery(&consumeReq); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	fmt.Print(consumeReq)
	messages, err := service.SubscriptionsService.ConsumeMessages(consumeReq.Subscriber, consumeReq.MaxMessages)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, &messages)
}
