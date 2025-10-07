-- Enhanced Training Module Schema
-- Migration 013: Training Enhanced Schema

-- Update existing materials table to support all material types
ALTER TABLE materials ADD COLUMN IF NOT EXISTS material_type VARCHAR(50) NOT NULL DEFAULT 'document';
ALTER TABLE materials ADD COLUMN IF NOT EXISTS duration_minutes INTEGER;
ALTER TABLE materials ADD COLUMN IF NOT EXISTS tags TEXT[];
ALTER TABLE materials ADD COLUMN IF NOT EXISTS is_required BOOLEAN DEFAULT false;
ALTER TABLE materials ADD COLUMN IF NOT EXISTS passing_score INTEGER DEFAULT 80;
ALTER TABLE materials ADD COLUMN IF NOT EXISTS attempts_limit INTEGER;
ALTER TABLE materials ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Create training courses table (курсы как наборы материалов)
CREATE TABLE IF NOT EXISTS training_courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create course materials junction table
CREATE TABLE IF NOT EXISTS course_materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL REFERENCES training_courses(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    order_index INTEGER NOT NULL DEFAULT 0,
    is_required BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(course_id, material_id)
);

-- Enhanced training assignments table
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS course_id UUID REFERENCES training_courses(id);
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS assigned_by UUID REFERENCES users(id);
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'normal';
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS progress_percentage INTEGER DEFAULT 0;
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS time_spent_minutes INTEGER DEFAULT 0;
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS last_accessed_at TIMESTAMP;
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS reminder_sent_at TIMESTAMP;
ALTER TABLE train_assignments ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Create training progress tracking table
CREATE TABLE IF NOT EXISTS training_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignment_id UUID NOT NULL REFERENCES train_assignments(id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    progress_percentage INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0,
    last_position INTEGER,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assignment_id, material_id)
);

-- Enhanced quiz questions table
ALTER TABLE quiz_questions ADD COLUMN IF NOT EXISTS question_type VARCHAR(20) DEFAULT 'multiple_choice';
ALTER TABLE quiz_questions ADD COLUMN IF NOT EXISTS points INTEGER DEFAULT 1;
ALTER TABLE quiz_questions ADD COLUMN IF NOT EXISTS explanation TEXT;
ALTER TABLE quiz_questions ADD COLUMN IF NOT EXISTS order_index INTEGER DEFAULT 0;

-- Enhanced quiz attempts table
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS assignment_id UUID REFERENCES train_assignments(id);
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS max_score INTEGER;
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS time_spent_minutes INTEGER;
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS passed BOOLEAN DEFAULT false;

-- Create certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    assignment_id UUID NOT NULL REFERENCES train_assignments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    material_id UUID REFERENCES materials(id) ON DELETE CASCADE,
    course_id UUID REFERENCES training_courses(id) ON DELETE CASCADE,
    certificate_number VARCHAR(100) NOT NULL,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    is_valid BOOLEAN DEFAULT true,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create training notifications table
CREATE TABLE IF NOT EXISTS training_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    assignment_id UUID NOT NULL REFERENCES train_assignments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP
);

-- Create training analytics table
CREATE TABLE IF NOT EXISTS training_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    material_id UUID REFERENCES materials(id) ON DELETE CASCADE,
    course_id UUID REFERENCES training_courses(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL,
    metric_value DECIMAL(10,2) NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create role-based training assignments table
CREATE TABLE IF NOT EXISTS role_training_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    material_id UUID REFERENCES materials(id) ON DELETE CASCADE,
    course_id UUID REFERENCES training_courses(id) ON DELETE CASCADE,
    is_required BOOLEAN DEFAULT false,
    due_days INTEGER,
    assigned_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_materials_material_type ON materials(material_type);
CREATE INDEX IF NOT EXISTS idx_materials_is_required ON materials(is_required);
CREATE INDEX IF NOT EXISTS idx_training_courses_tenant_id ON training_courses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_course_materials_course_id ON course_materials(course_id);
CREATE INDEX IF NOT EXISTS idx_course_materials_material_id ON course_materials(material_id);
CREATE INDEX IF NOT EXISTS idx_train_assignments_course_id ON train_assignments(course_id);
CREATE INDEX IF NOT EXISTS idx_train_assignments_status ON train_assignments(status);
CREATE INDEX IF NOT EXISTS idx_train_assignments_due_at ON train_assignments(due_at);
CREATE INDEX IF NOT EXISTS idx_training_progress_assignment_id ON training_progress(assignment_id);
CREATE INDEX IF NOT EXISTS idx_training_progress_material_id ON training_progress(material_id);
CREATE INDEX IF NOT EXISTS idx_certificates_assignment_id ON certificates(assignment_id);
CREATE INDEX IF NOT EXISTS idx_certificates_user_id ON certificates(user_id);
CREATE INDEX IF NOT EXISTS idx_training_notifications_user_id ON training_notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_training_notifications_type ON training_notifications(type);
CREATE INDEX IF NOT EXISTS idx_training_analytics_tenant_id ON training_analytics(tenant_id);
CREATE INDEX IF NOT EXISTS idx_training_analytics_metric_type ON training_analytics(metric_type);
CREATE INDEX IF NOT EXISTS idx_role_training_assignments_role_id ON role_training_assignments(role_id);

-- Insert training-related permissions
INSERT INTO permissions (code, module, description) VALUES
('training.materials.view', 'training', 'View training materials'),
('training.materials.create', 'training', 'Create training materials'),
('training.materials.edit', 'training', 'Edit training materials'),
('training.materials.delete', 'training', 'Delete training materials'),
('training.courses.view', 'training', 'View training courses'),
('training.courses.create', 'training', 'Create training courses'),
('training.courses.edit', 'training', 'Edit training courses'),
('training.courses.delete', 'training', 'Delete training courses'),
('training.assign', 'training', 'Assign training materials'),
('training.progress.view', 'training', 'View training progress'),
('training.certificates.view', 'training', 'View certificates'),
('training.certificates.generate', 'training', 'Generate certificates'),
('training.analytics.view', 'training', 'View training analytics'),
('training.quizzes.create', 'training', 'Create quiz questions'),
('training.quizzes.edit', 'training', 'Edit quiz questions'),
('training.quizzes.delete', 'training', 'Delete quiz questions')
ON CONFLICT (code) DO NOTHING;
