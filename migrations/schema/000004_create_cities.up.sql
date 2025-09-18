CREATE TABLE IF NOT EXISTS medical_center.cities (
	id serial PRIMARY KEY,
	name varchar(100) NOT NULL UNIQUE
);
