package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenance/wequeue-management_api/app/client"
	"github.com/wenance/wequeue-management_api/app/errors"
	"github.com/wenance/wequeue-management_api/app/model"
	"github.com/wenance/wequeue-management_api/app/service"
)

type TopicController struct {
}

var EngineService client.EngineService

// Create godoc
// @Summary Add a topic
// @Description It creates a new topic in Bondi. The security_token field should be saved for publishing.
// @Tags topics
// @Accept json
// @Produce json
// @Param body body model.Topic true "Topic created for publishing messages"
// @Success 201 {object} model.Topic
// @Failure 400 {object} errors.APIError
// @Failure 500 {object} errors.APIError
// @Router /topics [post]
func (t TopicController) Create(c *gin.Context) {
	var json model.Topic
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(http.StatusBadRequest, "validation_error", err.Error()))
		return
	}
	engine := client.GetEngineService(json.Engine)
	topic, err := service.TopicsService.CreateTopic(json.Name, json.Owner, json.Description, engine)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, &topic)
}

// Get godoc
// @Summary Get a topic
// @Description Get the topic information
// @Tags topics
// @Accept json
// @Produce json
// @Param topic path string true "The name of the Topic"
// @Success 200 {object} model.Topic
// @Failure 404 {object} errors.APIError "The topic doesn't exist"
// @Failure 500 {object} errors.APIError
// @Router /topics/{topic} [get]
// @OperationId get-topic
func (t TopicController) Get(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	topicName := c.Param("topic")

	topic, apierr := service.TopicsService.GetTopic(topicName, adminToken)
	if apierr != nil {
		c.JSON(apierr.Status, &apierr)
		return
	}
	topic.ResourceID = ""
	c.JSON(http.StatusOK, &topic)
}

// List godoc
// @Summary List topics
// @Description List all the topics information
// @Tags topics
// @Accept json
// @Produce json
// @Success 200 {array} model.Topic
// @Failure 404 {object} errors.APIError "The topic doesn't exist"
// @Failure 500 {object} errors.APIError
// @Router /topics [get]
// @OperationId list-topic
func (t TopicController) List(c *gin.Context) {
	topics, apierr := service.TopicsService.ListTopics()
	if apierr != nil {
		c.JSON(apierr.Status, &apierr)
		return
	}

	c.JSON(http.StatusOK, &gin.H{"topics": topics})
}

// Delete godoc
// @Summary Delete a topic
// @Description Deletes the topic from database and resources
// @Tags topics
// @Accept json
// @Produce json
// @Param topic path string true "The name of the Topic"
// @Success 204
// @Failure 404 {object} errors.APIError "The topic doesn't exist"
// @Failure 500 {object} errors.APIError
// @Router /topics/{topic} [get]
// @OperationId get-topic
func (t TopicController) Delete(c *gin.Context) {
	adminToken := c.GetHeader("X-Admin-Token")
	topicName := c.Param("topic")

	apierr := service.TopicsService.DeleteTopic(topicName, adminToken)
	if apierr != nil {
		c.JSON(apierr.Status, &apierr)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetTopicSubscriptions godoc
// @Summary List topic subscribers
// @Description List all the topic subscribers
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} model.Subscriber
// @Failure 404 {object} errors.APIError "The topic doesn't exist"
// @Failure 500 {object} errors.APIError
// @Router /topics/{topic}/subscribers [get]
// @OperationId list-topic-subscribers
func (t TopicController) GetTopicSubscriptions(c *gin.Context) {

	topicName := c.Param("topic")
	subscriptions, apierr := service.SubscriptionsService.GetTopicSubscriptions(topicName)
	if apierr != nil {
		c.JSON(apierr.Status, &apierr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"topic": topicName, "subscribers": &subscriptions})
}
