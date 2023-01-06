package routes

import (
	"goGinServer/controllers"

	"github.com/gin-gonic/gin"
)

func Users(router *gin.Engine) {
	users := router.Group("/sql")

	users.GET("/groupby", controllers.SqlGroupByController())
}
