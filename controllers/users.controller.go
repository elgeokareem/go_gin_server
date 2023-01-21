package controllers

import (
	"fmt"
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header["Authorization"]
		tokenParsed, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
			return []byte("your-256-bit-secret"), nil
		})

		fmt.Println(tokenParsed)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "Error")
		}

		// hash := sha256.New()

		// we can use a service here.
		data := services.RegisterUserService()

		c.JSON(http.StatusOK, data)
	}
}

type LOGIN struct {
	EMAIL    string `json:"email" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData LOGIN
		err := c.BindJSON(&loginData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in service"})
		}

		c.JSON(http.StatusOK, gin.H{"status": loginData.EMAIL})
	}
}
