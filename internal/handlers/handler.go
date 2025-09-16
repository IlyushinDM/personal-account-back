// Package handlers определяет HTTP-слой приложения (API).
// Он отвечает за обработку входящих запросов, вызов соответствующих сервисов
// и форматирование ответов для отправки клиенту.
package handlers

import (
	"lk/internal/services"

	"github.com/gin-gonic/gin"

	//_ "lk/docs" // ссылка на сгенерированную документацию

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
		}

		// Защищенные эндпоинты, требующие валидного JWT
		authorized := apiV1.Group("/")
		authorized.Use(h.userIdentity) // Применяем middleware ко всей группе
		{
			// Профиль пользователя (FR-2.x)
			profile := authorized.Group("/profile")
			{
				profile.GET("/", h.getProfile)
				profile.PATCH("/", h.updateProfile)
			}

			// Специалисты и справочники (FR-3.x)
			departments := authorized.Group("/departments")
			{
				departments.GET("/", h.getDepartmentsTree) // Получить дерево отделений и специальностей
			}

			specialists := authorized.Group("/specialists")
			{
				specialists.GET("/", h.findSpecialists)      // Поиск врачей по параметрам
				specialists.GET("/:id", h.getSpecialistByID) // Получить детальный профиль врача
			}

			// Записи на прием (FR-3.x)
			appointments := authorized.Group("/appointments")
			{
				appointments.GET("/", h.getUserAppointments)     // Получить все записи пользователя
				appointments.POST("/", h.createAppointment)      // Создать новую запись
				appointments.DELETE("/:id", h.cancelAppointment) // Отменить запись
			}
		}
	}
	return router
}
