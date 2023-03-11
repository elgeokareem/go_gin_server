package routes

import (
	"goGinServer/controllers"

	"github.com/gin-gonic/gin"
)

func Dashboard(router *gin.Engine) {
	router.GET("/", controllers.DashboardController())
}
