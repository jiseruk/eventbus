package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		var message map[string]interface{}
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("Message received: %v \n", message)

		payload := message["payload"].(map[string]interface{})
		if payload != nil && payload["fail"] == true {
			c.JSON(http.StatusInternalServerError, &message)
			return
		}
		c.JSON(http.StatusOK, &message)

	})
	router.Run(":9000")
}
