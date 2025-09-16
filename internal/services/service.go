// Package services содержит бизнес-логику приложения.
package services

import (
	"context"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
)

// Authorization определяет методы для регистрации и входа пользователя.
type Authorization interface {
	CreateUser(ctx context.Context, phone, password, fullName,
		gender, birthDateStr string, cityID uint32) (uint64, error)
	GenerateToken(ctx context.Context, phone, password string) (string, error)
	ParseToken(token string) (uint64, error)
}

// UserService определяет методы для работы с данными пользователя.
type UserService interface {
	GetFullUserProfile(ctx context.Context, userID uint64) (models.UserProfile, []models.Appointment, error)
}

// DoctorService определяет методы для работы с информацией о врачах.
type DoctorService interface {
	GetDoctorByID(ctx context.Context, doctorID uint64) (models.Doctor, error)
	GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32) ([]models.Doctor, error)
	SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error)
}

// AppointmentService определяет методы для работы с записями на прием.
type AppointmentService interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error)
}

// DirectoryService для работы со справочниками.
type DirectoryService interface {
	GetAllDepartmentsWithSpecialties(ctx context.Context) ([]models.DepartmentWithSpecialties, error)
}

// Service - это контейнер для всех сервисов приложения.
type Service struct {
	Authorization Authorization
	User          UserService
	Doctor        DoctorService
	Appointment   AppointmentService
	Directory     DirectoryService
}

// ServiceDependencies содержит все зависимости, необходимые для создания сервисов.
type ServiceDependencies struct {
	Repos      *repository.Repository
	SigningKey string
	TokenTTL   time.Duration
}

// NewService создает новый экземпляр главного сервиса, инициализируя все реализации.
func NewService(deps ServiceDependencies) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repos.User, deps.Repos.Transactor,
			deps.SigningKey, deps.TokenTTL),
		User:        NewUserService(deps.Repos.User, deps.Repos.Appointment),
		Doctor:      NewDoctorService(deps.Repos.Doctor),
		Appointment: NewAppointmentService(deps.Repos.Appointment),
		Directory:   NewDirectoryService(deps.Repos.Directory),
	}
}
