CREATE TABLE IF NOT EXISTS medical_center.services (
	service_id bigserial PRIMARY KEY,
	name varchar(255) NOT NULL,
	price numeric(10,2) NOT NULL,
	duration_minutes smallint NOT NULL,
	description text,
	doctor_id bigint NOT NULL,
	CONSTRAINT services_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_services_doctor_id ON medical_center.services(doctor_id);
