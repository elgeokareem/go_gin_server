package services

import (
	"fmt"
	"goGinServer/utils"
)

func RegisterUserService(email string, password string) {
	// Hash the password with salt
	hashedPassword := utils.HashPassword(password)
	fmt.Println("hashed password", hashedPassword)
}
