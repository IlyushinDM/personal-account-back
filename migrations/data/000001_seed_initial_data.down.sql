-- TRUNCATE удаляет все данные и сбрасывает счетчики, CASCADE автоматически обрабатывает зависимости
TRUNCATE TABLE
    medical_center.cities,
    medical_center.departments,
    medical_center.specialties,
    medical_center.appointmentstatuses,
    medical_center.analysisstatuses,
    medical_center.clinics,
    medical_center.doctors,
    medical_center.services,
    medical_center.users,
    medical_center.user_profiles,
    medical_center.schedules,
    medical_center.legal_documents,
    medical_center.refresh_tokens
RESTART IDENTITY CASCADE;