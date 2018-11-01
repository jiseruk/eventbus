package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
)

type TopicController struct {
}

var EngineService client.EngineService

// Create godoc
// @Summary Add a topic
// @Description add by json a topic
// @Tags topics
// @Accept json
// @Produce json
// @Param body body model.Topic true "Topic created for publishing messages"
// @Success 201 {object} model.Topic
// @Failure 400 {object} app.APIError
// @Failure 500 {object} app.APIError
// @Router /topics [post]
func (t TopicController) Create(c *gin.Context) {
	var json model.Topic
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, app.NewAPIError(http.StatusBadRequest, "json_error", err.Error()))
		return
	}
	engine := client.GetEngineService(json.Engine)
	topic, err := service.TopicsService.CreateTopic(json.Name, engine)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &topic)
}
