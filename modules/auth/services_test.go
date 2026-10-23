package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"goGinServer/db/models"
)

func TestParseAndValidateToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	user := models.User{ID: 42, Email: "test@example.com"}
	tokenString, err := GenerateJWT(user)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	claims, err := ParseAndValidateToken(tokenString)
	if err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
	if claims.ID != "42" || claims.Email != user.Email {
		t.Fatalf("unexpected claims: %+v", claims)
	}

	if _, err := ParseAndValidateToken("not-a-jwt"); !errors.Is(err, jwt.ErrTokenMalformed) {
		t.Fatalf("malformed token: want ErrTokenMalformed, got %v", err)
	}

	t.Setenv("JWT_SECRET", "other-secret")
	if _, err := ParseAndValidateToken(tokenString); !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Fatalf("wrong signature: want ErrTokenSignatureInvalid, got %v", err)
	}

	t.Setenv("JWT_SECRET", "test-secret")
	expired := MyCustomClaims{
		ID: "42",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expired).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	if _, err := ParseAndValidateToken(expiredString); !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("expired token: want ErrTokenExpired, got %v", err)
	}
}
