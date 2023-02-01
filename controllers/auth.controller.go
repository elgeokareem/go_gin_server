package controllers

import (
	"goGinServer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LOGIN struct {
	EMAIL    string `json:"email" binding:"required"`
	PASSWORD string `json:"password" binding:"required"`
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData LOGIN
		err := c.BindJSON(&loginData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in service"})
			return
		}

		user, isUserInDb := services.CheckIfUserIsRegistered(loginData.EMAIL)

		if !isUserInDb {
			c.JSON(http.StatusNotFound, gin.H{"status": "client not in DB"})
			return
		}

		// Check if password matches
		matchPassword := services.DoPasswordsMatch(user.Password, loginData.PASSWORD)

		if !matchPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "password doesn't match"})
			return
		}

		token, err := services.GenerateJWT(user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in service"})
			return
		}

		c.SetCookie("gin_cookie", token, 3600, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"status": "client logged in successfully"})
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData LOGIN
		err := c.BindJSON(&loginData)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in service"})
		}

		// Check user doesn't exists already
		_, isUserInDb := services.CheckIfUserIsRegistered(loginData.EMAIL)

		if isUserInDb {
			c.JSON(http.StatusConflict, gin.H{"status": "client already registered"})
			return
		}

		// Save the user to DB
		services.RegisterUserService(loginData.EMAIL, loginData.PASSWORD)

		c.JSON(http.StatusCreated, gin.H{"status": "client registered successfully"})
	}
}
