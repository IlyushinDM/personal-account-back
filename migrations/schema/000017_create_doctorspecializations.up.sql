CREATE TABLE IF NOT EXISTS medical_center.doctorspecializations (
	id bigserial PRIMARY KEY,
	doctor_id bigint NOT NULL,
	area varchar(255) NOT NULL,
	CONSTRAINT doctorspecializations_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);
