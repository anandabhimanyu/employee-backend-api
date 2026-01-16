package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret   string
	duration time.Duration
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret:   secret,
		duration: 24 * time.Hour, // ✅ token valid for 24 hours
	}
}

// Generate creates JWT with user_id + role
func (j *JWTManager) Generate(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(j.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// Verify validates token and returns user_id
func (j *JWTManager) Verify(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id")
	}

	return int(userID), nil
}

// Secret exposes JWT secret for middleware
func (j *JWTManager) Secret() string {
	return j.secret
}
