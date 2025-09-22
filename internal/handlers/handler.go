// Package handlers определяет HTTP-слой приложения (API).
// Он отвечает за обработку входящих запросов, вызов соответствующих сервисов
// и форматирование ответов для отправки клиенту.
package handlers

import (
	"lk/internal/services"

	"github.com/gin-gonic/gin"

	_ "lk/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler - это контейнер для всех зависимостей слоя обработчиков.
type Handler struct {
	services *services.Service
}

// NewHandler создает новый экземпляр обработчика.
func NewHandler(services *services.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// InitRoutes настраивает и возвращает роутер Gin со всеми эндпоинтами приложения.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Подключаем middleware для перехвата паник и возврата 500 ошибки.
	router.Use(gin.Recovery())

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
			// TODO: Реализовать эндпоинты FR-1.3, FR-1.4, FR-1.5, FR-1.6, FR-1.7
			// auth.GET("/gosuslugi", h.gosuslugiLogin)
			// auth.GET("/gosuslugi/callback", h.gosuslugiCallback)
			// auth.POST("/forgot-password", h.forgotPassword)
			// auth.POST("/reset-password", h.resetPassword)
			// auth.POST("/refresh", h.refreshToken)
			// auth.POST("/logout", h.logout)
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
	return router
}
