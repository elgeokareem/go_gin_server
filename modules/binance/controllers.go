package binance

import (
	"fmt"
	"goGinServer/modules/auth"

	"github.com/gin-gonic/gin"
)

func GetSpotHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := auth.GetUserFromContext(c)

		data, err := GetSpotData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching data"})
		}

		fmt.Printf("Data: %+v \n", data)
		// Now i need to store the data in the db
		err = InsertSpotDataInDb(data)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error inserting data in DB"})
			return
		}

		c.JSON(201, gin.H{"ok": true})
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
