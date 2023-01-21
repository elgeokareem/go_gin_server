package services

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/model"
	"goGinServer/utils"

	"gorm.io/gorm"
)

func CheckIfUserIsRegistered(email string) bool {
	db := db.DbConnection()
	user := model.User{Email: email}

	userDb := db.Limit(1).Find(&user)

	return !errors.Is(userDb.Error, gorm.ErrRecordNotFound)
}

func RegisterUserService(email string, password string) {
	// Hash the password with salt
	hashedPassword := utils.HashPassword(password)

	user := model.User{Email: email, Password: hashedPassword}

	// Add user data to DB
	db := db.DbConnection()
	db.Create(&user)
}
