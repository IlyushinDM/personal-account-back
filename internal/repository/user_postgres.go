package repository

import (
	"context"
	"fmt"
	"strings"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// UserPostgres реализует UserRepository для PostgreSQL.
type UserPostgres struct {
	db *sqlx.DB
}

// NewUserPostgres создает новый экземпляр репозитория для пользователей.
func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// CreateUser вставляет нового пользователя в таблицу users и возвращает его ID.
// * Эта функция должна вызываться внутри транзакции.
func (r *UserPostgres) CreateUser(ctx context.Context, tx *sqlx.Tx, user models.User) (uint64, error) {
	var id uint64
	query := "INSERT INTO medical_center.users (phone, password_hash) VALUES ($1, $2) RETURNING id"
	err := tx.QueryRowxContext(ctx, query, user.Phone, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateUserProfile вставляет новый профиль пользователя и возвращает его ID.
// * Эта функция должна вызываться внутри транзакции.
func (r *UserPostgres) CreateUserProfile(ctx context.Context, tx *sqlx.Tx, profile models.UserProfile) (uint64, error) {
	var id uint64
	query := `INSERT INTO medical_center.user_profiles
        (user_id, first_name, last_name, patronymic, birth_date, gender, city_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := tx.QueryRowxContext(ctx, query, profile.UserID, profile.FirstName, profile.LastName, profile.Patronymic,
		profile.BirthDate, profile.Gender, profile.CityID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByPhone находит пользователя по номеру телефона.
func (r *UserPostgres) GetUserByPhone(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM medical_center.users WHERE phone=$1"
	err := r.db.GetContext(ctx, &user, query, phone)
	return user, err
}

// GetUserByID находит пользователя по его ID.
func (r *UserPostgres) GetUserByID(ctx context.Context, id uint64) (models.User, error) {
	var user models.User
	query := "SELECT * FROM medical_center.users WHERE id=$1"
	err := r.db.GetContext(ctx, &user, query, id)
	return user, err
}

// GetUserProfileByUserID находит профиль пользователя по ID пользователя.
func (r *UserPostgres) GetUserProfileByUserID(ctx context.Context, userID uint64) (models.UserProfile, error) {
	var profile models.UserProfile
	query := "SELECT * FROM medical_center.user_profiles WHERE user_id=$1"
	err := r.db.GetContext(ctx, &profile, query, userID)
	return profile, err
}

// UpdateUserProfile обновляет данные профиля пользователя.
func (r *UserPostgres) UpdateUserProfile(ctx context.Context, profile models.UserProfile) (models.UserProfile, error) {
	var updatedProfile models.UserProfile

	// Формируем запрос динамически, чтобы обновлять только переданные поля
	queryParts := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if profile.Email.Valid {
		queryParts = append(queryParts, fmt.Sprintf("email=$%d", argID))
		args = append(args, profile.Email)
		argID++
	}
	if profile.CityID > 0 {
		queryParts = append(queryParts, fmt.Sprintf("city_id=$%d", argID))
		args = append(args, profile.CityID)
		argID++
	}

	if len(queryParts) == 0 {
		// Если нечего обновлять, просто возвращаем текущий профиль
		return r.GetUserProfileByUserID(ctx, profile.UserID)
	}

	query := fmt.Sprintf("UPDATE medical_center.user_profiles SET %s WHERE user_id=$%d RETURNING *",
		strings.Join(queryParts, ", "), argID)
	args = append(args, profile.UserID)

	err := r.db.GetContext(ctx, &updatedProfile, query, args...)
	return updatedProfile, err
}

// UpdateAvatar обновляет URL аватара для профиля пользователя.
func (r *UserPostgres) UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error {
	query := `UPDATE medical_center.user_profiles SET avatar_url = $1 WHERE user_id = $2`
	_, err := r.db.ExecContext(ctx, query, avatarURL, userID)
	return err
}
