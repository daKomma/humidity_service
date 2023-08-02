package server

import "humidity_service/main/db"

func Init() {
	db.NewDb()

	router := NewRouter()

	NewCron()

	router.Run()
}
