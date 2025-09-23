-- Добавляем колонку для хранения tsvector
ALTER TABLE medical_center.doctors
ADD COLUMN IF NOT EXISTS fts_document tsvector;

-- Создаем функцию, которая будет обновлять tsvector-колонку
CREATE OR REPLACE FUNCTION medical_center.update_doctor_fts_document()
RETURNS TRIGGER AS $$
DECLARE
    specialty_name text;
BEGIN
    -- Получаем название специальности
    SELECT name INTO specialty_name FROM medical_center.specialties WHERE id = NEW.specialty_id;

    -- Объединяем ФИО и специальность, обрабатывая возможные NULL значения,
    -- и преобразуем в tsvector с использованием словаря для русского языка.
    NEW.fts_document := to_tsvector('russian',
        COALESCE(NEW.last_name, '') || ' ' ||
        COALESCE(NEW.first_name, '') || ' ' ||
        COALESCE(NEW.patronymic, '') || ' ' ||
        COALESCE(specialty_name, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер, который вызывает функцию перед вставкой или обновлением
DROP TRIGGER IF EXISTS tsvector_update ON medical_center.doctors;
CREATE TRIGGER tsvector_update
BEFORE INSERT OR UPDATE ON medical_center.doctors
FOR EACH ROW EXECUTE FUNCTION medical_center.update_doctor_fts_document();

-- Заполняем fts_document для уже существующих записей
-- Важно выполнить это после создания триггера, чтобы функция уже существовала
UPDATE medical_center.doctors SET fts_document = to_tsvector('russian',
    COALESCE(last_name, '') || ' ' ||
    COALESCE(first_name, '') || ' ' ||
    COALESCE(patronymic, '') || ' ' ||
    (SELECT name FROM medical_center.specialties WHERE id = doctors.specialty_id)
);

-- Создаем GIN-индекс для быстрой работы FTS
CREATE INDEX IF NOT EXISTS idx_doctors_fts ON medical_center.doctors USING gin(fts_document);