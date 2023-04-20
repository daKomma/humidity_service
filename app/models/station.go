package models

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Station struct {
	url         string
	updated     time.Time
	humidity    float32
	temperature float32
}

func (s *Station) NewStation(url string) {
	s.url = url
}

func (s *Station) UpdateData() {
	resp, err := http.Get(s.url)

	if err != nil {
		log.Fatalln(err)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

}
