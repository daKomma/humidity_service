package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	timestamp := time.Now().Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"message":   "healthy i guess",
		"timestamp": timestamp,
	})
}
