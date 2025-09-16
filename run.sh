#!/usr/bin/env bash
set -euo pipefail

# Проверка наличия команды
has_cmd() { command -v "$1" >/dev/null 2>&1; }

# Загрузка конфигурации из .env
load_env_config() {
  if [[ -f .env && -r .env ]]; then
    echo "→ Загрузка конфигурации из .env"
    set -a
    # shellcheck disable=SC1091
    . ./.env
    set +a
    return 0
  fi
  return 1
}

# Загрузка конфигурации из переменных окружения с дефолтами
load_env_defaults() {
  echo "→ Использование переменных окружения с дефолтными значениями"
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_NAME="${DB_NAME:-medical_center}"
  DB_SSLMODE="${DB_SSLMODE:-disable}"
  return 0
}

# Загрузка конфигурации
load_config() {
  if load_env_config; then
    return 0
  else
    load_env_defaults
    return 0
  fi
}

# Проверка обязательных переменных
check_vars() {
  local missing=()

  [[ -z "${DB_USER:-}" ]] && missing+=("DB_USER")
  [[ -z "${DB_PASSWORD:-}" ]] && missing+=("DB_PASSWORD")
  [[ -z "${DB_HOST:-}" ]] && missing+=("DB_HOST")
  [[ -z "${DB_PORT:-}" ]] && missing+=("DB_PORT")
  [[ -z "${DB_NAME:-}" ]] && missing+=("DB_NAME")
  [[ -z "${DB_SSLMODE:-}" ]] && missing+=("DB_SSLMODE")

  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "[ERROR] Отсутствуют обязательные переменные: ${missing[*]}" >&2
    return 1
  fi
}

# Проверка наличия необходимых инструментов
check_tools() {
  local missing=()

  echo "→ Проверка необходимых инструментов..."

  if ! has_cmd go; then
    missing+=("go")
  fi

  if ! has_cmd migrate; then
    missing+=("migrate (установите github.com/golang-migrate/migrate/v4/cmd/migrate)")
  fi

  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "[ERROR] Отсутствуют необходимые инструменты:" >&2
    for tool in "${missing[@]}"; do
      echo "  - ${tool}" >&2
    done
    return 1
  fi
}

# Создание БД если не существует
create_db_if_not_exists() {
  if ! has_cmd psql; then
    echo "[WARNING] psql не найден. Пропускаю проверку/создание базы данных" >&2
    return 0
  fi

  echo "→ Проверка существования БД '${DB_NAME}'..."

  export PGPASSWORD="${DB_PASSWORD}"

  local exists
  if ! exists=$(psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -tAc \
    "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" 2>/dev/null); then
    echo "[ERROR] Не удалось подключиться к PostgreSQL серверу" >&2
    echo "  Проверьте параметры подключения: host=${DB_HOST}, port=${DB_PORT}, user=${DB_USER}" >&2
    unset PGPASSWORD
    return 1
  fi

  if [[ "${exists}" != "1" ]]; then
    echo "→ Создаю БД '${DB_NAME}'..."
    if ! psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -c "CREATE DATABASE \"${DB_NAME}\";" >/dev/null 2>&1; then
      echo "[ERROR] Не удалось создать базу данных '${DB_NAME}'" >&2
      unset PGPASSWORD
      return 1
    fi
    echo "  БД '${DB_NAME}' успешно создана"
  else
    echo "  БД '${DB_NAME}' уже существует"
  fi

  unset PGPASSWORD
}

# Применение миграций
apply_migrations() {
  local db_url="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

  echo "→ Применение миграций..."

  if ! migrate -path "migrations" -database "${db_url}" up; then
    echo "[ERROR] Не удалось применить миграции" >&2
    echo "  Проверьте:" >&2
    echo "  - Правильность строки подключения" >&2
    echo "  - Наличие директории 'migrations'" >&2
    echo "  - Доступность PostgreSQL сервера" >&2
    return 1
  fi

  echo "  Миграции успешно применены"
}

# Сборка бинарника
build_binary() {
  echo "→ Сборка бинарника..."

  if [[ ! -d "cmd/app" ]]; then
    echo "[ERROR] Директория 'cmd' не найдена" >&2
    return 1
  fi

  mkdir -p bin
  local bin_name="personal-account-back"

  # Для Windows
  case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac

  if ! go build -o "bin/${bin_name}" ./cmd/app; then
    echo "[ERROR] Не удалось собрать бинарник" >&2
    echo "  Проверьте корректность Go кода и зависимости" >&2
    return 1
  fi

  echo "  Бинарник успешно собран: bin/${bin_name}"
}

# === ОСНОВНОЙ ПОТОК ===

echo "[0/4] Подготовка к запуску..."
if ! check_tools; then
  exit 1
fi

echo "[1/4] Загрузка конфигурации..."
if ! load_config; then
  exit 1
fi

if ! check_vars; then
  exit 1
fi

echo "[2/4] Подготовка базы данных..."
if ! create_db_if_not_exists; then
  exit 1
fi

if ! apply_migrations; then
  exit 1
fi

echo "[3/4] Сборка бинарника..."
if ! build_binary; then
  exit 1
fi

echo "[4/4] Запуск приложения..."
bin_name="personal-account-back"
case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac

if [[ ! -f "bin/${bin_name}" ]]; then
  echo "[ERROR] Бинарник не найден: bin/${bin_name}" >&2
  exit 1
fi

echo "→ Приложение запущено"
exec "bin/${bin_name}"