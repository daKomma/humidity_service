package models

import "sync"

type Manager struct {
	Stations []*Station
}

var (
	manager *Manager
	once    sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		manager = new(Manager)

		manager.loadFromFile()
	})

	return manager
}

func (m *Manager) loadFromFile() {
	// TODO add load from File
}

func (m *Manager) Add(station *Station) {
	m.Stations = append(m.Stations, station)
}
