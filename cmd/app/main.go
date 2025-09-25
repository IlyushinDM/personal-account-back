// Package main является точкой входа для серверного приложения "Личный Кабинет" (lk).
// Он инициализирует конфигурацию, базу данных, все слои приложения,
// а затем запускает HTTP-сервер с корректным завершением работы.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"lk/internal/config"
	"lk/internal/logger"
	"lk/internal/models"
	"lk/internal/repository"
	"lk/internal/server"
	"lk/internal/services"
	"lk/internal/storage"
	httptransport "lk/internal/transport/http"
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
	// 1. Инициализация логгера
	logger.Init(os.Getenv("LOG_DIR"))
	logger.Default().Info("логгер инициализирован")

	cfg := config.MustLoad()

	location, err := time.LoadLocation(cfg.ClinicTimezone)
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось загрузить указанный часовой пояс")
	}
	logger.Default().Info(fmt.Sprintf("приложение работает в часовом поясе: %s", location.String()))

	// 2. Инициализация клиентов к внешним системам (DB, Redis, S3)
	gormDB, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.NewGORMLogger(),
	})
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к базе данных")
	}
	db, err := gormDB.DB()
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось получить *sql.DB из gorm")
	}
	defer db.Close()
	logger.Default().Info("соединение с базой данных установлено")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к Redis")
	}
	defer redisClient.Close()
	logger.Default().Info("соединение с Redis установлено")

	storageClient, err := storage.NewMinIOClient(cfg.Minio)
	if err != nil {
		logger.Default().WithError(err).Fatal("не удалось подключиться к MinIO")
	}
	logger.Default().Info("соединение с MinIO установлено")

	// 3. Dependency Injection: собираем все зависимости
	repos := repository.NewRepository(gormDB, redisClient)
	serviceDeps := services.ServiceDependencies{
		Repos:      repos,
		Storage:    storageClient,
		Location:   location,
		SigningKey: cfg.Auth.JWTSecretKey,
		TokenTTL:   cfg.Auth.TokenTTL,
	}
	services := services.NewService(serviceDeps)

	handler := httptransport.NewHandler(services, repos.User)
	logger.Default().Info("слои приложения инициализированы")

	// 4. Инициализация HTTP-роутера и middleware
	router := gin.New()
	router.Use(logger.GinLogger(), gin.Recovery(), httptransport.ErrorMiddleware())
	handler.InitRoutes(router)

	// 5. Запуск HTTP-сервера с Graceful Shutdown
	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPServer.Port, router); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Default().WithError(err).Fatal("ошибка при запуске http-сервера")
		}
	}()
	logger.Default().Info(fmt.Sprintf("сервер запущен на порту: %s", cfg.HTTPServer.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Default().Info("сервер выключается")
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Default().WithError(err).Fatal("произошла ошибка при выключении сервера")
	}
	logger.Sync()
}
