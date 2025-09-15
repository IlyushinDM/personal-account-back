CREATE TABLE IF NOT EXISTS medical_center.appointmentstatuses (
	status_id serial PRIMARY KEY,
	name varchar(50) NOT NULL UNIQUE
);
