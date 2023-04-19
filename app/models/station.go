package models

import "time"

type Station struct {
	url         string
	updated     time.Time
	humidity    float32
	temperature float32
}

func (s *Station) NewStation(url string) {
	s.url = url
}
