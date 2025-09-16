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

	"lk/internal/config"
	"lk/internal/handlers"
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

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Используйте "Bearer " перед вашим токеном."
func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Файл .env не найден, используются системные переменные окружения")
	}

	cfg := config.MustLoad()

	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		fmt.Printf("не удалось подключиться к базе данных: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("соединение с базой данных установлено")

	repos := repository.NewRepository(db)
	serviceDeps := services.ServiceDependencies{
		Repos:      repos,
		SigningKey: cfg.Auth.JWTSecretKey,
		TokenTTL:   cfg.Auth.TokenTTL,
	}
	services := services.NewService(serviceDeps)
	handlers := handlers.NewHandler(services)
	fmt.Println("слои приложения инициализированы")

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.HTTPServer.Port, handlers.InitRoutes()); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("ошибка при запуске http-сервера: %v\n", err)
		}
	}()
	fmt.Printf("сервер запущен на порту: %s\n", cfg.HTTPServer.Port)

	// Плавное завершение работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("сервер выключается")
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("произошла ошибка при выключении сервера: %v\n", err)
	}
}
