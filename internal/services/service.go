// Package services содержит бизнес-логику приложения.
package services

import (
	"context"
	"mime/multipart"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
)

// Authorization определяет методы для регистрации и входа пользователя.
type Authorization interface {
	// CreateUser создает нового пользователя и возвращает пару access/refresh токенов.
	CreateUser(ctx context.Context, phone, password, fullName,
		gender, birthDateStr string, cityID uint32) (map[string]string, error)
	// GenerateToken проверяет учетные данные и возвращает пару access/refresh токенов.
	GenerateToken(ctx context.Context, phone, password string) (map[string]string, error)
	// ParseToken проверяет access-токен и возвращает ID пользователя.
	ParseToken(token string) (uint64, error)
	// RefreshToken принимает валидный refresh-токен и возвращает новую пару токенов.
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
	// Logout инвалидирует refresh-токен на сервере.
	Logout(ctx context.Context, refreshToken string) error
	// ForgotPassword инициирует процесс восстановления пароля (например, отправляет код).
	ForgotPassword(ctx context.Context, phone string) error
	// ResetPassword сбрасывает пароль пользователя с использованием кода подтверждения.
	ResetPassword(ctx context.Context, phone, code, newPassword string) error
}

// UserService определяет методы для работы с данными пользователя.
type UserService interface {
	GetFullUserProfile(ctx context.Context, userID uint64) (models.UserProfile, []models.Appointment, error)
	UpdateUserProfile(ctx context.Context, userID uint64, input models.UserProfile) (models.UserProfile, error)
	UpdateAvatar(ctx context.Context, userID uint64, fileHeader *multipart.FileHeader) (string, error)
}

// DoctorService определяет методы для работы с информацией о врачах.
type DoctorService interface {
	GetDoctorByID(ctx context.Context, doctorID uint64) (models.Doctor, error)
	GetDoctorsBySpecialty(ctx context.Context, specialtyID uint32, params models.PaginationParams) (
		models.PaginatedDoctorsResponse, error)
	SearchDoctors(ctx context.Context, query string) ([]models.Doctor, error)
	SearchDoctorsByService(ctx context.Context, serviceQuery string) ([]models.Doctor, error)
	GetSpecialistRecommendations(ctx context.Context, doctorID uint64) (models.Recommendation, error)
}

// DirectoryService для работы со справочниками.
type DirectoryService interface {
	GetAllDepartmentsWithSpecialties(ctx context.Context) ([]models.DepartmentWithSpecialties, error)
	GetSpecialties(ctx context.Context, departmentID *uint32) ([]models.Specialty, error)
}

// AppointmentService определяет методы для работы с записями на прием.
type AppointmentService interface {
	CreateAppointment(ctx context.Context, appointment models.Appointment) (uint64, error)
	GetUserAppointments(ctx context.Context, userID uint64) ([]models.Appointment, error)
	CancelAppointment(ctx context.Context, userID, appointmentID uint64) error
	GetAvailableDates(ctx context.Context, doctorID, serviceID uint64, month string) (
		models.AvailableDatesResponse, error)
	GetAvailableSlots(ctx context.Context, doctorID, serviceID uint64, date string) (
		models.AvailableSlotsResponse, error)
	GetUpcomingForUser(ctx context.Context, userID uint64) ([]models.Appointment, error)
}

// InfoService определяет методы для работы с общей информацией.
type InfoService interface {
	GetServiceRecommendations(ctx context.Context, serviceID uint64) (models.Recommendation, error)
	GetClinicInfo(ctx context.Context) (models.ClinicInfo, error)
	GetLegalDocuments(ctx context.Context) ([]models.LegalDocument, error)
}

// PrescriptionService определяет методы для работы с назначениями.
type PrescriptionService interface {
	GetActiveForUser(ctx context.Context, userID uint64) ([]models.Prescription, error)
	ArchiveForUser(ctx context.Context, userID, prescriptionID uint64) error
}

// MedicalCardService определяет методы для работы с медкартой.
type MedicalCardService interface {
	GetVisits(ctx context.Context, userID uint64, params models.PaginationParams) (
		models.PaginatedVisitsResponse, error)
	GetAnalyses(ctx context.Context, userID uint64, status *string) ([]models.LabAnalysis, error)
	GetArchivedPrescriptions(ctx context.Context, userID uint64) ([]models.Prescription, error)
	GetSummary(ctx context.Context, userID uint64) (models.MedicalCardSummary, error)
	ArchivePrescription(ctx context.Context, userID, prescriptionID uint64) error
	DownloadFile(ctx context.Context, userID, fileID uint64) ([]byte, string, error)
}

// Service - это контейнер для всех сервисов приложения.
type Service struct {
	Authorization Authorization
	User          UserService
	Doctor        DoctorService
	Appointment   AppointmentService
	Directory     DirectoryService
	Info          InfoService
	Prescription  PrescriptionService
	MedicalCard   MedicalCardService
}

// ServiceDependencies содержит все зависимости, необходимые для создания сервисов.
type ServiceDependencies struct {
	Repos      *repository.Repository
	SigningKey string
	TokenTTL   time.Duration
}

// NewService создает новый экземпляр главного сервиса, инициализируя все реализации.
func NewService(deps ServiceDependencies) *Service {
	// Создаем сервисы, от которых могут зависеть другие
	prescriptionRepo := deps.Repos.Prescription

	// Создаем все сервисы с их зависимостями
	authService := NewAuthService(
		deps.Repos.User,
		deps.Repos.Token,
		deps.Repos.Cache,
		deps.Repos.Transactor,
		deps.SigningKey,
		deps.TokenTTL,
	)
	prescriptionService := NewPrescriptionService(prescriptionRepo)

	return &Service{
		Authorization: authService,
		User:          NewUserService(deps.Repos.User, deps.Repos.Appointment),
		Doctor:        NewDoctorService(deps.Repos.Doctor),
		Appointment:   NewAppointmentService(deps.Repos.Appointment, deps.Repos.Doctor),
		Directory:     NewDirectoryService(deps.Repos.Directory),
		Info:          NewInfoService(deps.Repos.Service, deps.Repos.Info),
		Prescription:  prescriptionService,
		MedicalCard:   NewMedicalCardService(deps.Repos.MedicalCard, prescriptionRepo),
	}
}
