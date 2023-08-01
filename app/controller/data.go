package controller

import (
	"humidity_service/main/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Controller for managing data
type DataController struct{}

type DataBody struct {
	Uuid string `json:"uuid"`
}

// GetLive godoc
// @Summary Get live data
// @Description Get live data from one or all stations
// @Tags Data
// @Produce json
// @Param uuid path string false "Station ID"
// @Success 200 {array} models.StationData
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /data/live [get]
func (d *DataController) GetLive(c *gin.Context) {
	manager := models.GetManager()

	// get station id from url parameter
	id := c.Param("id")

	var stations []models.Station
	var err error

	if id != "/" {
		id, _ = strings.CutPrefix(id, "/")
		stations, err = manager.GetStation(id)
	} else {
		stations, err = manager.GetAllStation()
	}

	if err != nil {
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, manager.LiveData(stations))
}

// GetSpecific godoc
// @Summary Get data
// @Description Get all data of one or all stations
// @Tags Data
// @Accept  json
// @Produce json
// @Param request body controller.DataBody true "query params"
// @Success 200 {array} models.StationData
// @Failure 400 {object}  controller.JSONBadReqResult
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /data [post]
func (d DataController) GetSpecific(c *gin.Context) {

	var body DataBody

	// get body and if error handle it
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, JSONBadReqResult{
			Code:    http.StatusBadRequest,
			Message: "Wrong body",
		})
		return
	}

	manager := models.GetManager()

	var data []models.StationData
	var err error

	if body.Uuid != "" {
		data, err = manager.GetDBStationData(body.Uuid)
	} else {
		data, err = manager.GetAllData()
	}

	if err != nil {
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// Update godoc
// @Summary Update data
// @Description Updates data from all stations and stores it in the DB
// @Tags Data
// @Produce json
// @Success 200 {object} controller.JSONSuccessResult
// @Failure 404 {object}  controller.JSONNotFoundResult
// @Router /update [post]
func (d *DataController) Update(c *gin.Context) {
	manager := models.GetManager()

	stations, err := manager.GetAllStation()

	if err != nil {
		c.JSON(http.StatusNotFound, JSONNotFoundResult{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	manager.Update(stations)

	c.JSON(http.StatusOK, JSONSuccessResult{
		Code:    http.StatusOK,
		Message: "All stations updated.",
	})
	return
}
