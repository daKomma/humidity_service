package controller

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"humidity_service/main/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Controller for managing the Stations
type RegisterController struct{}

// Get one station by uuid or all stations
func (r RegisterController) Get(c *gin.Context) {
	// get uuid of station from url parameter
	stationId := c.Param("id")

	// Struct to parse the SQL response
	type Station struct {
		Uuid    string    `json:"uuid"`
		Url     string    `json:"url"`
		Created time.Time `json:"created"`
	}

	// get DB and close the connection at the end
	db := db.NewDb()

	defer db.Close()

	if stationId != "/" {
		stationId, _ = strings.CutPrefix(stationId, "/")

		query := "select * from Stations where uuid = ?"

		rows := db.QueryRow(query, stationId)

		station := Station{}

		// Scan row and parse into variable
		err := rows.Scan(&station.Uuid, &station.Url, &station.Created)

		// if no lines than log and return 404
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

		// Check for errors and handle those
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

		// Fill array of Stations
		for rows.Next() {
			rows.Scan(&station.Uuid, &station.Url, &station.Created)
			resStations = append(resStations, station)
		}

		c.JSON(http.StatusOK, resStations)
		return
	}
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

	// define values
	uuid := uuid.New()
	createdTime := time.Now().UTC()

	// store new Station in Database
	db := db.NewDb()

	defer db.Close()

	// insert station into DB
	insertStatement := `INSERT INTO Stations (uuid, url, created)
	VALUES (?, ?, ?)`
	_, err = db.Exec(insertStatement, uuid, checkedUrl.String(), createdTime)

	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"url": body.Url, "uuid": uuid})
}

// Remove station from database
func (r RegisterController) Remove(c *gin.Context) {
	stationId := c.Param("id")

	db := db.NewDb()
	defer db.Close()

	query := "delete from Stations where uuid = ?"

	_, err := db.Exec(query, stationId)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "uuid": stationId})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "uuid": stationId})
}
