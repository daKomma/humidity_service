package test_test

import (
	"humidity_service/main/models"
	"testing"
)

func TestGetManager(t *testing.T) {
	manager := models.GetManager()
	newManager := models.GetManager()

	if manager != newManager {
		t.Errorf("Manager are not same!")
	}
}

// func TestAdd(t *testing.T) {
// 	station := new(models.Station)

// 	station.NewStation("http://localhost:3000")

// 	manager := models.GetManager()

// 	manager.Add(station)

// 	if len(manager.Stations) != 1 {
// 		t.Errorf("ERROR: expected 1 but got %d", len(manager.Stations))
// 	}
// }
