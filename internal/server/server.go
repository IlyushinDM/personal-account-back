// Package server инкапсулирует логику запуска и корректного завершения работы HTTP-сервера,
// отделяя управление жизненным циклом сервера от основной логики инициализации приложения.
package server

import (
	"context"
	"net/http"
	"time"
)

// Server представляет наш HTTP-сервер.
type Server struct {
	httpServer *http.Server
}

// Run запускает HTTP-сервер на указанном порту с переданным обработчиком (роутером).
// Метод настраивает сервер с адекватными таймаутами для защиты от атак типа Slowloris.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,          // 1 MB - защита от атак с большими заголовками.
		ReadTimeout:    10 * time.Second, // Защита от Slowloris-атак.
		WriteTimeout:   10 * time.Second, // Защита от Slowloris-атак.
		IdleTimeout:    1 * time.Minute,  // Позволяет keep-alive соединениям закрываться.
	}

	// Запуск сервера. Эта функция блокирующая.
	return s.httpServer.ListenAndServe()
}

// Shutdown корректно останавливает HTTP-сервер (graceful shutdown).
// Он дожидается завершения текущих запросов перед остановкой.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
