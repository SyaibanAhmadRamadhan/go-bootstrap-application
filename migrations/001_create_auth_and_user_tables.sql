-- Migration: Create users table
-- Created: 2025-11-23

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    phone VARCHAR(20) NULL,
    gender VARCHAR(20) NULL CHECK (gender IN ('male', 'female', 'other')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Migration: Create auth_tokens table
-- Created: 2025-11-23

CREATE TABLE IF NOT EXISTS auth_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    token_type VARCHAR(20) NOT NULL CHECK (token_type IN ('access', 'refresh')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked', 'expired')),
    expires_at TIMESTAMP NOT NULL,
    refresh_token TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX idx_auth_tokens_token_type ON auth_tokens(token_type);
CREATE INDEX idx_auth_tokens_status ON auth_tokens(status);
CREATE INDEX idx_auth_tokens_expires_at ON auth_tokens(expires_at);
CREATE INDEX idx_auth_tokens_created_at ON auth_tokens(created_at);

-- Sample data for testing
INSERT INTO users (email, password_hash, name, role, status) VALUES
('admin@example.com', '$2a$10$YourHashedPasswordHere', 'Admin User', 'admin', 'active'),
('user@example.com', '$2a$10$YourHashedPasswordHere', 'Regular User', 'user', 'active');
