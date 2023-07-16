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

// Init all routs of the router
func NewRouter() *gin.Engine {
	once.Do(func() {
		router = gin.Default()

		router.UseRawPath = true

		router.UnescapePathValues = false

		health := new(controller.HealthController)

		router.GET("/health", health.Status)

		// data routes
		data := router.Group("data")
		{
			dataController := new(controller.DataController)

			data.POST("/", dataController.GetSpecific)

			data.POST("/update", dataController.Update)

			data.GET("/live/*id", dataController.GetLive)
		}

		// station routes
		station := router.Group("station")
		{
			register := new(controller.RegisterController)

			station.GET("/*id", register.Get)

			station.POST("/register", register.Add)

			station.DELETE("/:id", register.Remove)
		}

	})
	return router
}
