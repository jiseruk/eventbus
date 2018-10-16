package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"net/http"
)

type TopicController struct{
}

var EngineService client.EngineService

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
	c.JSON(http.StatusCreated, topic)
}
