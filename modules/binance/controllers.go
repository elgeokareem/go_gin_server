package binance

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetSpotHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := GetSpotData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching data"})
		}

		fmt.Printf("Data: %+v \n", data)
		c.JSON(200, gin.H{"data": data})
	}
}

func GetFundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := GetFundData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching data"})
		}

		fmt.Printf("Data: %+v \n", data)
		c.JSON(200, gin.H{"data": data})
	}
}
