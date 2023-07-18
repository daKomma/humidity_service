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

	var stations []models.Station

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
	type Body struct {
		Uuid string `json:"uuid" required:true`
	}

	var body Body

	// get body and if error handle it
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manager := models.GetManager()

	var data []models.StationData
	var err error

	if body.Uuid != "" {
		data, err = manager.GetStationData(body.Uuid)
	} else {
		data, err = manager.GetAllData()
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (d *DataController) Update(c *gin.Context) {
	manager := models.GetManager()

	stations, _ := manager.GetAllStation()

	manager.Update(stations)

	c.JSON(http.StatusOK, gin.H{})
	return
}
