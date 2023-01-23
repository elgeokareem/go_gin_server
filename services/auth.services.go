package services

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func GenerateJWT(user model.User) (string, error) {
	type MyCustomClaims struct {
		User  string `json:"user"`
		Email string `json:"email"`
		jwt.RegisteredClaims
	}

	claims := MyCustomClaims{
		user.Email,
		user.Email,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return ss, err
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
