CREATE TABLE IF NOT EXISTS medical_center.doctorskills (
	skill_id bigserial PRIMARY KEY,
	doctor_id bigint NOT NULL,
	skill_name varchar(255) NOT NULL,
	CONSTRAINT doctorskills_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);
