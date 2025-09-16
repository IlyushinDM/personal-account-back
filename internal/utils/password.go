// Package utils предоставляет набор вспомогательных функций и утилит,
// которые используются во всем приложении, но не содержат бизнес-логики.
// Включает в себя логирование, хэширование паролей и управление токенами.
package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword создает хэш из строки пароля с использованием bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash сравнивает пароль в виде строки с его хэшем.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
