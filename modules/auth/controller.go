package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData LOGIN
		err := c.BindJSON(&loginData)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error in service"})
			c.Abort()
			return
		}

		user, isUserInDb := CheckIfUserIsRegistered(loginData.EMAIL)

		if !isUserInDb {
			c.JSON(http.StatusNotFound, gin.H{"status": "client not in DB"})
			c.Abort()
			return
		}

		// Check if password matches
		matchPassword := DoPasswordsMatch(user.Password, loginData.PASSWORD)

		if !matchPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "password doesn't match"})
			c.Abort()
			return
		}

		token, err := GenerateJWT(user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in service"})
			c.Abort()
			return
		}

		completeToken := "Bearer " + token

		// create a cookie to send the token using the c.SetCookie() syntax
		// the expiration time is set to 24 hours
		// the cookie is only accessible by the server
		// the cookie is only accessible through the HTTP protocol
		// the last parameter is set to true because we're using HTTPS
		// the last parameter is set to true because we're using HTTPS
		c.SetCookie("token", completeToken, int(3600), "/", "127.0.0.1", false, false)

		c.JSON(http.StatusOK, gin.H{"status": "client logged in successfully"})
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerData REGISTER
		err := c.BindJSON(&registerData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Check user doesn't exists already
		_, isUserInDb := CheckIfUserIsRegistered(registerData.EMAIL)

		if isUserInDb {
			c.JSON(http.StatusConflict, gin.H{"status": "client already registered"})
			c.Abort()
			return
		}

		// Save the user to DB
		err = RegisterUserService(registerData.EMAIL, registerData.PASSWORD)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "client registered successfully"})
	}
}
