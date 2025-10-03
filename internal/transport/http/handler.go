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
	services *services.Service
	userRepo repository.UserRepository // Зависимость для middleware
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(services *services.Service, userRepo repository.UserRepository) *Handler {
	return &Handler{
		services: services,
		userRepo: userRepo,
	}
}

// InitRoutes настраивает переданный роутер Gin, добавляя в него все эндпоинты приложения.
func (h *Handler) InitRoutes(router *gin.Engine) {
	// Эндпоинт для Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Группа для API версии 1
	apiV1 := router.Group("/api/v1")
	{
		// Публичные эндпоинты для аутентификации
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

		// Защищенные эндпоинты, требующие валидного JWT
		authorized := apiV1.Group("/")
		authorized.Use(h.userIdentity)
		{
			// --- ПРОФИЛЬ ПОЛЬЗОВАТЕЛЯ (FR-2.x) ---
			profile := authorized.Group("/profile")
			{
				profile.GET("/", h.getProfile)
				profile.PATCH("/", h.updateProfile)
				profile.POST("/avatar", h.updateAvatar)
			}

			// --- СПРАВОЧНИКИ И ОБЩАЯ ИНФОРМАЦИЯ (FR-2.x, FR-3.x) ---
			authorized.GET("/clinic-info", h.getClinicInfo)
			authorized.GET("/legal/documents", h.getLegalDocuments)
			authorized.GET("/specialties", h.getSpecialties)
			departments := authorized.Group("/departments")
			{
				departments.GET("/", h.getDepartmentsTree)
			}

			// --- СПЕЦИАЛИСТЫ (FR-3.x) ---
			specialists := authorized.Group("/specialists")
			{
				specialists.GET("/", h.findSpecialists)
				specialists.GET("/:id", h.getSpecialistByID)
				specialists.GET("/:id/recommendations", h.getSpecialistRecommendations)
			}

			// --- ОТЗЫВЫ (FR-3.x) ---
			reviews := authorized.Group("/reviews")
			{
				reviews.GET("/", h.getReviews)
			}

			// --- УСЛУГИ (FR-5.x) ---
			services := authorized.Group("/services")
			{
				services.GET("/:id/recommendations", h.getServiceRecommendations)
			}

			// --- ЗАПИСИ НА ПРИЕМ (FR-2.x, FR-5.x) ---
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

			// --- НАЗНАЧЕНИЯ (FR-2.x) ---
			prescriptions := authorized.Group("/prescriptions")
			{
				prescriptions.GET("/active", h.getActivePrescriptions)
				prescriptions.POST("/:id/archive", h.archivePrescription)
			}

			// --- МЕДКАРТА (FR-4.x) ---
			medicalCard := authorized.Group("/medical-card")
			{
				medicalCard.GET("/visits", h.getVisits)
				medicalCard.GET("/analyses", h.getAnalyses)
				medicalCard.GET("/archive/prescriptions", h.getArchivedPrescriptions)
				medicalCard.GET("/summary", h.getSummary)
				medicalCard.POST("/archive/prescriptions", h.archivePrescriptionFromCard)
			}

			// Отдельный роут для скачивания файлов (FR-4.4)
			authorized.GET("/files/:id", h.downloadFile)
		}
	}
}
