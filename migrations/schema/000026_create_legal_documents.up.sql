CREATE TABLE IF NOT EXISTS medical_center.legal_documents (
    id bigserial PRIMARY KEY,
    type varchar(100) NOT NULL,
    title varchar(255) NOT NULL,
    url varchar(512) NOT NULL,
    version varchar(20) NOT NULL,
    update_date varchar(20) NOT NULL
);