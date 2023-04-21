package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Station struct {
	url         string
	added       time.Time
	updated     time.Time
	humidity    float32
	temperature float32
}

func (s *Station) NewStation(url string) {
	s.url = url
	s.added = time.Now()
}

func (s *Station) UpdateData() {
	type Response struct {
		Hum  float32 `json:"hum"`
		Temp float32 `json:"temp"`
	}

	resp, err := http.Get(s.url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln(err)
	}

	s.humidity = result.Hum
	s.temperature = result.Temp
	s.updated = time.Now()
}
