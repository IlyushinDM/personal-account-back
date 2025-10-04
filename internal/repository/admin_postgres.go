package repository

import (
	"context"
	"time"

	"lk/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminPostgres struct {
	db *gorm.DB
}

func NewAdminPostgres(db *gorm.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

// GetByLogin находит администратора по логину.
func (r *AdminPostgres) GetByLogin(ctx context.Context, login string) (models.Admin, error) {
	var admin models.Admin
	err := r.db.WithContext(ctx).Where("login = ?", login).First(&admin).Error
	return admin, err
}

// GetByID находит администратора по ID.
func (r *AdminPostgres) GetByID(ctx context.Context, id uint64) (models.Admin, error) {
	var admin models.Admin
	err := r.db.WithContext(ctx).First(&admin, id).Error
	return admin, err
}

// --- User ---

func (r *AdminPostgres) GetAllUsers(ctx context.Context, params models.PaginationParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := query.Order("id DESC").Limit(params.Limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *AdminPostgres) UpdateUser(ctx context.Context, user models.User, profile models.UserProfile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
			return err
		}
		if err := tx.Model(&profile).Where("user_id = ?", user.ID).Updates(
			profile).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *AdminPostgres) DeleteUser(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Select(clause.Associations).Delete(&models.User{}, userID).Error
}

func (r *AdminPostgres) GetUserAppointments(ctx context.Context, userID uint64, params models.PaginationParams) (
	[]models.Appointment, int64, error,
) {
	var appointments []models.Appointment
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Appointment{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (params.Page - 1) * params.Limit
	err := query.Order("appointment_date DESC").Limit(
		params.Limit).Offset(offset).Find(&appointments).Error
	return appointments, total, err
}

func (r *AdminPostgres) GetUserAnalyses(ctx context.Context, userID uint64, params models.PaginationParams) (
	[]models.LabAnalysis, int64, error,
) {
	var analyses []models.LabAnalysis
	var total int64
	query := r.db.WithContext(ctx).Model(&models.LabAnalysis{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (params.Page - 1) * params.Limit
	err := query.Order("assigned_date DESC").Limit(params.Limit).Offset(offset).Find(&analyses).Error
	return analyses, total, err
}

// --- Doctor ---

func (r *AdminPostgres) GetAllSpecialists(ctx context.Context, params models.PaginationParams) (
	[]models.Doctor, int64, error,
) {
	var doctors []models.Doctor
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Doctor{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (params.Page - 1) * params.Limit
	err := query.Order("last_name ASC").Limit(params.Limit).Offset(offset).Find(&doctors).Error
	return doctors, total, err
}

func (r *AdminPostgres) CreateDoctor(ctx context.Context, doctor models.Doctor) (uint64, error) {
	result := r.db.WithContext(ctx).Create(&doctor)
	return doctor.ID, result.Error
}

func (r *AdminPostgres) UpdateDoctor(ctx context.Context, doctor models.Doctor) error {
	return r.db.WithContext(ctx).Save(&doctor).Error
}

func (r *AdminPostgres) DeleteDoctor(ctx context.Context, doctorID uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Doctor{}, doctorID).Error
}

func (r *AdminPostgres) GetDoctorSchedule(ctx context.Context, doctorID uint64) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.WithContext(ctx).Where(
		"doctor_id = ?", doctorID).Order("date ASC").Find(&schedules).Error
	return schedules, err
}

func (r *AdminPostgres) UpdateDoctorSchedule(ctx context.Context, doctorID uint64, schedules []models.Schedule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Удаляем старое расписание
		if err := tx.Where("doctor_id = ?", doctorID).Delete(
			&models.Schedule{}).Error; err != nil {
			return err
		}
		// Вставляем новое
		if len(schedules) > 0 {
			if err := tx.Create(&schedules).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// --- Appointment ---

func (r *AdminPostgres) GetAllAppointments(
	ctx context.Context, params models.PaginationParams, filters map[string]interface{},
) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Appointment{})
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (params.Page - 1) * params.Limit
	err := query.Order(
		"appointment_date DESC").Limit(params.Limit).Offset(offset).Find(&appointments).Error
	return appointments, total, err
}

func (r *AdminPostgres) DeleteAppointment(ctx context.Context, appointmentID uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Appointment{}, appointmentID).Error
}

func (r *AdminPostgres) GetAppointmentStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	var total, cancelled, completed int64

	r.db.WithContext(ctx).Model(&models.Appointment{}).Count(&total)
	r.db.WithContext(ctx).Model(&models.Appointment{}).Where(
		"status_id IN (?, ?)", models.StatusCancelledByClinic, models.StatusCancelledByPatient).Count(
		&cancelled)
	r.db.WithContext(ctx).Model(&models.Appointment{}).Where(
		"status_id = ?", models.StatusCompleted).Count(&completed)

	stats["total"] = total
	stats["cancelled"] = cancelled
	stats["completed"] = completed
	return stats, nil
}

// --- Service & Department ---

func (r *AdminPostgres) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service
	err := r.db.WithContext(ctx).Order("name ASC").Find(&services).Error
	return services, err
}

func (r *AdminPostgres) CreateService(ctx context.Context, service models.Service) (uint64, error) {
	result := r.db.WithContext(ctx).Create(&service)
	return service.ID, result.Error
}

func (r *AdminPostgres) UpdateService(ctx context.Context, service models.Service) error {
	return r.db.WithContext(ctx).Save(&service).Error
}

func (r *AdminPostgres) DeleteService(ctx context.Context, serviceID uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Service{}, serviceID).Error
}

func (r *AdminPostgres) CreateDepartment(ctx context.Context, department models.Department) (uint32, error) {
	result := r.db.WithContext(ctx).Create(&department)
	return department.ID, result.Error
}

func (r *AdminPostgres) UpdateDepartment(ctx context.Context, department models.Department) error {
	return r.db.WithContext(ctx).Save(&department).Error
}

func (r *AdminPostgres) DeleteDepartment(ctx context.Context, departmentID uint32) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, departmentID).Error
}

// --- Статистика ---

func (r *AdminPostgres) GetDashboardStats(ctx context.Context) (models.AdminDashboardStats, error) {
	var stats models.AdminDashboardStats

	r.db.WithContext(ctx).Model(&models.User{}).Where(
		"is_active = ?", true).Count(&stats.ActiveUsers)

	todayStart := time.Now().Truncate(24 * time.Hour)
	r.db.WithContext(ctx).Model(&models.User{}).Where(
		"created_at >= ?", todayStart).Count(&stats.NewUsersToday)

	r.db.WithContext(ctx).Model(&models.Appointment{}).Count(&stats.Appointments)

	type revenueStats struct {
		TotalRevenue   float64
		CompletedTotal int64
	}
	var revenue revenueStats
	r.db.WithContext(ctx).Model(&models.Appointment{}).
		Select("COUNT(*) as completed_total, SUM(price_at_booking) as total_revenue").
		Where("status_id = ?", models.StatusCompleted).
		Scan(&revenue)

	stats.TotalRevenue = revenue.TotalRevenue
	stats.CompletedTotal = revenue.CompletedTotal

	return stats, nil
}
