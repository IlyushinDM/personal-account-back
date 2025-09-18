// Package main является точкой входа для серверного приложения "Личный Кабинет" (lk).
// Он инициализирует конфигурацию, базу данных, все слои приложения,
// а затем запускает HTTP-сервер с корректным завершением работы.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lk/internal/config"
	"lk/internal/handlers"
	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/server"
	"lk/internal/services"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title API Личного Кабинета
// @version 1.0
// @description Серверная часть (бэкенд) для личного кабинета пациента. Предоставляет RESTful API.
// @description
// @description ### Аутентификация
// @description Для доступа к защищенным эндпоинтам необходимо передавать JWT в заголовке `Authorization`.
// @description Формат: `Authorization: Bearer <ваш_access_token>`

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Используйте "Bearer " перед вашим токеном."

// swagger:meta
// nolint
type _ struct {
	// Этот блок нужен для того, чтобы swag корректно распознал
	// типы из пакета models, используемые в аннотациях.
	_ models.DepartmentWithSpecialties
	_ models.Doctor
	_ models.Appointment
	_ models.UserProfile
	_ models.Specialty
	_ models.PaginatedDoctorsResponse
	_ models.Recommendation
	_ models.AvailableDatesResponse
	_ models.AvailableSlotsResponse
}

func main() {
	// Загружаем переменные окружения из .env файла, если он существует.
	// Это удобно для локальной разработки.
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются системные переменные окружения")
	}

	// Инициализируем конфигурацию приложения.
	// Приложение завершит работу, если обязательные переменные не установлены.
	cfg := config.MustLoad()

	// Устанавливаем соединение с базой данных PostgreSQL.
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()
	log.Println("соединение с базой данных установлено")

	// Инициализируем слои приложения (репозиторий, сервисы, обработчики).
	// Внедряем зависимости сверху вниз (DI: Dependency Injection).
	repos := repository.NewRepository(db)
	serviceDeps := services.ServiceDependencies{
		Repos:      repos,
		SigningKey: cfg.Auth.JWTSecretKey,
		TokenTTL:   cfg.Auth.TokenTTL,
	}
	services := services.NewService(serviceDeps)
	handlers := handlers.NewHandler(services)
	log.Println("слои приложения инициализированы")

	// Запускаем HTTP-сервер в отдельной горутине.
	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPServer.Port, handlers.InitRoutes()); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка при запуске http-сервера: %v", err)
		}
	}()
	log.Printf("сервер запущен на порту: %s", cfg.HTTPServer.Port)

	// Настраиваем graceful shutdown (плавное завершение работы).
	// Ждем сигнала от операционной системы (SIGINT, SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// При получении сигнала начинаем процесс остановки сервера.
	log.Println("сервер выключается")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("произошла ошибка при выключении сервера: %v", err)
	}
}
