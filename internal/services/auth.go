package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrUserExists         = errors.New("user with this phone already exists")
	ErrInvalidCredentials = errors.New("invalid phone or password")
)

// authService - это конкретная реализация интерфейса Authorization.
type authService struct {
	userRepo   repository.UserRepository
	transactor repository.Transactor
	signingKey string
	tokenTTL   time.Duration
}

// NewAuthService является конструктором для сервиса авторизации.
func NewAuthService(userRepo repository.UserRepository, transactor repository.Transactor,
	signingKey string, tokenTTL time.Duration,
) Authorization {
	return &authService{
		userRepo:   userRepo,
		transactor: transactor,
		signingKey: signingKey,
		tokenTTL:   tokenTTL,
	}
}

// CreateUser - бизнес-логика регистрации нового пользователя.
func (s *authService) CreateUser(ctx context.Context, phone, password, fullName, gender,
	birthDateStr string, cityID uint32,
) (uint64, error) {
	// 1. Проверяем, не существует ли уже пользователь с таким телефоном.
	_, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err == nil {
		return 0, ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("ERROR: database error while checking user existence for phone %s: %v", phone, err)
		return 0, fmt.Errorf("database error while checking user: %w", err)
	}

	// 2. Хэшируем пароль.
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
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
		return 0, fmt.Errorf("invalid birth date format (expected YYYY-MM-DD): %w", err)
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
		return 0, err
	}

	return userID, nil
}

// GenerateToken - бизнес-логика входа пользователя.
func (s *authService) GenerateToken(ctx context.Context, phone, password string) (string, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := utils.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	return utils.GenerateToken(user.ID, s.signingKey, s.tokenTTL)
}

// ParseToken проверяет токен и возвращает ID пользователя из него.
func (s *authService) ParseToken(accessToken string) (uint64, error) {
	return utils.ParseToken(accessToken, s.signingKey)
}
