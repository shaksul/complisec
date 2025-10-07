-- Migration: Add user activity tracking tables
-- Created: 2025-01-07

-- Create user_activities table
CREATE TABLE IF NOT EXISTS user_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_tenant_id ON user_activities(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_activities_created_at ON user_activities(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_activities_action ON user_activities(action);

-- Create login_attempts table
CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255),
    ip_address INET NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    failure_reason VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for login_attempts
CREATE INDEX IF NOT EXISTS idx_login_attempts_user_id ON login_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_login_attempts_tenant_id ON login_attempts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_login_attempts_email ON login_attempts(email);
CREATE INDEX IF NOT EXISTS idx_login_attempts_created_at ON login_attempts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_login_attempts_ip_address ON login_attempts(ip_address);

-- Create user_sessions table for tracking active sessions
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create indexes for user_sessions
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_tenant_id ON user_sessions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_active ON user_sessions(is_active);

-- Insert some sample activity data for testing
INSERT INTO user_activities (user_id, tenant_id, action, description, ip_address, user_agent, metadata) 
SELECT 
    u.id,
    u.tenant_id,
    'login',
    'Пользователь вошел в систему',
    '192.168.1.100',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    '{"method": "POST", "endpoint": "/api/auth/login"}'
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

INSERT INTO user_activities (user_id, tenant_id, action, description, ip_address, user_agent, metadata) 
SELECT 
    u.id,
    u.tenant_id,
    'document_view',
    'Просмотр документа "Политика безопасности"',
    '192.168.1.100',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    '{"document_id": "doc-001", "document_title": "Политика безопасности"}'
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

INSERT INTO user_activities (user_id, tenant_id, action, description, ip_address, user_agent, metadata) 
SELECT 
    u.id,
    u.tenant_id,
    'user_management',
    'Просмотр списка пользователей',
    '192.168.1.100',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    '{"page": 1, "filters": {"role": "all"}}'
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

INSERT INTO user_activities (user_id, tenant_id, action, description, ip_address, user_agent, metadata) 
SELECT 
    u.id,
    u.tenant_id,
    'risk_assessment',
    'Создание новой оценки риска',
    '192.168.1.100',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    '{"risk_id": "risk-001", "risk_type": "cyber_security"}'
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

-- Insert sample login attempts
INSERT INTO login_attempts (user_id, tenant_id, email, ip_address, user_agent, success, failure_reason) 
SELECT 
    u.id,
    u.tenant_id,
    u.email,
    '192.168.1.100',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    TRUE,
    NULL
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

INSERT INTO login_attempts (user_id, tenant_id, email, ip_address, user_agent, success, failure_reason) 
SELECT 
    u.id,
    u.tenant_id,
    u.email,
    '192.168.1.101',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36',
    TRUE,
    NULL
FROM users u 
WHERE u.email = 'admin@demo.local'
LIMIT 1;

INSERT INTO login_attempts (user_id, tenant_id, email, ip_address, user_agent, success, failure_reason) 
VALUES 
    (NULL, (SELECT id FROM tenants LIMIT 1), 'hacker@evil.com', '192.168.1.200', 'curl/7.68.0', FALSE, 'invalid_credentials'),
    (NULL, (SELECT id FROM tenants LIMIT 1), 'test@test.com', '192.168.1.201', 'Python-requests/2.28.1', FALSE, 'user_not_found');

