#!/usr/bin/env bash
# Скрипт для ПОЛНОГО УДАЛЕНИЯ базы данных и бинарного файла проекта.
# ! ВНИМАНИЕ: Это действие необратимо и приведет к потере всех данных в БД.
set -euo pipefail

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
    echo "→ Загрузка конфигурации из .env"
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
    echo "→ Проверка необходимых инструментов..."

    if ! has_cmd migrate; then missing+=("migrate (установите с тегом postgres)"); fi
    if ! has_cmd psql; then missing+=("psql (клиент PostgreSQL)"); fi

    if [[ ${#missing[@]} -gt 0 ]]; then
        echo "[ERROR] Отсутствуют необходимые инструменты:" >&2
        for tool in "${missing[@]}"; do echo "  - ${tool}" >&2; done
        return 1
    fi
}

# Проверка, существует ли база данных
db_exists() {
  export PGPASSWORD="${DB_PASSWORD}"
  # Проверяем, можно ли вообще подключиться к серверу
  if ! psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -c '\q' 2>/dev/null; then
    echo "[WARNING] Не удалось подключиться к серверу PostgreSQL. Возможно, база уже удалена или сервер не запущен."
    unset PGPASSWORD
    return 1 # Код возврата 1 означает "не существует" или "недоступна"
  fi
  
  local exists
  exists=$(psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -tAc \
    "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'")
  unset PGPASSWORD
  
  [[ "${exists}" == "1" ]]
}

# Удаление базы данных
drop_database() {
  echo "→ Удаление базы данных '${DB_NAME}'..."
  # Указываем путь к миграциям СХЕМЫ, так как там создается таблица schema_migrations,
  # на которую ориентируется команда drop.
  if ! migrate -path "migrations/schema" -database "${DATABASE_URL}" drop -f; then
    echo "[ERROR] Не удалось удалить базу данных с помощью migrate." >&2
    return 1
  fi
  echo "  База данных '${DB_NAME}' успешно удалена."
}

# Удаление скомпилированного бинарника
clean_binary() {
  echo "→ Удаление исполняемого файла..."
  
  local bin_name="personal-account-back"
  case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_name="${bin_name}.exe";; esac
  
  if [[ -f "bin/${bin_name}" ]]; then
    rm -f "bin/${bin_name}"
    echo "  Исполняемый файл удален: bin/${bin_name}"
  else
    echo "  Исполняемый файл не найден (bin/${bin_name})"
  fi
}

# --- ОСНОВНОЙ ПОТОК ВЫПОЛНЕНИЯ ---

echo "--- [0/2] Подготовка к очистке проекта ---"
if ! check_tools; then exit 1; fi
if ! load_config; then exit 1; fi

echo "--- [1/2] Очистка базы данных ---"
if db_exists; then
  drop_database
else
  echo "→ База данных '${DB_NAME}' не существует или сервер недоступен. Пропускаем удаление."
fi

echo "--- [2/2] Очистка артефактов сборки ---"
clean_binary

echo "========================================="
echo "  Очистка проекта успешно завершена."
echo "========================================="