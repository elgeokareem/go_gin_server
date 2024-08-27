package routes

import (
	"goGinServer/controllers"

	"github.com/gin-gonic/gin"
)

func CronEndpoints(router *gin.Engine) {
	cronRoutes := router.Group("/cron")

	cronRoutes.GET("/save-funds-wallet", controllers.CronSaveFundData())
}
