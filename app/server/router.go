package server

import (
	"humidity_service/main/controller"
	"sync"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var once sync.Once

func NewRouter() *gin.Engine {
	once.Do(func() {
		router = gin.Default()

		health := new(controller.HealthController)

		router.GET("/health", health.Status)

	})
	return router
}
