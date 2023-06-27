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

		router.UseRawPath = true

		router.UnescapePathValues = false

		health := new(controller.HealthController)

		register := new(controller.RegisterController)

		router.GET("/health", health.Status)

		data := router.Group("data")
		{
			dataController := new(controller.DataController)

			data.GET("/", dataController.GetSpecific)

			data.GET("/live/*id", dataController.GetLive)
		}

		station := router.Group("station")
		{
			station.GET("/*id", register.Get)

			station.POST("/register", register.Add)

			station.DELETE("/:id", register.Remove)
		}

	})
	return router
}
