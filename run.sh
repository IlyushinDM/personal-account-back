#!/usr/bin/env bash
# Скрипт-обертка для управления Docker Compose окружением.

set -euo pipefail

# --- ЗАГРУЗКА .env ---
# Проверяем, существует ли .env файл, и загружаем его переменные
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "[ERROR] .env file not found. Please create it before running."
  exit 1
fi
# --- КОНЕЦ БЛОКА .env ---

# Функция для вывода помощи
show_help() {
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  up          Builds, migrates, and starts all services. Use Ctrl+C to stop."
    echo "  down        Stops and removes all services and networks."
    echo "  logs        (In another terminal) Follows logs from all services."
    echo "  migrate     (In another terminal) Applies all database migrations."
    echo "  clean       Stops services and removes all data volumes."
    echo ""
}

# Основная логика
case "${1:-}" in
    up)
        echo "Starting up services..."
        # Запускаем миграции ПЕРЕД запуском приложения
        echo "Applying migrations..."
        # Теперь переменные доступны, и строка подключения соберется корректно
        migrate -path migrations/schema -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up
        
        echo "Starting docker-compose..."
        docker-compose up --build
        ;;
    down)
        echo "Stopping services..."
        docker-compose down
        ;;
    logs)
        echo "Following logs..."
        docker-compose logs -f "${2:-}"
        ;;
    migrate)
         echo "Applying migrations..."
         migrate -path migrations/schema -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up
        ;;
    clean)
        echo "Stopping services and removing all data..."
        docker-compose down -v
        ;;
    *)
        show_help
        ;;
esac