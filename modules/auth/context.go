package auth

import (
	"goGinServer/db/models"

	"github.com/gin-gonic/gin"
)

const UserContextKey = "User"

func GetUserFromContext(c *gin.Context) models.User {
	// No validations here because they are done in AuthMiddleware
	user, _ := c.Get(UserContextKey)
	u := user.(models.User)
	return u
}
