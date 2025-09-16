CREATE TABLE IF NOT EXISTS medical_center.clinics (
	clinic_id bigserial PRIMARY KEY,
	name varchar(150) NOT NULL,
	address varchar(512) NOT NULL,
	work_hours varchar(100) NOT NULL,
	phone varchar(20) NOT NULL
);
