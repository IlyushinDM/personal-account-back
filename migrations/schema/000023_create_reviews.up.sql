CREATE TABLE IF NOT EXISTS medical_center.reviews (
	id bigserial PRIMARY KEY,
	user_id bigint NOT NULL,
	doctor_id bigint NOT NULL,
	rating smallint NOT NULL,
	comment text,
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_moderated boolean NOT NULL DEFAULT false,
	CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id)
		REFERENCES medical_center.users(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION,
	CONSTRAINT reviews_doctor_id_fkey FOREIGN KEY (doctor_id)
		REFERENCES medical_center.doctors(id)
		ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON medical_center.reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_doctor_id ON medical_center.reviews(doctor_id);
