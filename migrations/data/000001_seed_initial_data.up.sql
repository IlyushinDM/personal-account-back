-- ! ПРОВЕРИТЬ ДАННЫЙ СКРИПТ НА АКТУАЛЬНОСТЬ ПЕРЕД ПРИМЕНЕНИЕМ ! 
-- Заполняем справочники
INSERT INTO medical_center.cities (id, name) VALUES (1, 'Санкт-Петербург'), (2, 'Москва') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.departments (id, name) VALUES (1, 'Терапия'), (2, 'Хирургия'), (3, 'Диагностика') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.specialties (id, name, department_id) VALUES (1, 'Терапевт', 1), (2, 'Кардиолог', 1), (3, 'Хирург', 2), (4, 'Врач УЗИ', 3) ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.appointmentstatuses (id, name) VALUES (1, 'Запланировано'), (2, 'Завершено'), (3, 'Отменено пациентом'), (4, 'Отменено клиникой') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.analysisstatuses (id, name) VALUES (1, 'Назначено'), (2, 'В работе'), (3, 'Готов') ON CONFLICT (id) DO NOTHING;

-- Создаем клинику
INSERT INTO medical_center.clinics (id, name, address, work_hours, phone) VALUES
(1, 'Клиника "Здоровье"', 'г. Санкт-Петербург, Невский пр., д. 1', 'Пн-Пт 08:00-20:00', '+78121234567') ON CONFLICT (id) DO NOTHING;

-- Создаем тестового доктора
INSERT INTO medical_center.doctors (id, first_name, last_name, patronymic, specialty_id, experience_years, rating, review_count, avatar_url) VALUES
(1, 'Иван', 'Иванов', 'Иванович', 2, 15, 4.8, 45, '/avatars/ivanov.jpg') ON CONFLICT (id) DO NOTHING;

-- Создаем услугу для этого доктора
INSERT INTO medical_center.services (id, name, price, duration_minutes, description, doctor_id) VALUES
(1, 'Консультация кардиолога', 2500.00, 30, 'Первичная консультация ведущего кардиолога.', 1) ON CONFLICT (id) DO NOTHING;

-- Создаем тестового пользователя
-- Пароль: 'password123'
INSERT INTO medical_center.users (id, phone, password_hash) VALUES
(1, '+79991234567', 'здесь нужно вставить токен') ON CONFLICT (id) DO NOTHING;  -- TODO: вставить токен!

-- Создаем профиль для пользователя (в правильную таблицу user_profiles)
INSERT INTO medical_center.user_profiles (id, user_id, first_name, last_name, patronymic, birth_date, gender, city_id, email) VALUES
(1, 1, 'Петр', 'Петров', 'Петрович', '1990-05-10', 'male', 1, 'test@example.com') ON CONFLICT (id) DO NOTHING;

-- Сбрасываем последовательности, чтобы новые записи начинались с корректных ID
SELECT setval('medical_center.cities_id_seq', (SELECT MAX(id) FROM medical_center.cities));
SELECT setval('medical_center.departments_id_seq', (SELECT MAX(id) FROM medical_center.departments));
SELECT setval('medical_center.specialties_id_seq', (SELECT MAX(id) FROM medical_center.specialties));
SELECT setval('medical_center.doctors_id_seq', (SELECT MAX(id) FROM medical_center.doctors));
SELECT setval('medical_center.services_id_seq', (SELECT MAX(id) FROM medical_center.services));
SELECT setval('medical_center.users_id_seq', (SELECT MAX(id) FROM medical_center.users));
SELECT setval('medical_center.user_profiles_id_seq', (SELECT MAX(id) FROM medical_center.user_profiles));