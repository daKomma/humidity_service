package server

import (
	"humidity_service/main/models"
	"os"
	"sync"

	"github.com/robfig/cron"
)

var (
	cronRunner *cron.Cron
	onceCron   sync.Once
)

// Init cron runner
// Singleton
func NewCron() *cron.Cron {
	onceCron.Do(func() {
		cronRunner = cron.New()

		manager := models.GetManager()

		cronRunner.AddFunc(os.Getenv("CRON_INTERVAL"), func() {
			stations, _ := manager.GetAllStation()
			manager.Update(stations)
		})

		cronRunner.Start()
	})

	return cronRunner
}
