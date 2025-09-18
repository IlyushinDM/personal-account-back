CREATE TABLE IF NOT EXISTS medical_center.prescriptions (
	id bigserial PRIMARY KEY,
	appointment_id bigint NOT NULL,
	doctor_id bigint NOT NULL,
	content text NOT NULL,
	status varchar(20) NOT NULL DEFAULT 'active',
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	completed_at timestamp without time zone,
	CONSTRAINT prescriptions_appointment_id_fkey FOREIGN KEY (appointment_id)
		REFERENCES medical_center.appointments(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT prescriptions_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_prescriptions_appointment_id ON medical_center.prescriptions(appointment_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_doctor_id ON medical_center.prescriptions(doctor_id);
