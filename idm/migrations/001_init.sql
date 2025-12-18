CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE DEFAULT gen_random_uuid(),
    username VARCHAR(50),
    email VARCHAR(100),
    password_hash VARCHAR(255),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_public_id ON users(public_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX idx_refresh_tokens_user_id (user_id),
CREATE INDEX idx_refresh_tokens_token (token),
CREATE INDEX idx_refresh_tokens_expires_at (expires_at)