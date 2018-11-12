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

// Publish godoc
// @Summary Publish a message in a topic
// @Description add by json a message in a topic
// @Tags publishers
// @Accept json
// @Produce json
// @Param body body model.PublishMessage true "The message to publish"
// @Success 201 {object} model.PublishMessage
// @Failure 400 {object} app.APIError
// @Failure 500 {object} app.APIError
// @Router /messages [post]
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

	messageOutput, err := service.PublishersService.Publish(message)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &messageOutput)
}
