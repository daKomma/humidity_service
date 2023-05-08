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

	stations := make([]models.Station, 0, len(manager.Stations))

	for _, station := range manager.Stations {
		// stationJSON, _ := json.Marshal(station)
		// stations = append(stations, stationJSON...)
		stations = append(stations, *station)
	}

	c.IndentedJSON(http.StatusOK, stations)
}

func (d DataController) GetSpecific(c *gin.Context) {
	c.String(http.StatusOK, "Specific Route")
}
