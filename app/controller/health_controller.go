package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
}

func (t HealthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
