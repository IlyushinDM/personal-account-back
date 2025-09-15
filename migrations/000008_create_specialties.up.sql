CREATE TABLE IF NOT EXISTS medical_center.specialties (
	specialty_id serial PRIMARY KEY,
	name varchar(150) NOT NULL UNIQUE,
	department_id integer NOT NULL,
	CONSTRAINT specialties_department_id_fkey FOREIGN KEY (department_id)
		REFERENCES medical_center.departments(department_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);
