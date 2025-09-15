CREATE TABLE IF NOT EXISTS medical_center.doctorcourses (
	course_id bigserial PRIMARY KEY,
	doctor_id bigint NOT NULL,
	course_name varchar(255) NOT NULL,
	year smallint NOT NULL,
	CONSTRAINT doctorcourses_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE CASCADE
);
