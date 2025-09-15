CREATE TABLE IF NOT EXISTS medical_center.doctorclinics (
	doctor_id bigint NOT NULL,
	clinic_id bigint NOT NULL,
	PRIMARY KEY (doctor_id, clinic_id),
	CONSTRAINT doctorclinics_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT doctorclinics_clinic_id_fkey FOREIGN KEY (clinic_id)
		REFERENCES medical_center.clinics(clinic_id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_doctor_clinics_clinic_id ON medical_center.doctorclinics(clinic_id);
