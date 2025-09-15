CREATE TABLE IF NOT EXISTS medical_center.userprofiles (
	profile_id bigserial PRIMARY KEY,
	user_id bigint NOT NULL,
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	patronymic varchar(100),
	birth_date date NOT NULL,
	gender varchar(10) NOT NULL,
	city_id integer NOT NULL,
	email varchar(255),
	avatar_url varchar(512),
	CONSTRAINT userprofiles_user_id_key UNIQUE (user_id),
	CONSTRAINT userprofiles_email_key UNIQUE (email),
	CONSTRAINT userprofiles_user_id_fkey FOREIGN KEY (user_id)
		REFERENCES medical_center.users(user_id)
		ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT userprofiles_city_id_fkey FOREIGN KEY (city_id)
		REFERENCES medical_center.cities(city_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_city_id ON medical_center.userprofiles(city_id);
CREATE INDEX IF NOT EXISTS userprofiles_user_id_key ON medical_center.userprofiles(user_id);
