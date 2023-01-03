package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Users(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get user by id"})
	}
}
