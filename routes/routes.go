package routes

import (
	"goGinServer/modules/binance"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	router := gin.Default()
	// add cors
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true, // Important for setting cookies cross-origin
	}
	router.Use(cors.New(config))

	getRoutes(router)
	router.Run()
}

// Get all the routes
func getRoutes(router *gin.Engine) {
	Auth(router)
	binance.BinanceRoutes(router)
}
