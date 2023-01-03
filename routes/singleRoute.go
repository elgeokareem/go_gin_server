package routes

import (
	"goGinServer/controllers"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine) {
	router.GET("/", controllers.PongController())
}
