CREATE TABLE IF NOT EXISTS medical_center.doctorcertificates (
	id bigserial PRIMARY KEY,
	doctor_id bigint NOT NULL,
	cert_name varchar(255) NOT NULL,
	cert_number varchar(100),
	CONSTRAINT doctorcertificates_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);
