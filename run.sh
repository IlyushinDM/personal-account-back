#!/usr/bin/env bash
# Скрипт для полной подготовки и запуска REST API "Личный кабинет".
#
# Использование:
#   ./run.sh          - Запускает проект с чистой базой данных (только структура).
#   ./run.sh --seed   - Запускает проект и заполняет базу тестовыми данными.
#
set -euo pipefail

# --- ПАРСИНГ АРГУМЕНТОВ ---
SEED_DATA=false
for arg in "$@"; do
  case $arg in
    --seed)
      SEED_DATA=true
      shift # Убираем обработанный флаг из списка аргументов
      ;;
  esac
done

# --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---

# Проверка наличия команды в PATH
has_cmd() {
  command -v "$1" >/dev/null 2>&1
}

# Функция для разбора строки DATABASE_URL на компоненты
parse_db_url() {
  local url="${DATABASE_URL}"
  url="${url#*//}"
  local credentials
  credentials="${url%%@*}"
  url="${url#*@}"
  DB_USER="${credentials%:*}"
  DB_PASSWORD="${credentials#*:}"
  local host_port
  host_port="${url%%/*}"
  url="${url#*/}"
  DB_HOST="${host_port%:*}"
  DB_PORT="${host_port#*:}"
  DB_NAME="${url%%\?*}"
  
  if [[ -z "${DB_USER}" || -z "${DB_PASSWORD}" || -z "${DB_HOST}" || -z "${DB_PORT}" || -z "${DB_NAME}" ]]; then
    echo "[ERROR] Не удалось разобрать DATABASE_URL. Проверьте формат в .env" >&2
    return 1
  fi
}

# Загрузка конфигурации из .env файла
load_env_config() {
  if [[ -f .env && -r .env ]]; then
    set -a
    # shellcheck disable=SC1091
    . ./.env
    set +a
  else
    echo "[ERROR] Файл конфигурации .env не найден." >&2
    return 1
  fi
}

# Основная функция загрузки конфигурации
load_config() {
  load_env_config
  parse_db_url
}

# Проверка наличия необходимых инструментов
check_tools() {
  local missing=()
  echo ">> Проверка необходимых инструментов..."

  if ! has_cmd go; then missing+=("go"); fi
  if ! has_cmd migrate; then missing+=("migrate (установите с тегом postgres)"); fi
  if ! has_cmd psql; then missing+=("psql (клиент PostgreSQL)"); fi
  if ! has_cmd swag; then missing+=("swag (документация API)"); fi

  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "[ERROR] Отсутствуют необходимые инструменты:" >&2
    for tool in "${missing[@]}"; do echo "  - ${tool}" >&2; done
    return 1
  fi
}

# Создание базы данных, если она не существует
create_db_if_not_exists() {
  echo ">> Проверка существования БД '${DB_NAME}'..."
  export PGPASSWORD="${DB_PASSWORD}"

  if ! psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -c '\q' 2>/dev/null; then
    echo "[ERROR] Не удалось подключиться к серверу PostgreSQL. Проверьте DATABASE_URL и доступность сервера." >&2
    unset PGPASSWORD
    return 1
  fi
  
  if [[ "$(psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'")" != "1" ]]; then
    echo ">> Создаю БД '${DB_NAME}'..."
    if ! createdb -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" "${DB_NAME}"; then
      echo "[ERROR] Не удалось создать базу данных '${DB_NAME}'" >&2
      unset PGPASSWORD
      return 1
    fi
    echo "  БД '${DB_NAME}' успешно создана."
  else
    echo "  БД '${DB_NAME}' уже существует."
  fi
  unset PGPASSWORD
}

# Применение миграций схемы и опциональное заполнение данными
apply_migrations_and_seed() {
  echo "→ Применение миграций СХЕМЫ..."
  if ! migrate -path "migrations/schema" -database "${DATABASE_URL}" up; then
    echo "[ERROR] Не удалось применить миграции схемы." >&2
    return 1
  fi
  echo "  Миграции схемы успешно применены."

  if [ "$SEED_DATA" = true ]; then
    echo "→ Заполнение БД тестовыми данными (обнаружен флаг --seed)..."
    
    # Принудительно устанавливаем кодировку клиента для psql
    export PGCLIENTENCODING=UTF8
    export PGPASSWORD="${DB_PASSWORD}"

    if ! psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "${DB_NAME}" -f "migrations/data/000001_seed_initial_data.up.sql"; then
      echo "[ERROR] Не удалось выполнить скрипт заполнения данными." >&2
      unset PGPASSWORD
      return 1
    fi
    unset PGPASSWORD
    echo "  Тестовые данные успешно загружены в БД."
  else
    echo "→ Заполнение тестовыми данными пропущено. Используйте './run.sh --seed' для заполнения."
  fi
}

# Генерация документации Swagger
generate_swagger_docs() {
    echo ">> Генерация документации Swagger..."
    if ! swag init -g cmd/app/main.go --parseDependency --parseInternal; then
        echo "[ERROR] Не удалось сгенерировать документацию Swagger." >&2
        return 1
    fi
    echo "  Документация успешно сгенерирована в папку ./docs"
}

# Сборка Go-приложения
build_binary() {
  echo ">> Сборка бинарника..."
  mkdir -p bin
  local bin_name="personal-account-back"
  case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac

  if ! go build -o "bin/${bin_name}" ./cmd/app; then
    echo "[ERROR] Сборка приложения не удалась." >&2
    return 1
  fi
  echo "  Бинарник успешно собран: bin/${bin_name}"
}

# --- ОСНОВНОЙ ПОТОК ВЫПОЛНЕНИЯ ---
echo "--- [0/5] Подготовка к запуску ---"
if ! check_tools; then exit 1; fi

echo "--- [1/5] Загрузка конфигурации ---"
if ! load_config; then exit 1; fi

echo "--- [2/5] Подготовка базы данных ---"
if ! create_db_if_not_exists; then exit 1; fi
if ! apply_migrations_and_seed; then exit 1; fi

echo "--- [3/5] Генерация кода и документации ---"
if ! generate_swagger_docs; then exit 1; fi

echo "--- [4/5] Сборка приложения ---"
if ! build_binary; then exit 1; fi

echo "--- [5/5] Запуск приложения ---"
bin_name="personal-account-back"
case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac

echo "========================================="
echo ">> Сервер запущен: http://localhost:${HTTP_PORT:-8080}"
echo ">> Документация Swagger: http://localhost:${HTTP_PORT:-8080}/swagger/index.html"
echo "Для остановки нажмите Ctrl+C"
echo "========================================="

exec "bin/${bin_name}"