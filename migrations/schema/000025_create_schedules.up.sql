CREATE TABLE IF NOT EXISTS medical_center.schedules (
    id bigserial PRIMARY KEY,
    doctor_id bigint NOT NULL,
    date date NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone NOT NULL,

    CONSTRAINT fk_schedules_doctor_id FOREIGN KEY (doctor_id)
        REFERENCES medical_center.doctors(id)
        ON DELETE CASCADE,

    -- У одного доктора может быть только одна запись на одну дату
    UNIQUE (doctor_id, date)
);

CREATE INDEX IF NOT EXISTS idx_schedules_doctor_id_date ON medical_center.schedules(doctor_id, date);