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

func NewRouter() *gin.Engine {
	once.Do(func() {
		docs.SwaggerInfo.BasePath = "/"
		router = gin.Default()

		router.UseRawPath = true

		router.UnescapePathValues = false

		health := new(controller.HealthController)

		router.GET("/health", health.Status)

		v1 := router.Group("/api/v1")
		{
			// data routes
			data := v1.Group("data")
			{
				dataController := new(controller.DataController)

				data.POST("/", dataController.GetSpecific)

				data.POST("/update", dataController.Update)

				data.GET("/live/*id", dataController.GetLive)
			}

			// station routes
			station := v1.Group("station")
			{
				register := new(controller.RegisterController)

				station.GET("/*id", register.Get)

				station.POST("/register", register.Add)

				station.DELETE("/:id", register.Remove)
			}
		}

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	})
	return router
}
