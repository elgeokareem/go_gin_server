package routes

import (
	"goGinServer/controllers"
	"goGinServer/middlewares"

	"github.com/gin-gonic/gin"
)

func LoginUser(router *gin.Engine) {
	users := router.Group("/auth")

	users.POST("/login", controllers.Login())
	users.POST("/register", middlewares.GetJwtToken(), controllers.Register())
}
