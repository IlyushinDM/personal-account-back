CREATE TABLE IF NOT EXISTS medical_center.labanalyses (
	analysis_id bigserial PRIMARY KEY,
	user_id bigint NOT NULL,
	appointment_id bigint,
	name varchar(255) NOT NULL,
	assigned_date date NOT NULL,
	status_id integer NOT NULL,
	result_file_url varchar(512),
	clinic_id bigint,
	CONSTRAINT labanalyses_user_id_fkey FOREIGN KEY (user_id)
		REFERENCES medical_center.users(user_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT labanalyses_appointment_id_fkey FOREIGN KEY (appointment_id)
		REFERENCES medical_center.appointments(appointment_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT labanalyses_clinic_id_fkey FOREIGN KEY (clinic_id)
		REFERENCES medical_center.clinics(clinic_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT labanalyses_status_id_fkey FOREIGN KEY (status_id)
		REFERENCES medical_center.analysisstatuses(status_id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_lab_analyses_status ON medical_center.labanalyses(status_id);
CREATE INDEX IF NOT EXISTS idx_lab_analyses_user_id ON medical_center.labanalyses(user_id);
