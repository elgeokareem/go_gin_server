package services

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Hash password
func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func CheckIfUserIsRegistered(email string) (model.User, bool) {
	db := db.DbConnection()
	user := model.User{Email: email}

	// userDb := db.Where(user).First(&user)
	userDb := db.Where("email = ?", email).First(&user)

	return user, !errors.Is(userDb.Error, gorm.ErrRecordNotFound)
}

func RegisterUserService(email string, password string) {
	// Hash the password with salt
	hashedPassword, err := HashPassword(password)

	if err != nil {
		panic(err)
	}

	user := model.User{Email: email, Password: hashedPassword}

	// Add user data to DB
	db := db.DbConnection()
	db.Create(&user)
}
