package server

import (
	"humidity_service/main/controller"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	once   sync.Once
)

func NewRouter() *gin.Engine {
	once.Do(func() {
		router = gin.Default()

		health := new(controller.HealthController)

		register := new(controller.RegisterController)

		router.GET("/health", health.Status)

		router.POST("/register", register.Add)

		data := router.Group("data")
		{
			dataController := new(controller.DataController)

			data.GET("/", dataController.GetSpecific)

			data.GET("/live", dataController.GetLive)
		}

	})
	return router
}
