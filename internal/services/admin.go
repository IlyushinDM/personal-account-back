package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// adminService реализует интерфейс AdminService.
type adminService struct {
	repos      *repository.Repository
	signingKey string
	tokenTTL   time.Duration
}

// NewAdminService создает новый сервис для администрирования.
func NewAdminService(repos *repository.Repository, signingKey string, tokenTTL time.Duration) AdminService {
	return &adminService{
		repos:      repos,
		signingKey: signingKey,
		tokenTTL:   tokenTTL,
	}
}

// --- Auth & Dashboard ---

// Login аутентифицирует администратора и возвращает access токен.
func (s *adminService) Login(ctx context.Context, login, password string) (map[string]string, error) {
	admin, err := s.repos.Admin.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewUnauthorizedError("invalid login or password", nil)
		}
		return nil, NewInternalServerError("database error while getting admin", err)
	}

	if err := utils.CheckPasswordHash(password, admin.PasswordHash); err != nil {
		return nil, NewUnauthorizedError("invalid login or password", nil)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(s.tokenTTL).Unix(),
		"iat":    time.Now().Unix(),
		"sub":    admin.ID,
		"role":   admin.Role,
		"is_adm": true,
	})

	accessToken, err := token.SignedString([]byte(s.signingKey))
	if err != nil {
		return nil, NewInternalServerError("failed to generate admin token", err)
	}

	return map[string]string{"accessToken": accessToken}, nil
}

// ParseAdminToken проверяет токен админа и возвращает его ID.
func (s *adminService) ParseAdminToken(accessToken string) (uint64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.signingKey), nil
		})
	if err != nil {
		return 0, NewUnauthorizedError("invalid admin token", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, NewUnauthorizedError("invalid admin token claims", nil)
	}

	isAdmin, ok := (*claims)["is_adm"].(bool)
	if !ok || !isAdmin {
		return 0, NewUnauthorizedError("token is not an admin token", nil)
	}

	subFloat, ok := (*claims)["sub"].(float64)
	if !ok {
		return 0, NewUnauthorizedError("invalid subject claim in admin token", nil)
	}

	return uint64(subFloat), nil
}

// GetDashboardStats получает статистику для дашборда.
func (s *adminService) GetDashboardStats(ctx context.Context) (models.AdminDashboardStats, error) {
	stats, err := s.repos.Admin.GetDashboardStats(ctx)
	if err != nil {
		return models.AdminDashboardStats{}, NewInternalServerError("failed to get dashboard stats", err)
	}
	return stats, nil
}

// --- User (Пациент) ---

func (s *adminService) GetAllUsers(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error) {
	users, total, err := s.repos.Admin.GetAllUsers(ctx, params)
	if err != nil {
		return nil, 0, NewInternalServerError("failed to get all users", err)
	}
	return users, total, nil
}

func (s *adminService) GetUserByID(ctx context.Context, userID uint64) (*models.User, *models.UserProfile, error) {
	user, err := s.repos.User.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, NewNotFoundError("user not found", err)
		}
		return nil, nil, NewInternalServerError("failed to get user by id", err)
	}

	profile, err := s.repos.User.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		// Если пользователь есть, а профиля нет - это внутренняя ошибка.
		return nil, nil, NewInternalServerError("failed to get user profile by id", err)
	}
	return &user, &profile, nil
}

func (s *adminService) UpdateUser(ctx context.Context, userID uint64, input UpdateUserInput) error {
	user, profile, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err // Ошибка уже обернута в GetUserByID
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	if input.FirstName != nil {
		profile.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		profile.LastName = *input.LastName
	}
	if input.Email != nil {
		profile.Email.String = *input.Email
		profile.Email.Valid = true
	} else {
		profile.Email.Valid = false
	}

	err = s.repos.Admin.UpdateUser(ctx, *user, *profile)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return NewConflictError("phone or email already in use", err)
		}
		return NewInternalServerError("failed to update user in db", err)
	}
	return nil
}

func (s *adminService) DeleteUser(ctx context.Context, userID uint64) error {
	if err := s.repos.Admin.DeleteUser(ctx, userID); err != nil {
		return NewInternalServerError("failed to delete user", err)
	}
	return nil
}

func (s *adminService) GetUserAppointments(ctx context.Context, userID uint64, params models.PaginationParams) (
	[]models.Appointment, int64, error,
) {
	return s.repos.Admin.GetUserAppointments(ctx, userID, params)
}

func (s *adminService) GetUserAnalyses(ctx context.Context, userID uint64, params models.PaginationParams) (
	[]models.LabAnalysis, int64, error,
) {
	return s.repos.Admin.GetUserAnalyses(ctx, userID, params)
}

// --- Doctor (Специалист aka врач) ---

func (s *adminService) GetAllSpecialists(ctx context.Context, params models.PaginationParams) (
	[]models.Doctor, int64, error,
) {
	return s.repos.Admin.GetAllSpecialists(ctx, params)
}

func (s *adminService) CreateSpecialist(ctx context.Context, input CreateDoctorInput) (uint64, error) {
	doctor := models.Doctor{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		SpecialtyID:     input.SpecialtyID,
		ExperienceYears: input.ExperienceYears,
	}
	if input.Patronymic != nil {
		doctor.Patronymic.String, doctor.Patronymic.Valid = *input.Patronymic, true
	}
	if input.Recommendations != nil {
		doctor.Recommendations.String, doctor.Recommendations.Valid = *input.Recommendations, true
	}

	id, err := s.repos.Admin.CreateDoctor(ctx, doctor)
	if err != nil {
		return 0, NewInternalServerError("failed to create specialist", err)
	}
	return id, nil
}

func (s *adminService) GetSpecialistByID(ctx context.Context, doctorID uint64) (models.Doctor, error) {
	return s.repos.Doctor.GetDoctorByID(ctx, doctorID)
}

func (s *adminService) UpdateSpecialist(ctx context.Context, doctorID uint64, input UpdateDoctorInput) error {
	doctor, err := s.repos.Doctor.GetDoctorByID(ctx, doctorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewNotFoundError("specialist to update not found", err)
		}
		return NewInternalServerError("failed to get specialist", err)
	}

	if input.FirstName != nil {
		doctor.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		doctor.LastName = *input.LastName
	}
	if input.Patronymic != nil {
		doctor.Patronymic.String, doctor.Patronymic.Valid = *input.Patronymic, true
	}
	if input.SpecialtyID != nil {
		doctor.SpecialtyID = *input.SpecialtyID
	}
	if input.ExperienceYears != nil {
		doctor.ExperienceYears = *input.ExperienceYears
	}
	if input.Recommendations != nil {
		doctor.Recommendations.String, doctor.Recommendations.Valid = *input.Recommendations, true
	}

	return s.repos.Admin.UpdateDoctor(ctx, doctor)
}

func (s *adminService) DeleteSpecialist(ctx context.Context, doctorID uint64) error {
	return s.repos.Admin.DeleteDoctor(ctx, doctorID)
}

func (s *adminService) GetSpecialistSchedule(ctx context.Context, doctorID uint64) ([]models.Schedule, error) {
	return s.repos.Admin.GetDoctorSchedule(ctx, doctorID)
}

func (s *adminService) UpdateSpecialistSchedule(ctx context.Context, doctorID uint64, input UpdateScheduleInput) error {
	schedules := make([]models.Schedule, len(input.Schedules))
	for i, item := range input.Schedules {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			return NewBadRequestError("invalid date format: "+item.Date, err)
		}
		startTime, err := time.Parse("15:04", item.StartTime)
		if err != nil {
			return NewBadRequestError("invalid start time format: "+item.StartTime, err)
		}
		endTime, err := time.Parse("15:04", item.EndTime)
		if err != nil {
			return NewBadRequestError("invalid end time format: "+item.EndTime, err)
		}

		schedules[i] = models.Schedule{
			DoctorID:  doctorID,
			Date:      date,
			StartTime: startTime,
			EndTime:   endTime,
		}
	}
	return s.repos.Admin.UpdateDoctorSchedule(ctx, doctorID, schedules)
}

// --- Appointment ---

func (s *adminService) GetAllAppointments(ctx context.Context, params models.PaginationParams, filters map[string]any) (
	[]models.Appointment, int64, error,
) {
	return s.repos.Admin.GetAllAppointments(ctx, params, filters)
}

func (s *adminService) GetAppointmentStats(ctx context.Context) (map[string]int64, error) {
	return s.repos.Admin.GetAppointmentStats(ctx)
}

func (s *adminService) GetAppointmentDetails(ctx context.Context, appointmentID uint64) (models.Appointment, error) {
	appointment, err := s.repos.Appointment.GetAppointmentByID(ctx, appointmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Appointment{}, NewNotFoundError("appointment not found", err)
		}
		return models.Appointment{}, NewInternalServerError("failed to get appointment details", err)
	}
	return appointment, nil
}

func (s *adminService) UpdateAppointmentStatus(ctx context.Context, appointmentID uint64, statusID uint32) error {
	return s.repos.Appointment.UpdateAppointmentStatus(ctx, appointmentID, statusID)
}

func (s *adminService) DeleteAppointment(ctx context.Context, appointmentID uint64) error {
	return s.repos.Admin.DeleteAppointment(ctx, appointmentID)
}

// --- Service & Department ---

func (s *adminService) GetAllServices(ctx context.Context) ([]models.Service, error) {
	return s.repos.Admin.GetAllServices(ctx)
}

func (s *adminService) CreateService(ctx context.Context, input CreateServiceInput) (uint64, error) {
	service := models.Service{
		Name:            input.Name,
		Price:           input.Price,
		DurationMinutes: input.DurationMinutes,
		DoctorID:        input.DoctorID,
	}
	if input.Description != nil {
		service.Description.String, service.Description.Valid = *input.Description, true
	}
	if input.Recommendations != nil {
		service.Recommendations.String, service.Recommendations.Valid = *input.Recommendations, true
	}

	id, err := s.repos.Admin.CreateService(ctx, service)
	if err != nil {
		return 0, NewInternalServerError("failed to create service", err)
	}
	return id, nil
}

func (s *adminService) UpdateService(ctx context.Context, serviceID uint64, input UpdateServiceInput) error {
	// TODO: Реализовать логику получения сервиса по ID и его обновления
	return NewInternalServerError("Not implemented yet", nil)
}

func (s *adminService) DeleteService(ctx context.Context, serviceID uint64) error {
	return s.repos.Admin.DeleteService(ctx, serviceID)
}

func (s *adminService) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	return s.repos.Directory.GetAllDepartments(ctx)
}

func (s *adminService) CreateDepartment(ctx context.Context, input CreateDepartmentInput) (uint32, error) {
	department := models.Department{Name: input.Name}
	id, err := s.repos.Admin.CreateDepartment(ctx, department)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return 0, NewConflictError("department with this name already exists", err)
		}
		return 0, NewInternalServerError("failed to create department", err)
	}
	return id, nil
}

func (s *adminService) UpdateDepartment(ctx context.Context, departmentID uint32, input UpdateDepartmentInput) error {
	if input.Name == nil {
		return NewBadRequestError("name is required for update", nil)
	}
	department := models.Department{ID: departmentID, Name: *input.Name}
	err := s.repos.Admin.UpdateDepartment(ctx, department)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return NewConflictError("department with this name already exists", err)
		}
		return NewInternalServerError("failed to update department", err)
	}
	return nil
}

func (s *adminService) DeleteDepartment(ctx context.Context, departmentID uint32) error {
	return s.repos.Admin.DeleteDepartment(ctx, departmentID)
}

// TODO: Реализовать остальные методы AdminService
// (управление анализами, назначениями, семьей, настройками, бекапами и т.д.)
