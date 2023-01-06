package controllers

import (
	"fmt"
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SqlGroupByController() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Entra al controlador de SQL")

		// we can use a service here.
		data := services.GroupByService()

		c.JSON(http.StatusOK, data)
	}
}
