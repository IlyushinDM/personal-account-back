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

var (
	ErrUserExists          = errors.New("user with this phone already exists")
	ErrInvalidCredentials  = errors.New("invalid phone or password")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrCodeMismatch        = errors.New("confirmation code is incorrect or expired")
)

const (
	refreshTokenTTL = 72 * time.Hour
	resetCodeTTL    = 5 * time.Minute
	resetCodePrefix = "reset_code:"
)

type resetCodeInfo struct {
	Code      string
	ExpiresAt time.Time
}

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
	// 1. Проверяем, не существует ли уже пользователь с таким телефоном.
	_, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err == nil {
		return nil, ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("ERROR: database error while checking user existence for phone %s: %v", phone, err)
		return nil, fmt.Errorf("database error while checking user: %w", err)
	}

	// 2. Хэшируем пароль.
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 3. Готовим данные для профиля.
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
		return nil, fmt.Errorf("invalid birth date format (expected YYYY-MM-DD): %w", err)
	}

	var userID uint64
	// 4. Выполняем создание пользователя и профиля в одной транзакции.
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
		return nil, err
	}

	// 5. Создаем сессию (пару токенов) для нового пользователя
	return s.createSession(ctx, userID)
}

// GenerateToken - бизнес-логика входа пользователя. Возвращает пару токенов.
func (s *authService) GenerateToken(ctx context.Context, phone, password string) (map[string]string, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.createSession(ctx, user.ID)
}

// RefreshToken обновляет пару токенов, используя старый refresh-токен.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	// Получаем userID из самого токена (предполагаем формат "userID.randomString")
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidRefreshToken
	}
	userID, err := utils.ParseUserID(parts[0])
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Получаем токен из БД
	storedToken, err := s.tokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Сравниваем хеши
	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	if base64.StdEncoding.EncodeToString(refreshTokenHash[:]) != storedToken.TokenHash {
		return nil, ErrInvalidRefreshToken
	}

	// Создаем новую сессию
	return s.createSession(ctx, userID)
}

// Logout завершает сессию пользователя, удаляя refresh-токен из БД.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 2 {
		return ErrInvalidRefreshToken
	}
	userID, err := utils.ParseUserID(parts[0])
	if err != nil {
		return ErrInvalidRefreshToken
	}

	// Дополнительная проверка, что токен валиден перед удалением
	storedToken, err := s.tokenRepo.GetByUserID(ctx, userID)
	if err != nil {
		// Если токена и так нет, считаем, что пользователь уже вышел
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return ErrInvalidRefreshToken
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	if base64.StdEncoding.EncodeToString(refreshTokenHash[:]) != storedToken.TokenHash {
		return ErrInvalidRefreshToken
	}

	return s.tokenRepo.Delete(ctx, userID)
}

// ForgotPassword инициирует сброс пароля.
func (s *authService) ForgotPassword(ctx context.Context, phone string) error {
	if _, err := s.userRepo.GetUserByPhone(ctx, phone); err != nil {
		log.Printf("INFO: Password reset requested for non-existent phone: %s", phone)
		return nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	key := resetCodePrefix + phone

	// Сохраняем код в Redis
	if err := s.cacheRepo.Set(ctx, key, code, resetCodeTTL); err != nil {
		log.Printf("ERROR: Failed to set reset code to redis for phone %s: %v", phone, err)
		return fmt.Errorf("internal server error")
	}

	// Имитируем отправку SMS
	log.Printf("!!! MOCK SMS !!! Password reset code for user %s is: %s", phone, code)

	return nil
}

// ResetPassword устанавливает новый пароль с использованием кода.
func (s *authService) ResetPassword(ctx context.Context, phone, code, newPassword string) error {
	key := resetCodePrefix + phone

	// Получаем код из Redis
	storedCode, err := s.cacheRepo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrCodeMismatch
		}
		log.Printf("ERROR: Failed to get reset code from redis for phone %s: %v", phone, err)
		return fmt.Errorf("internal server error")
	}

	if storedCode != code {
		return ErrCodeMismatch
	}

	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		return ErrInvalidCredentials
	}

	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, newPasswordHash); err != nil {
		return err
	}

	// Удаляем использованный код из Redis
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
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	subFloat, ok := (*claims)["sub"].(float64)
	if !ok {
		return 0, errors.New("invalid subject claim")
	}

	return uint64(subFloat), nil
}

// createSession - внутренний метод для генерации и сохранения пары токенов.
func (s *authService) createSession(ctx context.Context, userID uint64) (map[string]string, error) {
	// Генерируем access token
	accessToken, err := utils.GenerateToken(userID, s.signingKey, s.tokenTTL)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh token
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	// Формат: "userID.randomString"
	refreshTokenString := fmt.Sprintf("%d.%s", userID, base64.URLEncoding.EncodeToString(randomBytes))

	// Хешируем refresh token для хранения в БД
	refreshTokenHash := sha256.Sum256([]byte(refreshTokenString))
	refreshTokenHashString := base64.StdEncoding.EncodeToString(refreshTokenHash[:])

	// Сохраняем в БД
	err = s.tokenRepo.Create(ctx, models.RefreshToken{
		UserID:    userID,
		TokenHash: refreshTokenHashString,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshTokenString,
	}, nil
}
