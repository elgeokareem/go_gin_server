package services

import (
	"goGinServer/db"
	"goGinServer/db/model"
	"goGinServer/utils"
)

func RegisterUserService(email string, password string) {
	// Hash the password with salt
	hashedPassword := utils.HashPassword(password)

	user := model.User{Email: email, Password: hashedPassword}

	// Add user data to DB
	db := db.DbConnection()
	db.Create(&user)
}
