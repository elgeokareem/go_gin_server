package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetJwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenSecret := os.Getenv("JWT_SECRET")

		tokenString := c.Request.Header["Authorization"]
		fmt.Println("tokenString: ", c.Request.Header)

		if len(tokenString) == 0 {
			c.JSON(http.StatusBadRequest, "No token was found")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[0], func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

		if err != nil {
			fmt.Println("token: ", token)
			fmt.Println("token tokenString: ", tokenString)
			fmt.Println("error:", err)
			c.JSON(http.StatusInternalServerError, "Error")
			c.Abort()
			return
		}

		c.Next()
	}
}
