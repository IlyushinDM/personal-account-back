-- Имя таблицы изменено на user_profiles
CREATE TABLE IF NOT EXISTS medical_center.user_profiles (
	id bigserial PRIMARY KEY,
	user_id bigint NOT NULL,
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	patronymic varchar(100),
	birth_date date NOT NULL,
	gender varchar(10) NOT NULL,
	city_id integer NOT NULL,
	email varchar(255),
	avatar_url varchar(512),
	CONSTRAINT user_profiles_user_id_key UNIQUE (user_id), -- Имя constraint'а обновлено
	CONSTRAINT user_profiles_email_key UNIQUE (email),     -- Имя constraint'а обновлено
	CONSTRAINT user_profiles_user_id_fkey FOREIGN KEY (user_id)
		REFERENCES medical_center.users(id)
		ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT user_profiles_city_id_fkey FOREIGN KEY (city_id)
		REFERENCES medical_center.cities(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_city_id ON medical_center.user_profiles(city_id);
-- Индекс userprofiles_user_id_key уже создан через UNIQUE constraint