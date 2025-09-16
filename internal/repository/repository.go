// Package repository определяет слой доступа к данным (DAL) приложения.
package repository

import (
	"context"

	"lk/internal/models"

	"github.com/jmoiron/sqlx"
)

// Transactor определяет интерфейс для управления транзакциями.
type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

// UserRepository определяет методы для работы с пользователями.
type UserRepository interface {
	CreateUser(ctx context.Context, tx *sqlx.Tx, user models.User) (uint64, error)
	GetUserByPhone(ctx context.Context, phone string) (models.User, error)
	GetUserByID(ctx context.Context, id uint64) (models.User, error)
	CreateUserProfile(ctx context.Context, tx *sqlx.Tx, profile models.UserProfile) (uint64, error)
	GetUserProfileByUserID(ctx context.Context, userID uint64) (models.UserProfile, error)
}

// DoctorRepository определяет методы для работы с врачами.
type DoctorRepository interface {
	GetDoctorByID(ctx context.Context, id uint64) (models.Doctor, error)
	GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32) ([]models.Doctor, error)
	SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error)
}

// AppointmentRepository определяет методы для работы с записями на прием.
type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetAppointmentsByUserID(ctx context.Context, userID uint64) ([]models.Appointment, error)
}

// DirectoryRepository определяет методы для работы со справочниками.
type DirectoryRepository interface {
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	GetAllSpecialties(ctx context.Context) ([]models.Specialty, error)
}

// Repository - это контейнер для всех репозиториев приложения.
type Repository struct {
	User        UserRepository
	Doctor      DoctorRepository
	Appointment AppointmentRepository
	Directory   DirectoryRepository
	Transactor
}

// NewRepository создает новый экземпляр главного репозитория.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserPostgres(db),
		Doctor:      NewDoctorPostgres(db),
		Appointment: NewAppointmentPostgres(db),
		Directory:   NewDirectoryPostgres(db),
		Transactor:  NewTransactor(db),
	}
}
