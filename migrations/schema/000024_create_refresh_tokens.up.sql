CREATE TABLE IF NOT EXISTS medical_center.refresh_tokens (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL UNIQUE,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,

    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id)
        REFERENCES medical_center.users(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON medical_center.refresh_tokens(user_id);