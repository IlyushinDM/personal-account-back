// Package main является точкой входа для серверного приложения "Личный Кабинет" (lk).
// Он инициализирует конфигурацию, базу данных, все слои приложения,
// а затем запускает HTTP-сервер с корректным завершением работы.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lk/internal/config"
	"lk/internal/handlers"
	"lk/internal/logger"
	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/server"
	"lk/internal/services"
	"lk/internal/storage"
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
	if err := godotenv.Load(); err != nil {
		logger.Default().Warn("Файл .env не найден, используются системные переменные окружения")
	}

	logger.Init(os.Getenv("LOG_DIR"))
	logger.Default().Info("логгер инициализирован")

	// Инициализируем конфигурацию приложения.
	cfg := config.MustLoad()

	// Загружаем и валидируем часовой пояс
	location, err := time.LoadLocation(cfg.ClinicTimezone)
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось загрузить указанный часовой пояс")
	}
	logger.Default().Info(fmt.Sprintf("приложение работает в часовом поясе: %s", location.String()))

	// Устанавливаем соединение с базой данных PostgreSQL.
	DB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.NewGORMLogger(),
	})
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к базе данных")
	}

	db, err := DB.DB()
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось получить *sql.DB из gorm")
	}
	defer db.Close()
	logger.Default().Info("соединение с базой данных установлено")

	// Инициализируем клиент Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Проверяем соединение с Redis
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к Redis")
	}
	defer redisClient.Close()
	logger.Default().Info("соединение с Redis установлено")

	// Инициализируем клиент MinIO
	storageClient, err := storage.NewMinIOClient(cfg.Minio)
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к MinIO")
	}
	logger.Default().Info("соединение с MinIO установлено")

	// Инициализируем слои приложения (репозиторий, сервисы, обработчики).
	repos := repository.NewRepository(DB, redisClient)
	serviceDeps := services.ServiceDependencies{
		Repos:      repos,
		Storage:    storageClient,
		Location:   location,
		SigningKey: cfg.Auth.JWTSecretKey,
		TokenTTL:   cfg.Auth.TokenTTL,
	}
	services := services.NewService(serviceDeps)
	handlers := handlers.NewHandler(services)
	logger.Default().Info("слои приложения инициализированы")

	// Инициализируем роутер и применяем глобальные middleware
	router := gin.New()
	router.Use(logger.GinLogger(), gin.Recovery())

	// Наполняем роутер эндпоинтами
	handlers.InitRoutes(router)

	// Запускаем HTTP-сервер в отдельной горутине.
	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPServer.Port, router); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Default().WithError(err).Fatal("ошибка при запуске http-сервера")
		}
	}()
	log.Printf("сервер запущен на порту: %s", cfg.HTTPServer.Port)

	// Настраиваем graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Default().Info("сервер выключается")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Default().WithError(err).Fatal("произошла ошибка при выключении сервера")
	}
	logger.Sync()
}
