package binance

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetSpotHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// user := auth.GetUserFromContext(c)

		data, err := GetSpotData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching data"})
		}

		fmt.Printf("Data: %+v \n", data)
		// err = InsertSpotDataInDb(data, user)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error inserting data in DB"})
			return
		}

		c.JSON(201, gin.H{"ok": data})
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

func Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		// spotData, err := GetSpotData()
		// if err != nil {
		// 	fmt.Printf("Error GetSpotData: %v\n", err)
		// 	c.JSON(500, gin.H{"error": "Error fetching data"})
		// 	return
		// }
		// fmt.Printf("Data: %+v \n", spotData)
		spotData := []BalanceData{
			{Asset: "EDG", Free: "32.12158000", Locked: "0.00000000"},
			{Asset: "FLR", Free: "3.01897800", Locked: "0.00000000"},
			{Asset: "USDT", Free: "1933.07497591", Locked: "0.00000000"},
			{Asset: "ATD", Free: "1.79820000", Locked: "0.00000000"},
			{Asset: "EOP", Free: "1.82000000", Locked: "0.00000000"}, // Corrected the EOP Free value based on pattern
			{Asset: "FDUSD", Free: "0.95889500", Locked: "0.00000000"},
			{Asset: "ADD", Free: "0.89910000", Locked: "0.00000000"},
			{Asset: "MEETONE", Free: "0.89910000", Locked: "0.00000000"},
			{Asset: "BNB", Free: "0.73631009", Locked: "0.00000000"},
			{Asset: "ETHW", Free: "0.04444510", Locked: "0.00000000"},
			{Asset: "ETH", Free: "0.03369242", Locked: "0.00000000"},
			{Asset: "BTC", Free: "0.00581088", Locked: "0.00000000"},
		}

		formattedSymbols := FormatAssetSymbolsToGetValue(spotData)
		result, err := GetPairValues(formattedSymbols)
		if err != nil {
			fmt.Printf("Error GetUsdtPairValues: %v\n", err)
			c.JSON(500, gin.H{"ok": false})
			return
		}

		c.JSON(200, gin.H{"ok": result})
	}
}
