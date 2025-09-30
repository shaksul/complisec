-- Enhanced Documents Module Migration
-- Based on updated requirements from документы.md

-- Add missing fields to documents table
ALTER TABLE documents 
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS code VARCHAR(50),
ADD COLUMN IF NOT EXISTS category VARCHAR(100),
ADD COLUMN IF NOT EXISTS tags TEXT[],
ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id),
ADD COLUMN IF NOT EXISTS classification VARCHAR(20) DEFAULT 'Internal' CHECK (classification IN ('Public', 'Internal', 'Confidential')),
ADD COLUMN IF NOT EXISTS effective_from DATE,
ADD COLUMN IF NOT EXISTS review_period_months INTEGER DEFAULT 12,
ADD COLUMN IF NOT EXISTS asset_ids UUID[],
ADD COLUMN IF NOT EXISTS risk_ids UUID[],
ADD COLUMN IF NOT EXISTS control_ids UUID[],
ADD COLUMN IF NOT EXISTS storage_key TEXT,
ADD COLUMN IF NOT EXISTS mime_type VARCHAR(100),
ADD COLUMN IF NOT EXISTS size_bytes BIGINT,
ADD COLUMN IF NOT EXISTS checksum_sha256 VARCHAR(64),
ADD COLUMN IF NOT EXISTS ocr_text TEXT,
ADD COLUMN IF NOT EXISTS av_scan_status VARCHAR(20) DEFAULT 'pending' CHECK (av_scan_status IN ('pending', 'clean', 'infected', 'error')),
ADD COLUMN IF NOT EXISTS av_scan_result TEXT,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Update status values to match new workflow
ALTER TABLE documents 
DROP CONSTRAINT IF EXISTS documents_status_check;

ALTER TABLE documents 
ADD CONSTRAINT documents_status_check 
CHECK (status IN ('draft', 'in_review', 'approved', 'obsolete'));

-- Create document versions table
CREATE TABLE IF NOT EXISTS document_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    storage_key TEXT NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    checksum_sha256 VARCHAR(64),
    ocr_text TEXT,
    av_scan_status VARCHAR(20) DEFAULT 'pending' CHECK (av_scan_status IN ('pending', 'clean', 'infected', 'error')),
    av_scan_result TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Create approval workflows table
CREATE TABLE IF NOT EXISTS approval_workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    workflow_type VARCHAR(20) NOT NULL CHECK (workflow_type IN ('sequential', 'parallel')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'approved', 'rejected', 'cancelled')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    completed_at TIMESTAMPTZ
);

-- Create approval steps table
CREATE TABLE IF NOT EXISTS approval_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES approval_workflows(id) ON DELETE CASCADE,
    step_order INTEGER NOT NULL,
    approver_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'skipped')),
    comments TEXT,
    deadline TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Create quizzes table first (referenced by other tables)
CREATE TABLE IF NOT EXISTS quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    questions JSONB NOT NULL,
    passing_score INTEGER DEFAULT 80,
    time_limit_minutes INTEGER,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Create acknowledgment campaigns table
CREATE TABLE IF NOT EXISTS ack_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    audience_type VARCHAR(20) NOT NULL CHECK (audience_type IN ('all', 'role', 'department', 'custom')),
    audience_ids UUID[],
    deadline DATE,
    quiz_id UUID REFERENCES quizzes(id),
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'completed', 'cancelled')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    completed_at TIMESTAMPTZ
);

-- Create acknowledgment assignments table
CREATE TABLE IF NOT EXISTS ack_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES ack_campaigns(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'overdue')),
    quiz_score INTEGER,
    quiz_passed BOOLEAN DEFAULT false,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Create training materials table
CREATE TABLE IF NOT EXISTS training_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('document', 'video', 'presentation', 'other')),
    storage_key TEXT NOT NULL,
    mime_type VARCHAR(100),
    size_bytes BIGINT,
    checksum_sha256 VARCHAR(64),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Create training assignments table
CREATE TABLE IF NOT EXISTS training_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    material_id UUID NOT NULL REFERENCES training_materials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    deadline DATE,
    quiz_id UUID REFERENCES quizzes(id),
    quiz_passed BOOLEAN DEFAULT false,
    quiz_score INTEGER,
    status VARCHAR(20) DEFAULT 'assigned' CHECK (status IN ('assigned', 'in_progress', 'completed', 'overdue')),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);


-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);
CREATE INDEX IF NOT EXISTS idx_documents_type ON documents(type);
CREATE INDEX IF NOT EXISTS idx_documents_owner_id ON documents(owner_id);
CREATE INDEX IF NOT EXISTS idx_documents_code ON documents(code);
CREATE INDEX IF NOT EXISTS idx_documents_effective_from ON documents(effective_from);

CREATE INDEX IF NOT EXISTS idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX IF NOT EXISTS idx_document_versions_version ON document_versions(version_number);

CREATE INDEX IF NOT EXISTS idx_approval_workflows_document_id ON approval_workflows(document_id);
CREATE INDEX IF NOT EXISTS idx_approval_workflows_status ON approval_workflows(status);

CREATE INDEX IF NOT EXISTS idx_approval_steps_workflow_id ON approval_steps(workflow_id);
CREATE INDEX IF NOT EXISTS idx_approval_steps_approver_id ON approval_steps(approver_id);

CREATE INDEX IF NOT EXISTS idx_ack_campaigns_document_id ON ack_campaigns(document_id);
CREATE INDEX IF NOT EXISTS idx_ack_campaigns_status ON ack_campaigns(status);

CREATE INDEX IF NOT EXISTS idx_ack_assignments_campaign_id ON ack_assignments(campaign_id);
CREATE INDEX IF NOT EXISTS idx_ack_assignments_user_id ON ack_assignments(user_id);

CREATE INDEX IF NOT EXISTS idx_training_materials_tenant_id ON training_materials(tenant_id);
CREATE INDEX IF NOT EXISTS idx_training_materials_type ON training_materials(type);

CREATE INDEX IF NOT EXISTS idx_training_assignments_material_id ON training_assignments(material_id);
CREATE INDEX IF NOT EXISTS idx_training_assignments_user_id ON training_assignments(user_id);

-- Add full-text search indexes
CREATE INDEX IF NOT EXISTS idx_documents_fts ON documents USING gin(to_tsvector('russian', title || ' ' || COALESCE(description, '') || ' ' || COALESCE(ocr_text, '')));
CREATE INDEX IF NOT EXISTS idx_training_materials_fts ON training_materials USING gin(to_tsvector('russian', title || ' ' || COALESCE(description, '')));

-- Update existing documents to have default values
UPDATE documents SET 
    classification = 'Internal',
    review_period_months = 12,
    av_scan_status = 'clean'
WHERE classification IS NULL;

-- Add comments for documentation
COMMENT ON TABLE documents IS 'Main documents table with enhanced metadata and workflow support';
COMMENT ON TABLE document_versions IS 'Version history for documents with file storage and processing info';
COMMENT ON TABLE approval_workflows IS 'Document approval workflows (sequential or parallel)';
COMMENT ON TABLE approval_steps IS 'Individual approval steps within workflows';
COMMENT ON TABLE ack_campaigns IS 'Acknowledgment campaigns for document familiarization';
COMMENT ON TABLE ack_assignments IS 'Individual user assignments for acknowledgment campaigns';
COMMENT ON TABLE training_materials IS 'Training materials (documents, videos, presentations)';
COMMENT ON TABLE training_assignments IS 'User assignments for training materials';
COMMENT ON TABLE quizzes IS 'Quizzes for training and acknowledgment campaigns';
