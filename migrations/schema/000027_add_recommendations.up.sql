ALTER TABLE medical_center.doctors
ADD COLUMN IF NOT EXISTS recommendations text;

ALTER TABLE medical_center.services
ADD COLUMN IF NOT EXISTS recommendations text;