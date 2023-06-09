package server

func Init() {
	NewDb()

	router := NewRouter()

	// NewCron()

	router.Run("localhost:8080")
}
