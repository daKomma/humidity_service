package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Station struct {
	url         *url.URL
	added       time.Time
	updated     time.Time
	humidity    float32
	temperature float32
}
type StationResponse struct {
	Hum  float32 `json:"hum"`
	Temp float32 `json:"temp"`
}

func (s *Station) NewStation(rawUrl string) (*Station, error) {
	checkedUrl, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return s, err
	}

	s.url = checkedUrl
	s.added = time.Now()

	return s, nil
}

func (s *Station) UpdateData() {

	resp, err := http.Get(s.url.String())

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var result StationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	s.humidity = result.Hum
	s.temperature = result.Temp
	s.updated = time.Now()
}

func (s *Station) GetHumidity() float32 {
	return s.humidity
}

func (s *Station) GetTemperature() float32 {
	return s.temperature
}
