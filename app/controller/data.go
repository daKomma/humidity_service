package controller

import (
	"humidity_service/main/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DataController struct{}

func (d DataController) GetLive(c *gin.Context) {
	manager := models.GetManager()
	id := c.Param("id")

	if id != "/" {
		id, _ = strings.CutPrefix(id, "/")

		manager.Update(id)
	} else {
		manager.UpdateAll()
	}

	stations := make([]models.Station, 0, len(manager.Stations))

	for _, station := range manager.Stations {
		stations = append(stations, *station)
	}

	c.IndentedJSON(http.StatusOK, stations)
}

func (d DataController) GetSpecific(c *gin.Context) {
	c.String(http.StatusOK, "Specific Route")
}
