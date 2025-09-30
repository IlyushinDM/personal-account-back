package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

const (
	refreshTokenTTL = 72 * time.Hour
	resetCodeTTL    = 5 * time.Minute
	resetCodePrefix = "reset_code:"
)

// authService - это конкретная реализация интерфейса Authorization.
type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	cacheRepo  repository.CacheRepository
	transactor repository.Transactor
	signingKey string
	tokenTTL   time.Duration
}

// NewAuthService является конструктором для сервиса авторизации.
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	cacheRepo repository.CacheRepository,
	transactor repository.Transactor,
	signingKey string,
	tokenTTL time.Duration,
) Authorization {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		cacheRepo:  cacheRepo,
		transactor: transactor,
		signingKey: signingKey,
		tokenTTL:   tokenTTL,
	}
}

// CreateUser - бизнес-логика регистрации нового пользователя. Возвращает пару токенов.
func (s *authService) CreateUser(ctx context.Context, phone, password, fullName, gender,
	birthDateStr string, cityID uint32,
) (map[string]string, error) {
	_, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err == nil {
		return nil, NewConflictError("user with this phone already exists", nil)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, NewInternalServerError("database error while checking user", err)
	}

	if err := utils.ValidatePasswordStrength(password); err != nil {
		return nil, NewBadRequestError(err.Error(), err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, NewInternalServerError("failed to hash password", err)
	}

	nameParts := strings.Fields(fullName)
	var lastName, firstName, patronymic string
	if len(nameParts) > 0 {
		lastName = nameParts[0]
	}
	if len(nameParts) > 1 {
		firstName = nameParts[1]
	}
	if len(nameParts) > 2 {
		patronymic = strings.Join(nameParts[2:], " ")
	}

	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return nil, NewBadRequestError("invalid birth date format (expected YYYY-MM-DD)", err)
	}

	var userID uint64
	err = s.transactor.WithinTransaction(ctx, func(tx *gorm.DB) error {
		user := models.User{
			Phone:        phone,
			PasswordHash: hashedPassword,
		}
		newUserID, err := s.userRepo.CreateUser(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		userID = newUserID

		profile := models.UserProfile{
			UserID:     userID,
			FirstName:  firstName,
			LastName:   lastName,
			Patronymic: sql.NullString{String: patronymic, Valid: patronymic != ""},
			BirthDate:  birthDate,
			Gender:     gender,
			CityID:     cityID,
		}
		_, err = s.userRepo.CreateUserProfile(ctx, tx, profile)
		if err != nil {
			return fmt.Errorf("failed to create user profile: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, NewInternalServerError("transaction failed on user creation", err)
	}

	return s.createSession(ctx, userID)
}

// GenerateToken - бизнес-логика входа пользователя. Возвращает пару токенов.
func (s *authService) GenerateToken(ctx context.Context, phone, password string) (map[string]string, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewUnauthorizedError("invalid phone or password", nil)
		}
		return nil, NewInternalServerError("database error while getting user", err)
	}

	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, NewUnauthorizedError("invalid phone or password", nil)
	}

	return s.createSession(ctx, user.ID)
}

// RefreshToken обновляет пару токенов, используя старый refresh-токен.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 2 {
		return nil, NewUnauthorizedError("invalid refresh token format", nil)
	}
	userID, err := utils.ParseUserID(parts[0])
	if err != nil {
		return nil, NewUnauthorizedError("invalid user id in refresh token", err)
	}

	storedToken, err := s.tokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, NewUnauthorizedError("refresh token not found or expired", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	if base64.StdEncoding.EncodeToString(refreshTokenHash[:]) != storedToken.TokenHash {
		return nil, NewUnauthorizedError("invalid refresh token", nil)
	}

	return s.createSession(ctx, userID)
}

// Logout завершает сессию пользователя, удаляя refresh-токен из БД.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 2 {
		return NewUnauthorizedError("invalid refresh token format", nil)
	}
	userID, err := utils.ParseUserID(parts[0])
	if err != nil {
		return NewUnauthorizedError("invalid user id in refresh token", err)
	}

	if err := s.tokenRepo.Delete(ctx, userID); err != nil {
		return NewInternalServerError("failed to logout", err)
	}
	return nil
}

// ForgotPassword инициирует сброс пароля.
func (s *authService) ForgotPassword(ctx context.Context, phone string) error {
	if _, err := s.userRepo.GetUserByPhone(ctx, phone); err != nil {
		log.Printf("INFO: Password reset requested for non-existent phone: %s", phone)
		return nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	key := resetCodePrefix + phone

	if err := s.cacheRepo.Set(ctx, key, code, resetCodeTTL); err != nil {
		return NewInternalServerError("failed to set reset code to cache", err)
	}

	log.Printf("!!! MOCK SMS !!! Password reset code for user %s is: %s", phone, code)
	return nil
}

// ResetPassword устанавливает новый пароль с использованием кода.
func (s *authService) ResetPassword(ctx context.Context, phone, code, newPassword string) error {
	key := resetCodePrefix + phone
	storedCode, err := s.cacheRepo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NewUnauthorizedError("confirmation code is incorrect or expired", err)
		}
		return NewInternalServerError("failed to get reset code from cache", err)
	}

	if storedCode != code {
		return NewUnauthorizedError("confirmation code is incorrect or expired", nil)
	}

	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return NewBadRequestError(err.Error(), err)
	}

	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		return NewInternalServerError("failed to get user for password reset", err)
	}

	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return NewInternalServerError("failed to hash new password", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, newPasswordHash); err != nil {
		return NewInternalServerError("failed to update password in db", err)
	}

	_ = s.cacheRepo.Delete(ctx, key)
	return nil
}

// ParseToken проверяет токен и возвращает ID пользователя из него.
func (s *authService) ParseToken(accessToken string) (uint64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.signingKey), nil
		})
	if err != nil {
		return 0, NewUnauthorizedError("invalid token", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, NewUnauthorizedError("invalid token claims", nil)
	}

	subFloat, ok := (*claims)["sub"].(float64)
	if !ok {
		return 0, NewUnauthorizedError("invalid subject claim in token", nil)
	}

	return uint64(subFloat), nil
}

// createSession - внутренний метод для генерации и сохранения пары токенов.
func (s *authService) createSession(ctx context.Context, userID uint64) (map[string]string, error) {
	accessToken, err := utils.GenerateToken(userID, s.signingKey, s.tokenTTL)
	if err != nil {
		return nil, NewInternalServerError("failed to generate access token", err)
	}

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, NewInternalServerError("failed to generate refresh token", err)
	}
	refreshTokenString := fmt.Sprintf("%d.%s", userID, base64.URLEncoding.EncodeToString(randomBytes))
	refreshTokenHash := sha256.Sum256([]byte(refreshTokenString))
	refreshTokenHashString := base64.StdEncoding.EncodeToString(refreshTokenHash[:])

	err = s.tokenRepo.Create(ctx, models.RefreshToken{
		UserID:    userID,
		TokenHash: refreshTokenHashString,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	})
	if err != nil {
		return nil, NewInternalServerError("failed to save refresh token", err)
	}

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshTokenString,
	}, nil
}
