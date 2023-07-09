package controller

import (
	"humidity_service/main/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller for managing data
type DataController struct{}

// Get live data from one or all stations
func (d *DataController) GetLive(c *gin.Context) {
	manager := models.GetManager()

	// get station id from url parameter
	id := c.Param("id")

	if id != "/" {
		id, _ = strings.CutPrefix(id, "/")

	} else {
		// allStation, _ := manager.GetAllStation()
		// manager.Update(allStation)
		stations, _ := manager.GetAllStation()
		c.JSON(http.StatusOK, manager.LiveData(stations))
		return
	}

	// // Create array of stations and fill it
	// stations := make([]models.Station, 0, len(manager.Stations))

	// for _, station := range manager.Stations {
	// 	stations = append(stations, *station)
	// }

	// Send all found Stations
	// c.IndentedJSON(http.StatusOK, stations)

}

// TODO
func (d DataController) GetSpecific(c *gin.Context) {
	c.String(http.StatusOK, "Specific Route")
}

func (d *DataController) GetAll(c *gin.Context) {
	manager := models.GetManager()
	data, err := manager.GetAllData()

	// if no lines than log and return 404
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
	return
}
