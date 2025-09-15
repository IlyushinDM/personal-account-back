CREATE TABLE IF NOT EXISTS medical_center.users (
	user_id bigserial PRIMARY KEY,
	phone varchar(20) NOT NULL UNIQUE,
	password_hash varchar(255) NOT NULL,
	gosuslugi_id varchar(255) UNIQUE,
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	is_active boolean NOT NULL DEFAULT true
);
