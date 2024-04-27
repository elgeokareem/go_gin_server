package controllers

import (
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SqlGroupByController() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := services.GroupByService()

		c.JSON(http.StatusOK, data)
	}
}
