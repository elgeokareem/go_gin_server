package middlewares

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/models"
	"goGinServer/modules/auth"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header["Authorization"]

		if len(authHeader) == 0 {
			c.JSON(http.StatusBadRequest, "No token was found")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader[0], " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Wrong header format"},
			)
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ParseAndValidateToken(tokenString)
		if err != nil {
			// Check for specific JWT errors if needed
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed token"})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired or not active yet"})
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Couldn't handle this token"})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		userIDUint64, err := strconv.ParseUint(claims.ID, 10, 64)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Error processing user ID from token"},
			)
			c.Abort()
			return
		}

		userID := uint(userIDUint64)
		var user models.User

		result := db.Service.DB.First(&user, userID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// User with this ID was not found in the database
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				// Some other database error occurred
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": "Database error retrieving user"},
				)
			}
			c.Abort()
			return
		}
		c.Set(auth.UserContextKey, user)
		c.Next()
	}
}
