CREATE TABLE IF NOT EXISTS medical_center.doctors (
	id bigserial PRIMARY KEY,
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	patronymic varchar(100),
	specialty_id integer NOT NULL,
	experience_years smallint NOT NULL,
	rating numeric(3,2) NOT NULL DEFAULT 0.00,
	review_count integer NOT NULL DEFAULT 0,
	avatar_url varchar(512),
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT doctors_specialty_id_fkey FOREIGN KEY (specialty_id)
		REFERENCES medical_center.specialties(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_doctors_specialty_id ON medical_center.doctors(specialty_id);
