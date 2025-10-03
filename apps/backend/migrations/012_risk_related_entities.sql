-- Migration 012: Risk related entities according to ASSETS_RISKS.md Sprint 2
-- Add risk_controls, risk_comments, risk_history, risk_attachments, risk_tags

-- Create risk_controls table (связь с контролями)
CREATE TABLE risk_controls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    control_id UUID NOT NULL, -- будет ссылка на таблицу controls когда она появится
    control_name VARCHAR(255) NOT NULL,
    control_type VARCHAR(50) NOT NULL CHECK (control_type IN ('preventive', 'detective', 'corrective')),
    implementation_status VARCHAR(20) NOT NULL DEFAULT 'planned' CHECK (implementation_status IN ('planned', 'in_progress', 'implemented', 'not_applicable')),
    effectiveness VARCHAR(20) CHECK (effectiveness IN ('high', 'medium', 'low')),
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create risk_comments table (комментарии пользователей)
CREATE TABLE risk_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT false, -- внутренний комментарий (не виден всем)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create risk_history table (лог изменений)
CREATE TABLE risk_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    field_changed VARCHAR(100) NOT NULL, -- название измененного поля
    old_value TEXT, -- старое значение
    new_value TEXT, -- новое значение
    change_reason TEXT, -- причина изменения
    changed_by UUID NOT NULL REFERENCES users(id),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create risk_attachments table (вложения)
CREATE TABLE risk_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_hash VARCHAR(64), -- для проверки целостности
    description TEXT,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create risk_tags table (теги)
CREATE TABLE risk_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    tag_name VARCHAR(100) NOT NULL,
    tag_color VARCHAR(7) DEFAULT '#007bff', -- hex цвет
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(risk_id, tag_name) -- уникальная комбинация риск-тег
);

-- Create indexes for better performance
CREATE INDEX idx_risk_controls_risk_id ON risk_controls(risk_id);
CREATE INDEX idx_risk_controls_control_type ON risk_controls(control_type);
CREATE INDEX idx_risk_controls_implementation_status ON risk_controls(implementation_status);

CREATE INDEX idx_risk_comments_risk_id ON risk_comments(risk_id);
CREATE INDEX idx_risk_comments_user_id ON risk_comments(user_id);
CREATE INDEX idx_risk_comments_created_at ON risk_comments(created_at);

CREATE INDEX idx_risk_history_risk_id ON risk_history(risk_id);
CREATE INDEX idx_risk_history_changed_by ON risk_history(changed_by);
CREATE INDEX idx_risk_history_changed_at ON risk_history(changed_at);

CREATE INDEX idx_risk_attachments_risk_id ON risk_attachments(risk_id);
CREATE INDEX idx_risk_attachments_uploaded_by ON risk_attachments(uploaded_by);
CREATE INDEX idx_risk_attachments_uploaded_at ON risk_attachments(uploaded_at);

CREATE INDEX idx_risk_tags_risk_id ON risk_tags(risk_id);
CREATE INDEX idx_risk_tags_tag_name ON risk_tags(tag_name);

-- Add foreign key constraints for better data integrity
ALTER TABLE risk_controls ADD CONSTRAINT fk_risk_controls_created_by FOREIGN KEY (created_by) REFERENCES users(id);
ALTER TABLE risk_tags ADD CONSTRAINT fk_risk_tags_created_by FOREIGN KEY (created_by) REFERENCES users(id);

