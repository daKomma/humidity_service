package models

import (
	"sync"

	"github.com/google/uuid"
)

type Manager struct {
	Stations map[uuid.UUID]*Station
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
	m.Stations[station.id] = station
}

func (m *Manager) UpdateAll() {
	for _, station := range m.Stations {
		go station.UpdateData()
	}
}
