package services

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/model"
	"goGinServer/utils"

	"gorm.io/gorm"
)

func CheckIfUserIsRegistered(email string) (model.User, bool) {
	db := db.DbConnection()
	user := model.User{Email: email}

	// userDb := db.Where(user).First(&user)
	userDb := db.Where("email = ?", email).First(&user)

	return user, !errors.Is(userDb.Error, gorm.ErrRecordNotFound)
}

func RegisterUserService(email string, password string) {
	// Hash the password with salt
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		panic(err)
	}

	user := model.User{Email: email, Password: hashedPassword}

	// Add user data to DB
	db := db.DbConnection()
	db.Create(&user)
}
