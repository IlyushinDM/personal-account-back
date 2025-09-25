#!/usr/bin/env bash
# Скрипт для УДАЛЕНИЯ всех тестовых данных из БД (выполняет down-скрипт из папки data).
# Структура таблиц (схема) при этом НЕ затрагивается.

# Подключаем общие функции
# shellcheck disable=SC1091
source "scripts/_helpers.sh"

# --- ОСНОВНОЙ ПОТОК ВЫПОЛНЕНИЯ ---
load_config

echo ">> Удаление всех тестовых данных из БД..."

# Принудительно устанавливаем кодировку клиента для psql
export PGCLIENTENCODING=UTF8
# Переменные DB_* уже установлены в load_config
export PGPASSWORD="${DB_PASSWORD:-}"

if ! psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "${DB_NAME}" -f "migrations/data/000001_seed_initial_data.down.sql"; then
  echo "[ERROR] Не удалось выполнить скрипт очистки данных." >&2
  unset PGPASSWORD
  exit 1
fi
unset PGPASSWORD

echo ">> Все тестовые данные были успешно удалены из базы данных."