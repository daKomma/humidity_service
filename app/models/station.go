package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Station struct {
	Id          uuid.UUID
	Url         *url.URL
	Added       time.Time
	Updated     time.Time
	Humidity    float32
	Temperature float32
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

	s.Id = uuid.New()
	s.Url = checkedUrl
	s.Added = time.Now()

	return s, nil
}

func (s *Station) UpdateData() {

	resp, err := http.Get(s.Url.String())

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var result StationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	s.Humidity = result.Hum
	s.Temperature = result.Temp
	s.Updated = time.Now()
}
