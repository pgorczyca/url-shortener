package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive": true,
	})
}
