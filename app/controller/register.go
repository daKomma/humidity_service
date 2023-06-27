package controller

import (
	"net/http"
	"strings"

	"humidity_service/main/models"

	"github.com/gin-gonic/gin"
)

type RegisterController struct{}

func (r RegisterController) Get(c *gin.Context) {
	stationId := c.Param("id")

	if stationId != "/" {
		stationId, _ = strings.CutPrefix(stationId, "/")

	} else {

	}

	c.JSON(http.StatusOK, gin.H{"Id": stationId})
}

func (r RegisterController) Add(c *gin.Context) {
	type Body struct {
		Url string `json:"url" binding:"required"`
	}

	var body Body

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station := new(models.Station)
	station.NewStation(body.Url)

	manager := models.GetManager()
	manager.Add(station)

	c.JSON(http.StatusOK, gin.H{"url": body.Url, "uuid": station.Id})
}

func (r RegisterController) Remove(c *gin.Context) {
	stationId := c.Param("id")

	manager := models.GetManager()
	manager.Remove(stationId)

	c.JSON(http.StatusOK, nil)
}
