package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

// GenerateToken создает новый JWT для указанного ID пользователя.
func GenerateToken(userID uint64, secretKey string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
		"sub": userID,
	})

	return token.SignedString([]byte(secretKey))
}

// ParseToken проверяет JWT и возвращает ID пользователя, который в нем содержится.
func ParseToken(accessToken string, secretKey string) (uint64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	subFloat, ok := (*claims)["sub"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	userID := uint64(subFloat)
	if userID == 0 {
		return 0, ErrInvalidToken
	}

	return userID, nil
}
