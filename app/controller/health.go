package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller for managing Health information
type HealthController struct{}

// Status godoc
// @Summary health check
// @Description returns health status of the server
// @Produce json
// @Success 200 {string} asd
// @Router /health [get]
func (h HealthController) Status(c *gin.Context) {
	timestamp := time.Now().Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"message":   "healthy i guess",
		"timestamp": timestamp,
	})
}
