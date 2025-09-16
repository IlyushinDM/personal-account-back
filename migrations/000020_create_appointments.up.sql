CREATE TABLE IF NOT EXISTS medical_center.appointments (
	appointment_id bigserial PRIMARY KEY,
	user_id bigint NOT NULL,
	doctor_id bigint NOT NULL,
	service_id bigint NOT NULL,
	clinic_id bigint NOT NULL,
	appointment_date date NOT NULL,
	appointment_time time without time zone NOT NULL,
	status_id integer NOT NULL,
	price_at_booking numeric(10,2) NOT NULL,
	is_dms boolean NOT NULL DEFAULT false,
	pre_visit_instructions text,
	diagnosis text,
	recommendations text,
	result_file_url varchar(512),
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT appointments_user_id_fkey FOREIGN KEY (user_id)
		REFERENCES medical_center.users(user_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT appointments_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(doctor_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT appointments_service_id_fkey FOREIGN KEY (service_id)
		REFERENCES medical_center.services(service_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT appointments_clinic_id_fkey FOREIGN KEY (clinic_id)
		REFERENCES medical_center.clinics(clinic_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT appointments_status_id_fkey FOREIGN KEY (status_id)
		REFERENCES medical_center.appointmentstatuses(status_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON medical_center.appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_appointments_user_id ON medical_center.appointments(user_id);
