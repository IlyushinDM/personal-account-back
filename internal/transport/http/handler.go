// Package http определяет HTTP-слой приложения (API).
// Он отвечает за обработку входящих запросов, вызов соответствующих сервисов
// и форматирование ответов для отправки клиенту.
package http

import (
	"lk/internal/repository"
	"lk/internal/services"

	"github.com/gin-gonic/gin"

	_ "lk/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler - это контейнер для всех зависимостей слоя обработчиков.
type Handler struct {
	services  *services.Service
	userRepo  repository.UserRepository  // Зависимость для middleware пользователя
	adminRepo repository.AdminRepository // Зависимость для middleware админа
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(services *services.Service, userRepo repository.UserRepository, adminRepo repository.AdminRepository) *Handler {
	return &Handler{
		services:  services,
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
}

// InitRoutes настраивает переданный роутер Gin, добавляя в него все эндпоинты приложения.
func (h *Handler) InitRoutes(router *gin.Engine) {
	// Эндпоинт для Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа для API версии 1
	apiV1 := router.Group("/api/v1")
	{
		// --- ПУБЛИЧНАЯ ЧАСТЬ: АУТЕНТИФИКАЦИЯ ПАЦИЕНТОВ ---
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", h.signUp)
			auth.POST("/login", h.signIn)
			auth.POST("/refresh", h.refresh)
			auth.POST("/logout", h.logout)
			auth.POST("/forgot-password", h.forgotPassword)
			auth.POST("/reset-password", h.resetPassword)
			auth.GET("/gosuslugi", h.gosuslugiLogin)
			auth.POST("/gosuslugi/callback", h.gosuslugiCallback)
		}

		// --- ЗАЩИЩЕННАЯ ЧАСТЬ: ЛИЧНЫЙ КАБИНЕТ ПАЦИЕНТА ---
		authorized := apiV1.Group("/")
		authorized.Use(h.userIdentity)
		{
			// Профиль пользователя (FR-2.x)
			profile := authorized.Group("/profile")
			{
				profile.GET("/", h.getProfile)
				profile.PATCH("/", h.updateProfile)
				profile.POST("/avatar", h.updateAvatar)
			}

			// Справочники и общая информация
			authorized.GET("/clinic-info", h.getClinicInfo)
			authorized.GET("/legal/documents", h.getLegalDocuments)
			authorized.GET("/specialties", h.getSpecialties)
			authorized.GET("/departments", h.getDepartmentsTree)

			// Специалисты и услуги
			authorized.GET("/specialists", h.findSpecialists)
			authorized.GET("/specialists/:id", h.getSpecialistByID)
			authorized.GET("/specialists/:id/recommendations", h.getSpecialistRecommendations)
			authorized.GET("/services/:id/recommendations", h.getServiceRecommendations)

			// Записи на прием (FR-2.x, FR-5.x)
			appointments := authorized.Group("/appointments")
			{
				appointments.GET("/", h.getUserAppointments)
				appointments.POST("/", h.createAppointment)
				appointments.GET("/upcoming", h.getUpcomingAppointments)
				appointments.DELETE("/:id", h.cancelAppointment)
				appointments.GET("/available-dates", h.getAvailableDates)
				appointments.GET("/available-slots", h.getAvailableSlots)
				appointments.GET("/slots-by-range", h.getAvailableSlotsByRange)
			}

			// Назначения (FR-2.x)
			prescriptions := authorized.Group("/prescriptions")
			{
				prescriptions.GET("/active", h.getActivePrescriptions)
				prescriptions.POST("/:id/archive", h.archivePrescription)
			}

			// Медкарта (FR-4.x)
			medicalCard := authorized.Group("/medical-card")
			{
				medicalCard.GET("/visits", h.getVisits)
				medicalCard.GET("/analyses", h.getAnalyses)
				medicalCard.GET("/archive/prescriptions", h.getArchivedPrescriptions)
				medicalCard.GET("/summary", h.getSummary)
				medicalCard.POST("/archive/prescriptions", h.archivePrescriptionFromCard)
			}

			// Скачивание файлов (FR-4.4)
			authorized.GET("/files/:id", h.downloadFile)
		}

		// --- АДМИН-ПАНЕЛЬ ---
		admin := apiV1.Group("/admin")
		{
			// Публичный эндпоинт для входа администратора
			admin.POST("/login", h.adminLogin)

			// Группа, защищенная middleware администратора
			adminAuthorized := admin.Group("/")
			adminAuthorized.Use(h.adminIdentity)
			{
				// 1. Управление пользователями (пациентами)
				users := adminAuthorized.Group("/users")
				{
					users.GET("/", h.adminGetAllUsers)
					users.GET("/:id", h.adminGetUserByID)
					users.PATCH("/:id", h.adminUpdateUser)
					users.DELETE("/:id", h.adminDeleteUser)
					users.GET("/:id/appointments", h.adminGetUserAppointments)
					users.GET("/:id/analyses", h.adminGetUserAnalyses)
				}

				// 2. Управление врачами (специалистами)
				specialists := adminAuthorized.Group("/specialists")
				{
					specialists.GET("/", h.adminGetAllSpecialists)
					specialists.POST("/", h.adminCreateSpecialist)
					specialists.GET("/:id", h.adminGetSpecialistByID)
					specialists.PUT("/:id", h.adminUpdateSpecialist)
					specialists.DELETE("/:id", h.adminDeleteSpecialist)
					specialists.GET("/:id/schedule", h.adminGetSpecialistSchedule)
					specialists.POST("/:id/schedule", h.adminUpdateSpecialistSchedule)
				}

				// 3. Управление записями на приём
				appointments := adminAuthorized.Group("/appointments")
				{
					appointments.GET("/", h.adminGetAllAppointments)
					appointments.GET("/statistics", h.adminGetAppointmentStats)
					appointments.GET("/:id", h.adminGetAppointmentDetails)
					appointments.PATCH("/:id", h.adminUpdateAppointmentStatus)
					appointments.DELETE("/:id", h.adminDeleteAppointment)
				}

				// 4. Управление услугами и отделениями
				services := adminAuthorized.Group("/services")
				{
					services.GET("/", h.adminGetAllServices)
					services.POST("/", h.adminCreateService)
					services.PUT("/:id", h.adminUpdateService)
					services.DELETE("/:id", h.adminDeleteService)
				}
				departments := adminAuthorized.Group("/departments")
				{
					departments.GET("/", h.adminGetAllDepartments)
					departments.POST("/", h.adminCreateDepartment)
					departments.PUT("/:id", h.adminUpdateDepartment)
					departments.DELETE("/:id", h.adminDeleteDepartment)
				}

				// 5. Управление анализами и назначениями
				analyses := adminAuthorized.Group("/analyses")
				{
					analyses.GET("/", h.adminGetAllAnalyses)
					analyses.POST("/", h.adminCreateAnalysisResult)
					analyses.PATCH("/:id", h.adminUpdateAnalysis)
					analyses.DELETE("/:id", h.adminDeleteAnalysis)
				}
				prescriptions := adminAuthorized.Group("/prescriptions")
				{
					prescriptions.GET("/", h.adminGetAllPrescriptions)
					prescriptions.POST("/", h.adminCreatePrescription)
				}

				// 7. Управление семьей
				family := adminAuthorized.Group("/family-relations")
				{
					family.GET("/", h.adminGetFamilyRelations)
					family.DELETE("/:id", h.adminDeleteFamilyRelation)
				}

				// 8. Системные настройки и статистика
				adminAuthorized.GET("/dashboard", h.getAdminDashboard)
				adminAuthorized.GET("/audit-logs", h.adminGetAuditLogs)
				settings := adminAuthorized.Group("/clinic-settings")
				{
					settings.GET("/", h.adminGetClinicSettings)
					settings.PUT("/", h.adminUpdateClinicSettings) // PUT для полного обновления
				}

				// 9. Управление документами
				legal := adminAuthorized.Group("/legal-documents")
				{
					legal.GET("/", h.adminGetLegalDocs)
					legal.POST("/", h.adminCreateLegalDoc)
					legal.PUT("/:id", h.adminUpdateLegalDoc)
				}

				// 10. Резервное копирование
				backup := adminAuthorized.Group("/backup")
				{
					backup.POST("/", h.adminCreateBackup)
					backup.GET("/list", h.adminGetBackupList)
					backup.POST("/restore", h.adminRestoreFromBackup)
				}
			}
		}
	}
}
