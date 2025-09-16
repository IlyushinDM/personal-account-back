#!/usr/bin/env bash
set -euo pipefail

# Проверка наличия утилиты migrate
if ! command -v migrate &>/dev/null; then
  echo "[ERROR] Команда 'migrate' не найдена" >&2
  exit 1
fi

# Загрузка переменных окружения
load_env() {
  if [[ -f .env && -r .env ]]; then
    set -a
    # shellcheck disable=SC1091
    . ./.env
    set +a
  fi
  
  # Устанавливаем дефолтные значения
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_NAME="${DB_NAME:-medical_center}"
  DB_SSLMODE="${DB_SSLMODE:-disable}"
}

# Проверка на заполненность переменных
check_vars() {
  : "${DB_USER?}" "${DB_PASSWORD?}" "${DB_HOST?}" "${DB_PORT?}" "${DB_NAME?}" "${DB_SSLMODE?}"
}

# Формирование строки подключения
build_db_url() {
  DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
}

# Проверка существования БД
db_exists() {
  PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -U "$DB_USER" -p "$DB_PORT" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"
}

# Откат миграций
rollback_migrations() {
  echo "[STEP] Откат всех миграций..."
  migrate -path "migrations" -database "$DB_URL" down -all
}

# Удаление схемы
drop_schema() {
  echo "[STEP] Удаление схемы..."
  migrate -path "migrations" -database "$DB_URL" drop -f
}

# Удаление самой БД
drop_database() {
  echo "[STEP] Удаление базы данных: $DB_NAME"
  PGPASSWORD="$DB_PASSWORD" dropdb -h "$DB_HOST" -U "$DB_USER" -p "$DB_PORT" "$DB_NAME"
}

# Удаление исполняемого файла
clean_binary() {
  echo "[STEP] Удаление исполняемого файла..."
  
  local bin_name="personal-account-back"
  case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac
  
  if [[ -f "bin/${bin_name}" ]]; then
    rm -f "bin/${bin_name}"
    echo "[DONE] Исполняемый файл удален: bin/${bin_name}"
  else
    echo "[INFO] Исполняемый файл не найден: bin/${bin_name}"
  fi
}

# === ОСНОВНОЙ ПОТОК ===

load_env
check_vars
build_db_url

if db_exists; then
  rollback_migrations
  drop_schema
  drop_database
  echo "[DONE] База данных полностью очищена и удалена"
else
  echo "[INFO] База данных $DB_NAME не существует"
fi

# Удаляем исполняемый файл в любом случае
clean_binary