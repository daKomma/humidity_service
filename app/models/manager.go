package models

import (
	"encoding/json"
	"errors"
	"humidity_service/main/db"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Manager struct {
}

// TODO Add Place and Paused flag
// Struct to parse the SQL response
type Station struct {
	Uuid    string    `json:"uuid"`
	Url     string    `json:"url"`
	Created time.Time `json:"created"`
	Place   string    `json:"place"`
}

type Data struct {
	Hum  float32   `json:"hum"`
	Temp float32   `json:"temp"`
	Time time.Time `json:"time"`
}

type StationData struct {
	Station Station `json:"station"`
	Data    []Data  `json:"data"`
}

type StationResponse struct {
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
	})

	return manager
}

// Add station to Database
func (m *Manager) Add(url string, place string) ([]Station, error) {
	isValid, err := m.testStation(url)

	if !isValid {
		return nil, err
	}

	// define values
	uuid := uuid.New().String()
	createdTime := time.Now().UTC()

	db := db.NewDb()
	defer db.Close()

	// insert station into DB
	insertStatement := `INSERT INTO Stations (uuid, url, created, place)
	VALUES (?, ?, ?, ?)`

	_, err = db.Exec(insertStatement, uuid, url, createdTime, place)

	if err != nil {
		return nil, err
	}

	return m.GetStation(uuid)
}

func (m *Manager) testStation(url string) (bool, error) {
	_, err := m.getStationData(url)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Get station with given uuid
func (m *Manager) GetStation(uuid string) ([]Station, error) {
	args := []interface{}{}
	query := "select * from Stations where uuid = ?"
	args = append(args, uuid)
	return m.getStationFromDB(query, args)
}

// Get all stations
func (m *Manager) GetAllStation() ([]Station, error) {
	args := []interface{}{}
	query := "select * from Stations"
	return m.getStationFromDB(query, args)
}

// helper function to do request to database
func (m *Manager) getStationFromDB(query string, args []interface{}) ([]Station, error) {

	db := db.NewDb()

	defer db.Close()

	rows, err := db.Query(query, args...)

	// Check for errors and handle those
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resStations := []Station{}

	station := Station{}

	// Fill array of Stations
	for rows.Next() {
		rows.Scan(&station.Uuid, &station.Url, &station.Created, &station.Place)
		resStations = append(resStations, station)
	}

	if len(resStations) == 0 {
		return nil, errors.New("No Station(s) found")
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

func (m *Manager) LiveData(stations []Station) []StationData {

	var wg sync.WaitGroup

	countStations := len(stations)

	dataChanel := make(chan StationData, countStations)
	wg.Add(countStations)

	for s := range stations {
		m.getLiveData(&stations[s], dataChanel, func() { wg.Done() })
	}

	go func() {
		defer close(dataChanel)
		wg.Wait()
	}()

	resp := []StationData{}

	for data := range dataChanel {
		resp = append(resp, data)
	}

	return resp
}

func (m *Manager) getLiveData(station *Station, stData chan<- StationData, onExit func()) {
	go func() {
		defer onExit()
		stationData, err := m.getStationData(station.Url)

		if err != nil {
			// TODO do something... pause or ignore????
			return
		}

		data := []Data{{stationData.Hum, stationData.Temp, time.Now()}}

		stData <- StationData{*station, data}
	}()
}

// Update all Stations
func (m *Manager) Update(stations []Station) {
	var wg sync.WaitGroup

	countStations := len(stations)

	failedChanel := make(chan string, countStations)
	successChanel := make(chan bool, countStations)
	wg.Add(countStations)

	for s := range stations {
		m.addStationData(&stations[s], failedChanel, successChanel, func() { wg.Done() })
	}

	go func() {
		defer close(failedChanel)
		defer close(successChanel)
		wg.Wait()
	}()

	for status := range successChanel {
		if !status {
			log.Println("Failed to save Data")
		}
	}

	// TODO Add to set station on paused
	// for data := range failedChanel {

	// }
}

func (m *Manager) addStationData(station *Station, failed chan<- string, success chan<- bool, onExit func()) {
	go func() {
		defer onExit()
		stationData, err := m.getStationData(station.Url)

		if err != nil {
			failed <- station.Uuid
		}

		success <- m.saveStationData(station, stationData)
	}()
}

// get Data from Station
func (m *Manager) getStationData(url string) (*StationResponse, error) {
	// get data from the Station
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.New("Something is wrong with the URL")
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)

	// parse body data
	var result StationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("Something is wrong with the Path or Response")
	}

	return &result, nil
}

// stores station data in database
func (m *Manager) saveStationData(station *Station, data *StationResponse) bool {
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
func (m *Manager) GetDBStationData(uuid string) ([]StationData, error) {
	args := []interface{}{}
	query := "select hum, temp, time from Data where station = ?"
	args = append(args, uuid)

	data, err := m.getDataFromDB(query, args)

	if err != nil {
		return nil, err
	}

	station, err := m.GetStation(uuid)

	if err != nil {
		return nil, err
	}

	stationData := []StationData{{station[0], data}}

	return stationData, nil
}

// Get all stations
func (m *Manager) GetAllData() ([]StationData, error) {
	allStation, err := m.GetAllStation()

	if err != nil {
		return nil, err
	}

	countStations := len(allStation)

	var wg sync.WaitGroup

	var stationData []StationData

	stationChanel := make(chan []StationData, countStations)
	errorChanel := make(chan error, countStations)
	wg.Add(countStations)

	for s := range allStation {
		m.getStationDataAsChanel(&allStation[s], stationChanel, errorChanel, func() { wg.Done() })
	}

	go func() {
		defer close(stationChanel)
		defer close(errorChanel)
		wg.Wait()
	}()

	for err := range errorChanel {
		return nil, err
	}

	for data := range stationChanel {
		stationData = append(stationData, data...)
	}

	return stationData, nil
}

func (m *Manager) getStationDataAsChanel(station *Station, stations chan<- []StationData, errors chan<- error, onExit func()) {
	go func() {
		defer onExit()

		data, err := m.GetDBStationData(station.Uuid)

		if err != nil {
			errors <- err
			return
		}

		stations <- data
	}()
}

// helper function to do request to database
func (m *Manager) getDataFromDB(query string, args []interface{}) ([]Data, error) {

	db := db.NewDb()

	defer db.Close()

	rows, err := db.Query(query, args...)

	// Check for errors and handle those
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resData := []Data{}

	data := Data{}

	// Fill array of Stations
	for rows.Next() {
		rows.Scan(&data.Hum, &data.Temp, &data.Time)

		resData = append(resData, data)
	}

	return resData, nil
}
