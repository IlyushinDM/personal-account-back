// Package services содержит бизнес-логику приложения.
package services

import (
	"context"
	"mime/multipart"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/storage"
)

// Authorization определяет методы для регистрации и входа пользователя.
type Authorization interface {
	CreateUser(ctx context.Context, phone, password, fullName,
		gender, birthDateStr string, cityID uint32) (map[string]string, error)
	GenerateToken(ctx context.Context, phone, password string) (map[string]string, error)
	ParseToken(token string) (uint64, error)
	RefreshToken(ctx context.Context, refreshToken string) (map[string]string, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, phone string) error
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
	GetAvailableSlotsByRange(ctx context.Context, doctorID, serviceID uint64, startDate, endDate string) (
		models.AvailableRangeSlotsResponse, error)
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

// AdminService определяет все методы для администрирования системы.
type AdminService interface {
	// Auth & Dashboard
	Login(ctx context.Context, login, password string) (map[string]string, error)
	ParseAdminToken(token string) (uint64, error)
	GetDashboardStats(ctx context.Context) (models.AdminDashboardStats, error)

	// User
	GetAllUsers(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error)
	GetUserByID(ctx context.Context, userID uint64) (*models.User, *models.UserProfile, error)
	UpdateUser(ctx context.Context, userID uint64, input UpdateUserInput) error
	DeleteUser(ctx context.Context, userID uint64) error
	GetUserAppointments(ctx context.Context, userID uint64, params models.PaginationParams) (
		[]models.Appointment, int64, error)
	GetUserAnalyses(ctx context.Context, userID uint64, params models.PaginationParams) (
		[]models.LabAnalysis, int64, error)

	// Doctor
	GetAllSpecialists(ctx context.Context, params models.PaginationParams) ([]models.Doctor, int64, error)
	CreateSpecialist(ctx context.Context, input CreateDoctorInput) (uint64, error)
	GetSpecialistByID(ctx context.Context, doctorID uint64) (models.Doctor, error)
	UpdateSpecialist(ctx context.Context, doctorID uint64, input UpdateDoctorInput) error
	DeleteSpecialist(ctx context.Context, doctorID uint64) error
	GetSpecialistSchedule(ctx context.Context, doctorID uint64) ([]models.Schedule, error)
	UpdateSpecialistSchedule(ctx context.Context, doctorID uint64, input UpdateScheduleInput) error

	// Appointment
	GetAllAppointments(ctx context.Context, params models.PaginationParams, filters map[string]interface{}) (
		[]models.Appointment, int64, error)
	GetAppointmentStats(ctx context.Context) (map[string]int64, error)
	GetAppointmentDetails(ctx context.Context, appointmentID uint64) (models.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error
	DeleteAppointment(ctx context.Context, appointmentID uint64) error

	// Service & Department
	GetAllServices(ctx context.Context) ([]models.Service, error)
	CreateService(ctx context.Context, input CreateServiceInput) (uint64, error)
	UpdateService(ctx context.Context, serviceID uint64, input UpdateServiceInput) error
	DeleteService(ctx context.Context, serviceID uint64) error
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	CreateDepartment(ctx context.Context, input CreateDepartmentInput) (uint32, error)
	UpdateDepartment(ctx context.Context, departmentID uint32, input UpdateDepartmentInput) error
	DeleteDepartment(ctx context.Context, departmentID uint32) error

	// TODO: Реализовать другие методы бизнес-логики (Analyses, Prescriptions, Family, Settings, и т.д.)
}

// --- DTO для AdminService ---

type UpdateUserInput struct {
	Phone     *string `json:"phone"`
	IsActive  *bool   `json:"isActive"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
}

type CreateDoctorInput struct {
	FirstName       string  `json:"firstName" binding:"required"`
	LastName        string  `json:"lastName" binding:"required"`
	Patronymic      *string `json:"patronymic"`
	SpecialtyID     uint32  `json:"specialtyId" binding:"required"`
	ExperienceYears uint16  `json:"experienceYears" binding:"required"`
	Recommendations *string `json:"recommendations"`
}

type UpdateDoctorInput struct {
	FirstName       *string `json:"firstName"`
	LastName        *string `json:"lastName"`
	Patronymic      *string `json:"patronymic"`
	SpecialtyID     *uint32 `json:"specialtyId"`
	ExperienceYears *uint16 `json:"experienceYears"`
	Recommendations *string `json:"recommendations"`
}

type ScheduleItem struct {
	Date      string `json:"date" binding:"required"`      // YYYY-MM-DD
	StartTime string `json:"startTime" binding:"required"` // HH:MM
	EndTime   string `json:"endTime" binding:"required"`   // HH:MM
}

type UpdateScheduleInput struct {
	Schedules []ScheduleItem `json:"schedules"`
}

type CreateServiceInput struct {
	Name            string  `json:"name" binding:"required"`
	Price           float64 `json:"price" binding:"required"`
	DurationMinutes uint16  `json:"durationMinutes" binding:"required"`
	Description     *string `json:"description"`
	DoctorID        uint64  `json:"doctorId" binding:"required"`
	Recommendations *string `json:"recommendations"`
}

type UpdateServiceInput struct {
	Name            *string  `json:"name"`
	Price           *float64 `json:"price"`
	DurationMinutes *uint16  `json:"durationMinutes"`
	Description     *string  `json:"description"`
	DoctorID        *uint64  `json:"doctorId"`
	Recommendations *string  `json:"recommendations"`
}

type CreateDepartmentInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDepartmentInput struct {
	Name *string `json:"name"`
}

// --- Service Контейнер ---

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
	Admin         AdminService
}

// ServiceDependencies содержит все зависимости, необходимые для создания сервисов.
type ServiceDependencies struct {
	Repos      *repository.Repository
	Storage    storage.FileStorage
	Location   *time.Location
	SigningKey string
	TokenTTL   time.Duration
}

// NewService создает новый экземпляр главного сервиса, инициализируя все реализации.
func NewService(deps ServiceDependencies) *Service {
	authService := NewAuthService(
		deps.Repos.User,
		deps.Repos.Token,
		deps.Repos.Cache,
		deps.Repos.Transactor,
		deps.SigningKey,
		deps.TokenTTL,
	)

	return &Service{
		Authorization: authService,
		User:          NewUserService(deps.Repos.User, deps.Repos.Appointment, deps.Storage),
		Doctor:        NewDoctorService(deps.Repos.Doctor),
		Appointment:   NewAppointmentService(deps.Repos.Appointment, deps.Repos.Doctor, deps.Location),
		Directory:     NewDirectoryService(deps.Repos.Directory),
		Info:          NewInfoService(deps.Repos.Service, deps.Repos.Info),
		Prescription:  NewPrescriptionService(deps.Repos.Prescription),
		MedicalCard:   NewMedicalCardService(deps.Repos.MedicalCard, deps.Repos.Prescription, deps.Storage),
		Admin:         NewAdminService(deps.Repos, deps.SigningKey, deps.TokenTTL),
	}
}
