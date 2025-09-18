CREATE TABLE IF NOT EXISTS medical_center.analysisstatuses (
	id serial PRIMARY KEY,
	name varchar(50) NOT NULL UNIQUE
);
