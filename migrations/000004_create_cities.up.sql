CREATE TABLE IF NOT EXISTS medical_center.cities (
	city_id serial PRIMARY KEY,
	name varchar(100) NOT NULL UNIQUE
);
