package controllers

import (
	"fmt"
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CronSaveFundData() gin.HandlerFunc {
	return func(c *gin.Context) {
		fundData, err := services.GetFundData()
		// spotData, err := services.GetSpotData()

		if err != nil {
			fmt.Println("carajo se jodio el cron")
		}

		for _, data := range fundData {
			fmt.Printf("%+v\n", data)
		}

		c.JSON(http.StatusOK, gin.H{"message": "todo bien"})
	}
}
