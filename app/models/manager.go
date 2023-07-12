package models

import (
	"bufio"
	"encoding/json"
	"humidity_service/main/db"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Manager struct {
	Stations map[string]*Station
}

// Struct to parse the SQL response
type DBStation struct {
	Uuid    string    `json:"uuid"`
	Url     string    `json:"url"`
	Created time.Time `json:"created"`
}

type DBStationData struct {
	Hum         string    `json:"hum"`
	Temp        string    `json:"temp"`
	Time        time.Time `json:"time"`
	StationUUID time.Time `json:"station"`
}

type DBStationResponse struct {
	Hum  float32 `json:"hum"`
	Temp float32 `json:"temp"`
}

var (
	manager *Manager
	once    sync.Once
)

// Get instance of Manager
func GetManager() *Manager {
	once.Do(func() {
		manager = new(Manager)

		manager.Stations = make(map[string]*Station)

		manager.loadFromFile(os.Getenv("URLPATH"))
	})

	return manager
}

// Load station from given file
func (m *Manager) loadFromFile(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	// Foreach line create new Station object
	for fileScanner.Scan() {
		url, err := url.ParseRequestURI(fileScanner.Text())

		if err != nil {
			continue
		}

		station := new(Station)

		station.NewStation(url.String())

		m.Add("station")
	}
}

// Add station to Database
func (m *Manager) Add(url string) ([]DBStation, error) {
	// define values
	uuid := uuid.New().String()
	createdTime := time.Now().UTC()

	db := db.NewDb()
	defer db.Close()

	// insert station into DB
	insertStatement := `INSERT INTO Stations (uuid, url, created)
	VALUES (?, ?, ?)`

	_, err := db.Exec(insertStatement, uuid, url, createdTime)

	if err != nil {
		return nil, err
	}

	return m.GetStation(uuid)
}

// Get station with given uuid
func (m *Manager) GetStation(uuid string) ([]DBStation, error) {
	args := []interface{}{}
	query := "select * from Stations where uuid = ?"
	args = append(args, uuid)
	return m.getStationFromDB(query, args)
}

// Get all stations
func (m *Manager) GetAllStation() ([]DBStation, error) {
	args := []interface{}{}
	query := "select * from Stations"
	return m.getStationFromDB(query, args)
}

// helper function to do request to database
func (m *Manager) getStationFromDB(query string, args []interface{}) ([]DBStation, error) {

	db := db.NewDb()

	defer db.Close()

	rows, err := db.Query(query, args...)

	// Check for errors and handle those
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resStations := []DBStation{}

	station := DBStation{}

	// Fill array of Stations
	for rows.Next() {
		rows.Scan(&station.Uuid, &station.Url, &station.Created)
		resStations = append(resStations, station)
	}

	return resStations, nil
}

// Remove station from Database
func (m *Manager) Remove(uuid string) bool {
	db := db.NewDb()
	defer db.Close()

	query := "delete from Stations where uuid = ?"

	_, err := db.Exec(query, uuid)

	if err != nil {
		return false
	}

	return true
}

// Remove all station from Database => only in dev
func (m *Manager) RemoveAllStation() bool {
	db := db.NewDb()
	defer db.Close()

	query := "delete from Stations"
	_, err := db.Exec(query)

	if err != nil {
		return false
	}

	return true
}

func (m *Manager) LiveData(stations []DBStation) any {

	var wg sync.WaitGroup
	wg.Add(len(stations))

	type stationLiveData struct {
		DBStation         `json:"station"`
		DBStationResponse `json:"data"`
	}

	resp := []stationLiveData{}

	for s := range stations {
		go func(station *DBStation) {
			stationData := m.getStationData(station.Url)
			liveData := &stationLiveData{*station, stationData}

			resp = append(resp, *liveData)
			wg.Done()
		}(&stations[s])
	}

	wg.Wait()

	return resp
}

// Update all Stations
func (m *Manager) Update(stations []DBStation) {
	var wg sync.WaitGroup
	wg.Add(len(stations))

	for s := range stations {
		go func(station *DBStation) {
			stationData := m.getStationData(station.Url)
			isSaved := m.saveStationData(station, &stationData)

			log.Println("Station: %s status: %t", station.Uuid, isSaved)
			wg.Done()
		}(&stations[s])
	}

	wg.Wait()
}

// get Data from Station
func (m *Manager) getStationData(url string) DBStationResponse {
	// get data from the Station
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)

	// parse body data
	var result DBStationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	return result
}

// stores station data in database
func (m *Manager) saveStationData(station *DBStation, data *DBStationResponse) bool {
	// store new Values in Database
	db := db.NewDb()

	defer db.Close()

	insertStatement := `INSERT INTO Data (hum, temp, time, station)
	VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertStatement, data.Hum, data.Temp, time.Now().UTC(), station.Uuid)

	if err != nil {
		return false
	}

	return true
}

// Get station with given uuid
func (m *Manager) GetStationData(uuid string) ([]DBStationData, error) {
	args := []interface{}{}
	query := "select * from Data where station = ?"
	args = append(args, uuid)
	return m.getDataFromDB(query, args)
}

// Get all stations
func (m *Manager) GetAllData() ([]DBStationData, error) {
	args := []interface{}{}
	query := "select * from Data"
	return m.getDataFromDB(query, args)
}

// helper function to do request to database
func (m *Manager) getDataFromDB(query string, args []interface{}) ([]DBStationData, error) {

	db := db.NewDb()

	defer db.Close()

	rows, err := db.Query(query, args...)

	// Check for errors and handle those
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resData := []DBStationData{}

	data := DBStationData{}

	// Fill array of Stations
	for rows.Next() {
		rows.Scan(&data.Hum, &data.Temp, &data.Time, &data.StationUUID)
		resData = append(resData, data)
	}

	return resData, nil
}
