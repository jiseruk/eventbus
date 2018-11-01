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

// Create godoc
// @Summary Add a subscriber
// @Description add by json a subscriber to a topic
// @Tags subscribers
// @Accept json
// @Produce json
// @Param body body model.Subscriber true "Subscriber to a topic"
// @Success 201 {object} model.Subscriber
// @Failure 400 {object} app.APIError
// @Failure 500 {object} app.APIError
// @Router /subscribers [post]
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

// Consume godoc
// @Summary Consume pending messages
// @Description consume pending messages from the subscriber's dead letter queue
// @Tags subscribers
// @Accept json
// @Produce json
// @Param subscriber query string true "The Subscriber name"
// @Param max_messages query number true "Max messages to get"
// @Success 200 {object} model.Messages
// @Failure 400 {object} app.APIError
// @Failure 500 {object} app.APIError
// @Router /messages [get]
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

func (t SubscriptionController) DeleteMessages(c *gin.Context) {
	var deleteReq model.DeleteDeadLetterQueueMessagesRequest
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	messages, err := service.SubscriptionsService.DeleteMessages(deleteReq.Subscriber, deleteReq.Messages)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, &model.DeleteDeadLetterQueueMessagesResponse{Failed: messages})
}
