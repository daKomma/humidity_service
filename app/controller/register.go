package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterController struct{}

func (r RegisterController) Add(c *gin.Context) {
	c.String(http.StatusOK, "Register")
}
