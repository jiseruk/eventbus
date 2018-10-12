package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
	"net/http"
)

type TopicController struct{
	TopicService  service.TopicService
}

func (t TopicController) Create(c *gin.Context) {
	var json model.Topic
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	engine := client.GetEngineServiceImpl(json.Engine)
	t.TopicService.CreateTopic(json.Name, engine)
	c.JSON(http.StatusCreated, gin.H{"arn": "lala"})
}
