package server

import (
	"humidity_service/main/controller"
	docs "humidity_service/main/docs"
	"sync"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	once   sync.Once
)

// @Title Humidity service API
// @version 1.0
// @description This is a simple server to get and store data from multiple sensor stations.
// @host localhost:8080
func NewRouter() *gin.Engine {
	once.Do(func() {
		docs.SwaggerInfo.BasePath = "/"
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

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	})
	return router
}
