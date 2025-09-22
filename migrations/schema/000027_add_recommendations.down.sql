ALTER TABLE medical_center.doctors
DROP COLUMN IF EXISTS recommendations;

ALTER TABLE medical_center.services
DROP COLUMN IF EXISTS recommendations;