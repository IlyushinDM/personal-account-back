CREATE TABLE IF NOT EXISTS medical_center.admins (
    id bigserial PRIMARY KEY,
    login varchar(100) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    full_name varchar(255) NOT NULL,
    role varchar(50) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем GIN индекс для возможного будущего поиска по имени
CREATE INDEX IF NOT EXISTS idx_admins_full_name ON medical_center.admins USING gin(to_tsvector('russian', full_name));
