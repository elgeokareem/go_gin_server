package binance

import (
	"github.com/gin-gonic/gin"
)

func BinanceRoutes(router *gin.Engine) {
	binanceRoutes := router.Group("/binance")

	binanceRoutes.GET("/spot", GetSpotHandler())
	binanceRoutes.GET("/fund", GetFundHandler())
}
