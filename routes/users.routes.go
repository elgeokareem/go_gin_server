package routes

import (
	"goGinServer/controllers"

	"github.com/gin-gonic/gin"
)

func LoginUser(router *gin.Engine) {
	users := router.Group("/auth")

	users.POST("/login", controllers.Login())
	users.POST("/register", controllers.Register())
}
