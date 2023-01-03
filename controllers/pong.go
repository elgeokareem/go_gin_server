package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PongController() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is the a ver"})
	}
}
