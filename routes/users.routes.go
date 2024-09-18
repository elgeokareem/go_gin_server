package routes

import (
	"goGinServer/modules/auth"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.Engine) {
	users := router.Group("/auth")

	users.POST("/login", auth.Login())
	users.POST("/register", auth.Register())
}
