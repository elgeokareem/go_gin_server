package auth

import (
	"errors"
	"goGinServer/db"
	"goGinServer/db/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Hash password
func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	passwordBytes := []byte(password)

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

type MyCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {
	claims := MyCustomClaims{
		strconv.FormatUint(uint64(user.ID), 10),
		user.Email,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return ss, err
}

func ParseAndValidateToken(tokenString string) (*MyCustomClaims, error) {
	claims := &MyCustomClaims{}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, jwt.NewValidationError(
			"JWT_SECRET not configured",
			jwt.ValidationErrorUnverifiable,
		)
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError(
					"unexpected signing method",
					jwt.ValidationErrorSignatureInvalid,
				)
			}
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return nil, err // Handles expired tokens, malformed tokens, etc.
	}

	if !token.Valid {
		return nil, jwt.NewValidationError(
			"invalid token",
			jwt.ValidationErrorClaimsInvalid,
		)
	}

	return claims, nil
}

func CheckIfUserIsRegistered(email string) (models.User, bool) {
	user := models.User{Email: email}

	// userDb := db.Where(user).First(&user)
	userDb := db.Service.DB.Where("email = ?", email).First(&user)

	return user, !errors.Is(userDb.Error, gorm.ErrRecordNotFound)
}

func RegisterUserService(email string, password string) error {
	// Hash the password with salt
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{Email: email, Password: hashedPassword}

	// Add user data to DB
	result := db.Service.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
