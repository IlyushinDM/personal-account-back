#!/usr/bin/env bash
# Скрипт для УДАЛЕНИЯ всех тестовых данных из БД (выполняет down-скрипт из папки data).
# Структура таблиц (схема) при этом НЕ затрагивается.
set -euo pipefail

echo "→ Загрузка конфигурации из .env..."
if [[ -f .env && -r .env ]]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
else
  echo "[ERROR] Файл конфигурации .env не найден." >&2
  exit 1
fi

echo "→ Удаление всех тестовых данных из БД..."

# Парсим URL, чтобы получить компоненты для psql
url="${DATABASE_URL}"; url="${url#*//}"; credentials="${url%%@*}"; url="${url#*@}"; DB_USER="${credentials%:*}"; DB_PASSWORD="${credentials#*:}"; host_port="${url%%/*}"; url="${url#*/}"; DB_HOST="${host_port%:*}"; DB_PORT="${host_port#*:}"; DB_NAME="${url%%\?*}"

# Принудительно устанавливаем кодировку клиента для psql
export PGCLIENTENCODING=UTF8
export PGPASSWORD="${DB_PASSWORD}"

if ! psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "${DB_NAME}" -f "migrations/data/000001_seed_initial_data.down.sql"; then
  echo "[ERROR] Не удалось выполнить скрипт очистки данных." >&2
  unset PGPASSWORD
  exit 1
fi
unset PGPASSWORD

echo "Все тестовые данные были успешно удалены из базы данных."