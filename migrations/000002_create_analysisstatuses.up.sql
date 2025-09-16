CREATE TABLE IF NOT EXISTS medical_center.analysisstatuses (
	status_id serial PRIMARY KEY,
	name varchar(50) NOT NULL UNIQUE
);
