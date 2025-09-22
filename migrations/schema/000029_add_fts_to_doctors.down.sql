DROP INDEX IF EXISTS medical_center.idx_doctors_fts;
DROP TRIGGER IF EXISTS tsvector_update ON medical_center.doctors;
DROP FUNCTION IF EXISTS medical_center.update_doctor_fts_document();
ALTER TABLE medical_center.doctors
DROP COLUMN IF EXISTS fts_document;