package controller

import (
	"humidity_service/main/models"
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

	var stations []models.DBStation

	if id != "/" {
		id, _ = strings.CutPrefix(id, "/")
		stations, _ = manager.GetStation(id)
	} else {
		stations, _ = manager.GetAllStation()
	}

	c.JSON(http.StatusOK, manager.LiveData(stations))
	return
}

// TODO
func (d DataController) GetSpecific(c *gin.Context) {
	c.String(http.StatusOK, "Specific Route")
}

func (d *DataController) Update(c *gin.Context) {
	manager := models.GetManager()

	stations, _ := manager.GetAllStation()

	manager.Update(stations)

	c.JSON(http.StatusOK, gin.H{})
	return
}
