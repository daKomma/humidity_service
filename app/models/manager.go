package models

import "sync"

type Manager struct {
	stations []*Station
}

var (
	manager *Manager
	once    sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		manager.loadFromFile()
	})

	return manager
}

func (m *Manager) loadFromFile() {
	// TODO add load from File
}

func (m *Manager) Add(station *Station) {
	m.stations = append(m.stations, station)
}
