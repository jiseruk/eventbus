package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"net/http"
)

type PublisherController struct{
}

func (t PublisherController) Publish(c *gin.Context) {
	var json model.PublishMessage
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	messageId, err := service.PublishersService.Publish(json.Topic, json.Payload)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	json.MessageID = *messageId
	c.JSON(http.StatusCreated, &json)
}

