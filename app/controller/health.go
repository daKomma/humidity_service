package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller for managing Health information
type HealthController struct{}

type JSONHealthResult struct {
	Code      int       `json:"code" example:"200"`
	Message   string    `json:"message" example:"healthy i guess"`
	Status    string    `json:"status" example:"up"`
	Timestamp time.Time `json:"timestamp" example:"2023-07-29T07:52:50Z"`
}

// Status godoc
// @Summary Health check
// @Description returns health status of the server
// @Tags Health
// @Produce json
// @Success 200 {object} controller.JSONHealthResult
// @Router /health [get]
func (h HealthController) Status(c *gin.Context) {
	timestamp := time.Now()

	c.JSON(http.StatusOK, JSONHealthResult{
		Code:      http.StatusOK,
		Status:    "UP",
		Message:   "healthy i guess",
		Timestamp: timestamp,
	})
}
