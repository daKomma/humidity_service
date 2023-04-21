package controller

import (
	"net/http"

	"humidity_service/main/models"

	"github.com/gin-gonic/gin"
)

type RegisterController struct{}

type Body struct {
	Url string `json:"url" binding:"required"`
}

func (r RegisterController) Add(c *gin.Context) {
	var body Body

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station := new(models.Station)
	station.NewStation(body.Url)

	manager := models.GetManager()
	manager.Add(station)

	c.String(http.StatusOK, body.Url)
}
