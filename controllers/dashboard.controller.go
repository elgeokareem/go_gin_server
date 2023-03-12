package controllers

import (
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the binance wallet
		wallet := services.GetWalletData()

		// c.JSON(http.StatusOK, gin.H{"data": walletData})
		c.JSON(http.StatusOK, gin.H{"data": wallet})
	}
}
