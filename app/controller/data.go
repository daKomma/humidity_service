package controller

import (
	"humidity_service/main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataController struct{}

func (d DataController) GetLive(c *gin.Context) {
	manager := models.GetManager()
	manager.UpdateAll()

	c.String(http.StatusOK, "Live Route")
}

func (d DataController) GetSpecific(c *gin.Context) {
	c.String(http.StatusOK, "Specific Route")
}
