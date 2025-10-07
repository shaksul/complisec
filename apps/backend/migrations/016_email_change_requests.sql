-- Создание таблицы для запросов на смену email
CREATE TABLE email_change_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    old_email VARCHAR(255) NOT NULL,
    new_email VARCHAR(255) NOT NULL,
    verification_code VARCHAR(6) NOT NULL,
    old_email_verified BOOLEAN DEFAULT FALSE,
    new_email_verified BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'old_email_verified', 'new_email_verified', 'completed', 'expired', 'cancelled'))
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_email_change_requests_user_id ON email_change_requests(user_id);
CREATE INDEX idx_email_change_requests_tenant_id ON email_change_requests(tenant_id);
CREATE INDEX idx_email_change_requests_verification_code ON email_change_requests(verification_code);
CREATE INDEX idx_email_change_requests_status ON email_change_requests(status);
CREATE INDEX idx_email_change_requests_expires_at ON email_change_requests(expires_at);

-- Уникальный индекс для предотвращения дублирования активных запросов
CREATE UNIQUE INDEX idx_email_change_requests_active_user ON email_change_requests(user_id, tenant_id) 
WHERE status IN ('pending', 'old_email_verified', 'new_email_verified');

-- Создание таблицы для аудит-лога изменений email
CREATE TABLE email_change_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    old_email VARCHAR(255) NOT NULL,
    new_email VARCHAR(255) NOT NULL,
    change_type VARCHAR(20) NOT NULL CHECK (change_type IN ('requested', 'old_email_verified', 'new_email_verified', 'completed', 'cancelled', 'expired')),
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для аудит-лога
CREATE INDEX idx_email_change_audit_log_user_id ON email_change_audit_log(user_id);
CREATE INDEX idx_email_change_audit_log_tenant_id ON email_change_audit_log(tenant_id);
CREATE INDEX idx_email_change_audit_log_created_at ON email_change_audit_log(created_at);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_email_change_requests_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER trigger_update_email_change_requests_updated_at
    BEFORE UPDATE ON email_change_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_email_change_requests_updated_at();
