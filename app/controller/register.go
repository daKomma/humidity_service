package controller

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"humidity_service/main/models"

	"github.com/gin-gonic/gin"
)

// Controller for managing the Stations
type RegisterController struct{}

// Get one station by uuid or all stations
func (r RegisterController) Get(c *gin.Context) {
	// get uuid of station from url parameter
	stationId := c.Param("id")

	manager := models.GetManager()

	var stations []models.Station
	var err error

	if stationId != "/" {
		stationId, _ = strings.CutPrefix(stationId, "/")

		stations, err = manager.GetStation(stationId)
	} else {
		stations, err = manager.GetAllStation()
	}

	// if no lines than log and return 404
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stations)
	return
}

// Add station to the manager
func (r RegisterController) Add(c *gin.Context) {
	type Body struct {
		Url string `json:"url" binding:"required"`
	}

	var body Body

	// get body and if error handle it
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if url is valid
	checkedUrl, err := url.ParseRequestURI(body.Url)

	if err != nil {
		log.Println(err)
		resp := gin.H{
			"status": "url not valid",
			"url":    body.Url,
			"error":  err.Error()}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	manager := models.GetManager()

	station, err := manager.Add(checkedUrl.String())

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, station)
}

// Remove station from database
func (r RegisterController) Remove(c *gin.Context) {
	stationId := c.Param("id")

	manager := models.GetManager()

	if strings.Compare(stationId, "all") == 0 {
		success := manager.RemoveAllStation()

		c.JSON(http.StatusOK, gin.H{"success": success})
	}

	success := manager.Remove(stationId)

	c.JSON(http.StatusOK, gin.H{"success": success, "uuid": stationId})
}
