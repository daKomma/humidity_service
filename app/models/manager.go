package models

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"sync"
)

type Manager struct {
	Stations map[string]*Station
}

var (
	manager *Manager
	once    sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		manager = new(Manager)

		manager.Stations = make(map[string]*Station)

		manager.loadFromFile(os.Getenv("URLPATH"))
	})

	return manager
}

func (m *Manager) loadFromFile(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		url, err := url.ParseRequestURI(fileScanner.Text())

		if err != nil {
			continue
		}

		station := new(Station)

		station.NewStation(url.String())

		m.Add(station)
	}
}

func (m *Manager) Add(station *Station) {
	m.Stations[station.Id.String()] = station
}

func (m *Manager) Remove(id string) {
	delete(m.Stations, id)
}

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

func (m *Manager) Update(id string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		m.Stations[id].UpdateData()
		wg.Done()
	}()

	wg.Wait()
}
