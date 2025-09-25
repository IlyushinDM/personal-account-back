#!/usr/bin/env bash
#* Для run.sh seed нужно будет положить 000001_seed_initial_data.up.sql в папку, 
#* которую PostgreSQL сможет увидеть при старте, либо использовать psql, но это потребует psql на хосте. 
#* -!- Пока используется run-light -!-
# Скрипт-обертка для управления Docker Compose окружением.

set -euo pipefail

# Функция для вывода помощи
show_help() {
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  up          Builds and starts all services."
    echo "  up -d       Builds and starts all services in detached mode."
    echo "  down        Stops and removes all services and networks."
    echo "  logs        Follows logs from all services."
    echo "  logs [app]  Follows logs from a specific service (e.g., app, postgres)."
    echo "  migrate     Applies all database migrations."
    echo "  seed        Applies migrations and seeds the database with test data."
    echo "  clean       Stops services and removes all data volumes."
    echo ""
}

# Основная логика
case "${1:-}" in
    up)
        echo "Starting up services..."
        # Передаем все аргументы после 'up' (например, -d) в docker-compose
        docker-compose up --build "${@:2}"
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
        docker-compose exec -it app sh -c "migrate -path /migrations/schema -database \"\$$DATABASE_URL\" up"
        ;;
    seed)
        echo "Applying migrations and seeding data..."
        docker-compose exec -it app sh -c "migrate -path /migrations/schema -database \"\$$DATABASE_URL\" up"
        docker-compose exec -it postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/000001_seed_initial_data.up.sql
        echo "Data seeded successfully."
        ;;
    clean)
        echo "Stopping services and removing all data..."
        docker-compose down -v
        ;;
    *)
        show_help
        ;;
esac