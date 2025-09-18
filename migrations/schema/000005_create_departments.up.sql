CREATE TABLE IF NOT EXISTS medical_center.departments (
	id serial PRIMARY KEY,
	name varchar(100) NOT NULL UNIQUE
);
