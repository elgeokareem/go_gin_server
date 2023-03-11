package controllers

import (
	"fmt"
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardController() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Entra al controler")
		// services.Service()
		services.GetWalletData()

		c.JSON(http.StatusOK, gin.H{"message": "This is the end"})
	}
}
