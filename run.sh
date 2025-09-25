#!/usr/bin/env bash
# Скрипт-обертка для управления Docker Compose окружением.

set -euo pipefail

# --- Загрузка .env ---
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "[ERROR] .env file not found. Please create it before running."
  exit 1
fi

# Функция для вывода помощи
show_help() {
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  up          Generates docs, builds, migrates, and starts all services."
    echo "  down        Stops and removes all services and networks."
    echo "  logs        (In another terminal) Follows logs from a service (e.g., app)."
    echo "  seed        (After 'up') Seeds the database with test data."
    echo "  swag        Manually regenerate Swagger docs."
    echo "  clean       Stops services and removes all data volumes."
    echo ""
}

# Функция для генерации Swagger
generate_swag() {
    echo "--- [1/3] Generating Swagger docs ---"
    swag init --parseDependency --parseInternal -g ./cmd/app/main.go
}


# Основная логика
case "${1:-}" in
    up)
        generate_swag
        
        echo "--- [2/3] Starting up services (this will also apply migrations)..."
        docker-compose up --build
        
        echo "--- [3/3] Application is running ---"
        ;;
    down)
        echo "Stopping services..."
        docker-compose down
        ;;
    logs)
        echo "Following logs for '${2:-all services}'..."
        docker-compose logs -f "${2:-}"
        ;;
    seed)
        echo "Seeding data..."
        docker-compose exec -T postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" < ./migrations/data/000001_seed_initial_data.up.sql
        echo "Data seeded successfully."
        ;;
    swag)
        generate_swag
        ;;
    clean)
        echo "Stopping services and removing all data..."
        docker-compose down -v
        ;;
    *)
        show_help
        ;;
esac