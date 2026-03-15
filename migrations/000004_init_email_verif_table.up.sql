CREATE TABLE email_verification_codes (
    user_id     INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    code        VARCHAR(6) NOT NULL,
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);