package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/wenance/wequeue-management_api/app/errors"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
)

type PublisherController struct {
}

// Publish godoc
// @Summary Publish a message in a topic
// @Description add by json a message in a topic. The 'X-Publish-Token: $security_token' header is mandatory.
// @Tags publishers
// @Accept json
// @Produce json
// @Param body body model.PublishMessage true "The message to publish"
// @Success 201 {object} model.PublishMessage
// @Failure 400 {object} errors.APIError
// @Failure 401 {object} errors.APIError "The X-Publish-Token header is invalid"
// @Failure 500 {object} errors.APIError
// @Router /messages [post]
func (t PublisherController) Publish(c *gin.Context) {
	securityToken := c.GetHeader("X-Publish-Token")
	var message model.PublishMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		if e, ok := err.(validation.Errors); ok {
			c.JSON(http.StatusBadRequest, errors.NewAPIError(http.StatusBadRequest, "validation_error", e.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, errors.NewAPIError(http.StatusBadRequest, "json_error", "The request body is not a valid json"))
		return
	}

	messageOutput, err := service.PublishersService.Publish(message, securityToken, c.Request.Header)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &messageOutput)
}
