package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
)

type PublisherController struct {
}

func (t PublisherController) Publish(c *gin.Context) {
	var message model.PublishMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", "The request body is not a valid json"))
		return
	}

	if _, ok := message.Payload.(map[string]interface{}); !ok {
		if _, ok := message.Payload.([]interface{}); !ok {
			c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", "The payload should be a json"))
			return
		}
	}

	messageID, err := service.PublishersService.Publish(message.Topic, message.Payload)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	message.MessageID = *messageID
	c.JSON(http.StatusCreated, &message)
}
