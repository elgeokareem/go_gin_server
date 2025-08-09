package binance

import (
	"github.com/gin-gonic/gin"
)

func BinanceRoutes(router *gin.Engine) {
	binanceRoutes := router.Group("/binance")

	// Middleware
	// binanceRoutes.Use(middlewares.AuthMiddleware())

	binanceRoutes.GET("/spot", GetSpotHandler())
	binanceRoutes.GET("/fund", GetFundHandler())
	binanceRoutes.GET("/test", Test())
}
