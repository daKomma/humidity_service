package server

func Init() {
	router := NewRouter()

	NewCron()

	router.Run("localhost:8080")
}
