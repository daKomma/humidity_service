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

type JSONSuccessResult struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Wrong Parameter"`
}
type JSONBadReqResult struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Wrong Parameter"`
}

type JSONNotFoundResult struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Not found"`
}

type Body struct {
	Url   string `json:"url" binding:"required" example:"http://localhost:8080/data"`
	Place string `json:"place" binding:"required" example:"Garden"`
}

// Get godoc
// @Summary get station
// @Description Get all stations or one by its UUID
// @Produce json
// @Param uuid path string false "Station ID"
// @Success 200 {array} models.Station
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /station [get]
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
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stations)
	return
}

// Add godoc
// @Summary Create new Station
// @Description Add new station to the service and database
// @Accept  json
// @Produce json
// @Param request body controller.Body true "query params"
// @Success 200 {array} models.Station
// @Failure 400 {object}  controller.JSONBadReqResult
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /station/register [post]
func (r RegisterController) Add(c *gin.Context) {
	var body Body

	// get body and if error handle it
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, JSONBadReqResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// Check if url is valid
	checkedUrl, err := url.ParseRequestURI(body.Url)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, JSONBadReqResult{
			Code:    http.StatusBadRequest,
			Message: "url not valid: " + body.Url,
		})
		return
	}

	manager := models.GetManager()

	station, err := manager.Add(checkedUrl.String(), body.Place)

	if err != nil {
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, station)
}

// Remove godoc
// @Summary Remove a Station
// @Description Removes station with given UUID from the service and database
// @Produce json
// @Param uuid path string true "Station ID"
// @Success 200 {object} controller.JSONSuccessResult
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /station [delete]
func (r RegisterController) Remove(c *gin.Context) {
	stationId := c.Param("id")

	manager := models.GetManager()

	if strings.Compare(stationId, "all") == 0 {
		success := manager.RemoveAllStation()

		c.JSON(http.StatusOK, gin.H{"success": success})
	}

	_, err := manager.Remove(stationId)

	if err != nil {
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: stationId,
		})
		return
	}

	c.JSON(http.StatusOK, JSONSuccessResult{
		Code:    http.StatusOK,
		Message: stationId,
	})
}
