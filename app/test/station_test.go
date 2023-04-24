package test_test

import (
	"humidity_service/main/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Tests struct {
	name          string
	server        *httptest.Server
	response      *models.StationResponse
	expectedError error
}

func TestNewStation(t *testing.T) {
	url := "http://localhost:3000"
	station := new(models.Station)
	station, err := station.NewStation(url)

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestUpdateData(t *testing.T) {
	tests := []Tests{
		{
			name: "update-station-information",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "hum": 50.50, "temp": 26.5 }`))
			})),
			response: &models.StationResponse{
				Hum:  50.50,
				Temp: 26.5,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			station := new(models.Station)

			station.NewStation(test.server.URL)

			station.UpdateData()

			if station.GetHumidity() != test.response.Hum {
				t.Errorf("ERROR: wrong humidity. Expected: %f but got: %f", test.response.Hum, station.GetHumidity())
			}

			if station.GetTemperature() != test.response.Temp {
				t.Errorf("ERROR: wrong temperature. Expected: %f but got: %f", test.response.Temp, station.GetTemperature())
			}
		})
	}
}
