-- Заполняем справочники
INSERT INTO medical_center.cities (id, name) VALUES (1, 'Санкт-Петербург'), (2, 'Москва') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.departments (id, name) VALUES (1, 'Терапия'), (2, 'Хирургия'), (3, 'Диагностика') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.specialties (id, name, department_id) VALUES (1, 'Терапевт', 1), (2, 'Кардиолог', 1), (3, 'Хирург', 2), (4, 'Врач УЗИ', 3) ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.appointmentstatuses (id, name) VALUES (1, 'Запланировано'), (2, 'Завершено'), (3, 'Отменено пациентом'), (4, 'Отменено клиникой') ON CONFLICT (id) DO NOTHING;
INSERT INTO medical_center.analysisstatuses (id, name) VALUES (1, 'Назначено'), (2, 'В работе'), (3, 'Готов') ON CONFLICT (id) DO NOTHING;

-- Создаем клинику
INSERT INTO medical_center.clinics (id, name, address, work_hours, phone) VALUES
(1, 'Клиника "Здоровье"', 'г. Санкт-Петербург, Невский пр., д. 1', 'Пн-Пт 08:00-20:00', '+78121234567') ON CONFLICT (id) DO NOTHING;

-- Создаем тестовых докторов
INSERT INTO medical_center.doctors (id, first_name, last_name, patronymic, specialty_id, experience_years, rating, review_count, avatar_url, recommendations) VALUES
(1, 'Иван', 'Иванов', 'Иванович', 2, 15, 4.8, 45, '/avatars/ivanov.jpg', 'Пожалуйста, приходите за 15 минут до приема и возьмите с собой паспорт и результаты предыдущих обследований.'),
(2, 'Мария', 'Сергеева', 'Павловна', 1, 8, 4.9, 62, '/avatars/sergeeva.jpg', 'Рекомендуется не есть за 2 часа до приема. Пить воду можно.')
ON CONFLICT (id) DO NOTHING;

-- Создаем услуги для докторов
INSERT INTO medical_center.services (id, name, price, duration_minutes, description, doctor_id, recommendations) VALUES
(1, 'Консультация кардиолога', 2500.00, 30, 'Первичная консультация ведущего кардиолога.', 1, 'При себе иметь кардиограмму (ЭКГ), сделанную не позднее месяца назад.'),
(2, 'Первичный прием терапевта', 1800.00, 20, 'Осмотр, сбор анамнеза, назначение лечения.', 2, 'Вспомните все препараты, которые вы принимаете на постоянной основе.')
ON CONFLICT (id) DO NOTHING;

-- Создаем расписание для доктора Иванова (id=1)
INSERT INTO medical_center.schedules (doctor_id, date, start_time, end_time) VALUES
(1, '2025-09-25', '09:00:00', '17:00:00'),
(1, '2025-09-26', '09:00:00', '17:00:00'),
(1, '2025-09-29', '10:00:00', '15:00:00')
ON CONFLICT (doctor_id, date) DO NOTHING;

-- Создаем тестового пользователя
-- Пароль: 'password123'
INSERT INTO medical_center.users (id, phone, password_hash) VALUES
(1, '+79991234567', '$2a$10$5Aa9vplKAviMqZnl/4bOhutXhHiJjrNPN3AKefBsagxnuLejYzNcS') ON CONFLICT (id) DO NOTHING;

-- Создаем профиль для пользователя
INSERT INTO medical_center.user_profiles (id, user_id, first_name, last_name, patronymic, birth_date, gender, city_id, email) VALUES
(1, 1, 'Петр', 'Петров', 'Петрович', '1990-05-10', 'male', 1, 'test@example.com') ON CONFLICT (id) DO NOTHING;

-- Создаем юридические документы
INSERT INTO medical_center.legal_documents (type, title, url, version, update_date) VALUES
('terms_of_use', 'Пользовательское соглашение', '/legal/terms.pdf', '1.2', '2023-10-01'),
('privacy_policy', 'Политика конфиденциальности', '/legal/privacy.pdf', '2.0', '2023-09-15')
ON CONFLICT (id) DO NOTHING;


-- Сбрасываем последовательности, чтобы новые записи начинались с корректных ID
SELECT setval('medical_center.cities_id_seq', (SELECT MAX(id) FROM medical_center.cities), true);
SELECT setval('medical_center.departments_id_seq', (SELECT MAX(id) FROM medical_center.departments), true);
SELECT setval('medical_center.specialties_id_seq', (SELECT MAX(id) FROM medical_center.specialties), true);
SELECT setval('medical_center.doctors_id_seq', (SELECT MAX(id) FROM medical_center.doctors), true);
SELECT setval('medical_center.services_id_seq', (SELECT MAX(id) FROM medical_center.services), true);
SELECT setval('medical_center.users_id_seq', (SELECT MAX(id) FROM medical_center.users), true);
SELECT setval('medical_center.user_profiles_id_seq', (SELECT MAX(id) FROM medical_center.user_profiles), true);
SELECT setval('medical_center.schedules_id_seq', (SELECT MAX(id) FROM medical_center.schedules), true);
SELECT setval('medical_center.legal_documents_id_seq', (SELECT MAX(id) FROM medical_center.legal_documents), true);