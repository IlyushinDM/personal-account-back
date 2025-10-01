package models

import "time"

// AdminRole определяет тип для ролей администратора.
type AdminRole string

const (
	RoleSuperAdmin AdminRole = "superadmin"
	RoleAdmin      AdminRole = "admin"
)

// Admin представляет пользователя-администратора в системе.
type Admin struct {
	ID           uint64    `gorm:"primarykey" json:"id"`
	Login        string    `gorm:"unique;not null" json:"login"`
	PasswordHash string    `json:"-"`
	FullName     string    `gorm:"not null" json:"fullName"`
	Role         AdminRole `gorm:"type:varchar(50);not null" json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName возвращает имя таблицы в базе данных.
func (Admin) TableName() string {
	return "medical_center.admins"
}

// AdminDashboardStats представляет DTO для статистики на дашборде.
type AdminDashboardStats struct {
	ActiveUsers    int64   `json:"activeUsers"`
	NewUsersToday  int64   `json:"newUsersToday"`
	Appointments   int64   `json:"appointments"`
	CompletedTotal int64   `json:"completedTotal"`
	TotalRevenue   float64 `json:"totalRevenue"`
}
