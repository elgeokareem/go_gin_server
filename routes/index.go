package routes

import (
	"github.com/gin-gonic/gin"
)

func Routes() {
	router := gin.Default()

	getRoutes(router)
	router.Run()
}

// Get all the
func getRoutes(router *gin.Engine) {
	User(router)
}
