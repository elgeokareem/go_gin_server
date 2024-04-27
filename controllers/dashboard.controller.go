package controllers

import (
	"goGinServer/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardController() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallet, err := services.GetFundData() // funds
		// wallet, err := services.GetSpotData() // spot

		if err != nil {
			log.Printf("Error fetching wallet data: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": wallet})
	}
}
