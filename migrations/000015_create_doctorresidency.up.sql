CREATE TABLE IF NOT EXISTS medical_center.doctorresidency (
	residency_id bigserial PRIMARY KEY,
	doctor_id bigint NOT NULL,
	institution varchar(255) NOT NULL,
	specialty varchar(255) NOT NULL,
	start_year smallint NOT NULL,
	end_year smallint NOT NULL,
	CONSTRAINT doctorresidency_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);
