package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitServer() {
	router := gin.Default()
	// add cors
	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://locahost:5173"}
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	getRoutes(router)
	router.Run()
}

// Get all the routes
func getRoutes(router *gin.Engine) {
	Dashboard(router) // The dashboard should be for aggregating the data from the db, not making calls for current data
	CronEndpoints(router)
	Users(router)
	Auth(router)
}
