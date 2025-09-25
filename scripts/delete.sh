#!/usr/bin/env bash
# Скрипт для ПОЛНОГО УДАЛЕНИЯ базы данных и бинарного файла проекта.
# ! ВНИМАНИЕ: Это действие необратимо и приведет к потере всех данных в БД.

# Подключаем общие функции
# shellcheck disable=SC1091
source "scripts/_helpers.sh"

# --- ФУНКЦИИ СКРИПТА ---

check_tools() {
    local missing=()
    echo ">> Проверка необходимых инструментов..."

    if ! has_cmd migrate; then missing+=("migrate (установите с тегом postgres)"); fi
    if ! has_cmd psql; then missing+=("psql (клиент PostgreSQL)"); fi

    if [[ ${#missing[@]} -gt 0 ]]; then
        echo "[ERROR] Отсутствуют необходимые инструменты:" >&2
        for tool in "${missing[@]}"; do echo "  - ${tool}" >&2; done
        return 1
    fi
}

db_exists() {
  export PGPASSWORD="${DB_PASSWORD:-}"
  if ! psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -c '\q' 2>/dev/null; then
    echo "[WARNING] Не удалось подключиться к серверу PostgreSQL. Возможно, база уже удалена или сервер не запущен."
    unset PGPASSWORD
    return 1
  fi
  
  local exists
  exists=$(psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d postgres -tAc \
    "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'")
  unset PGPASSWORD
  
  [[ "${exists}" == "1" ]]
}

drop_database() {
  echo ">> Удаление базы данных '${DB_NAME}'..."
  if ! migrate -path "migrations/schema" -database "${DATABASE_URL}" drop -f; then
    echo "[ERROR] Не удалось удалить базу данных с помощью migrate." >&2
    return 1
  fi
  echo "  База данных '${DB_NAME}' успешно удалена."
}

clean_binary() {
  echo ">> Удаление исполняемого файла..."
  
  local bin_path="bin/${BINARY_NAME}"
  case "${MSYSTEM:-}" in MINGW*|MSYS*) bin_path="${bin_path}.exe";; esac
  
  if [[ -f "${bin_path}" ]]; then
    rm -f "${bin_path}"
    echo "  Исполняемый файл удален: ${bin_path}"
  else
    echo "  Исполняемый файл не найден (${bin_path})"
  fi
}

# --- ОСНОВНОЙ ПОТОК ВЫПОЛНЕНИЯ ---

echo "--- [0/2] Подготовка к очистке проекта ---"
check_tools
load_config

echo "--- [1/2] Очистка базы данных ---"
if db_exists; then
  drop_database
else
  echo ">> База данных '${DB_NAME}' не существует или сервер недоступен. Пропускаем удаление."
fi

echo "--- [2/2] Очистка артефактов сборки ---"
clean_binary

echo "--- Очистка проекта успешно завершена ---"