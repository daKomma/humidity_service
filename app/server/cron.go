package server

import (
	"fmt"
	"sync"

	"github.com/robfig/cron"
)

var (
	cronRunner *cron.Cron
	onceCron   sync.Once
)

func NewCron() *cron.Cron {
	onceCron.Do(func() {
		cronRunner = cron.New()

		cronRunner.AddFunc("* * * * * *", func() { fmt.Printf("I am running every second!") })

		cronRunner.Start()
	})

	return cronRunner
}
