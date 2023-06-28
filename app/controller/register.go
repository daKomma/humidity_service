package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"humidity_service/main/db"
	"humidity_service/main/models"

	"github.com/gin-gonic/gin"
)

type RegisterController struct{}

func (r RegisterController) Get(c *gin.Context) {
	stationId := c.Param("id")

	type Station struct {
		Uuid    string    `json:"uuid"`
		Url     string    `json:"url"`
		Created time.Time `json:"created"`
	}

	db := db.NewDb()

	defer db.Close()

	if stationId != "/" {
		stationId, _ = strings.CutPrefix(stationId, "/")

		query := "select * from Stations where uuid = ?"

		rows := db.QueryRow(query, stationId)

		station := Station{}

		err := rows.Scan(&station.Uuid, &station.Url, &station.Created)

		if err != nil && err == sql.ErrNoRows {
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{stationId: "Not Found"})
			return
		}

		c.JSON(http.StatusOK, station)
		return
	} else {
		query := "select * from Stations"

		rows, err := db.Query(query)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Println(err)
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}

			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		resStations := []Station{}

		station := Station{}

		for rows.Next() {
			rows.Scan(&station.Uuid, &station.Url, &station.Created)
			resStations = append(resStations, station)
		}

		c.JSON(http.StatusOK, resStations)
		return
	}
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
