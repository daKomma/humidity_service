package models

import (
	"bufio"
	"humidity_service/main/db"
	"log"
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

	return m.Get(uuid)
}

// Get station with given uuid
func (m *Manager) Get(uuid string) ([]DBStation, error) {
	args := []interface{}{}
	query := "select * from Stations where uuid = ?"
	args = append(args, uuid)
	return m.getStationFromDB(query, args)
}

// Get all stations
func (m *Manager) GetAll() ([]DBStation, error) {
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

// Remove station from manager
func (m *Manager) Remove(id string) {
	delete(m.Stations, id)
}

// Update all Stations
func (m *Manager) UpdateAll() {
	var wg sync.WaitGroup
	wg.Add(len(m.Stations))

	for _, station := range m.Stations {
		go func() {
			station.UpdateData()
			wg.Done()
		}()
	}

	wg.Wait()
}

// Update station with given uuid
func (m *Manager) Update(id string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		m.Stations[id].UpdateData()
		wg.Done()
	}()

	wg.Wait()
}
