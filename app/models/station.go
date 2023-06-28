package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"humidity_service/main/db"

	"github.com/google/uuid"
)

type Station struct {
	Id          uuid.UUID
	Url         *url.URL
	Added       time.Time
	Updated     time.Time
	Humidity    float32
	Temperature float32
}
type StationResponse struct {
	Hum  float32 `json:"hum"`
	Temp float32 `json:"temp"`
}

// Creates new Station and stores it in database
func (s *Station) NewStation(rawUrl string) (*Station, error) {
	// Check if url is valid
	checkedUrl, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return s, err
	}

	// set values
	s.Id = uuid.New()
	s.Url = checkedUrl
	s.Added = time.Now().UTC()

	// store new Station in Database
	db := db.NewDb()

	defer db.Close()

	insertStatement := `INSERT INTO Stations (uuid, url, created)
	VALUES (?, ?, ?)`
	_, err = db.Exec(insertStatement, s.Id, s.Url.String(), s.Added)

	if err != nil {
		log.Println(err)
	}

	return s, nil
}

// Update data of this Station
func (s *Station) UpdateData() {
	// get data from the Station
	resp, err := http.Get(s.Url.String())

	if err != nil {
		log.Fatalln(err)
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)

	var result StationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	s.Humidity = result.Hum
	s.Temperature = result.Temp
	s.Updated = time.Now().UTC()

	// store new Values in Database
	db := db.NewDb()

	defer db.Close()

	insertStatement := `INSERT INTO Data (hum, temp, time, station)
	VALUES (?, ?, ?, ?)`
	db.Exec(insertStatement, s.Humidity, s.Temperature, s.Updated, s.Id)
}

// Custom MarshalJson for Stations
func (s *Station) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id          string    `json:"id"`
		Url         string    `json:"url"`
		Added       time.Time `json:"added"`
		Updated     time.Time `json:"updated"`
		Humidity    float32   `json:"hum"`
		Temperature float32   `json:"temp"`
	}{
		Id:          s.Id.String(),
		Url:         s.Url.String(),
		Added:       s.Added,
		Updated:     s.Updated,
		Humidity:    s.Humidity,
		Temperature: s.Temperature,
	})
}
