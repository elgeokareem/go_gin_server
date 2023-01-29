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

		token := c.Request.Header["Authorization"]
		tokenParsed, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

		fmt.Println(tokenParsed)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "Error")
		}

		c.JSON(http.StatusOK, "todo bien")
	}
}
