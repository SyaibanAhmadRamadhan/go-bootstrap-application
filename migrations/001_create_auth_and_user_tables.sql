-- Migration: Create users table
-- Created: 2025-11-23

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    status ENUM('active', 'inactive', 'suspended') NOT NULL DEFAULT 'active',
    phone VARCHAR(20) NULL,
    gender ENUM('male', 'female', 'other') NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_role (role),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Migration: Create auth_tokens table
-- Created: 2025-11-23

CREATE TABLE IF NOT EXISTS auth_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    token_type ENUM('access', 'refresh') NOT NULL,
    status ENUM('active', 'revoked', 'expired') NOT NULL DEFAULT 'active',
    expires_at TIMESTAMP NOT NULL,
    refresh_token TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_token_type (token_type),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sample data for testing
INSERT INTO users (email, password_hash, name, role, status) VALUES
('admin@example.com', '$2a$10$YourHashedPasswordHere', 'Admin User', 'admin', 'active'),
('user@example.com', '$2a$10$YourHashedPasswordHere', 'Regular User', 'user', 'active');
