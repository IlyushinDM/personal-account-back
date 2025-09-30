// Package utils предоставляет набор вспомогательных функций и утилит,
// которые используются во всем приложении, но не содержат бизнес-логики.
// Включает в себя логирование, хэширование паролей и управление токенами.
package utils

import (
	"fmt"
	"strings"
	"unicode"

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

// PasswordValidationError содержит детали ошибки валидации пароля.
type PasswordValidationError struct {
	Message string
	Details []string
}

func (e *PasswordValidationError) Error() string {
	return e.Message
}

// ValidatePasswordStrength проверяет сложность пароля.
func ValidatePasswordStrength(password string) error {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, "минимум 8 символов")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "заглавные буквы")
	}
	if !hasLower {
		errors = append(errors, "строчные буквы")
	}
	if !hasDigit {
		errors = append(errors, "цифры")
	}
	if !hasSpecial {
		errors = append(errors, "специальные символы")
	}

	if len(errors) > 0 {
		return &PasswordValidationError{
			Message: "пароль должен содержать: " + strings.Join(errors, ", "),
			Details: errors,
		}
	}

	return nil
}
