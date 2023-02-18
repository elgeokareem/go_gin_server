package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func Routes() {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	getRoutes(router)
	router.Run()
}

// Get all the routes
func getRoutes(router *gin.Engine) {
	User(router)
	Users(router)
	LoginUser(router)
}
