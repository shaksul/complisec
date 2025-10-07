-- Documents Module Schema
-- Migration 014: Documents Schema

-- Create folders table (логические папки)
CREATE TABLE IF NOT EXISTS folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB
);

-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    folder_id UUID REFERENCES folders(id) ON DELETE SET NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    version INTEGER DEFAULT 1,
    metadata JSONB
);

-- Create document tags table
CREATE TABLE IF NOT EXISTS document_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(document_id, tag)
);

-- Create document links table (связи с другими модулями)
CREATE TABLE IF NOT EXISTS document_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    module VARCHAR(50) NOT NULL, -- 'risk', 'asset', 'incident', 'training', 'compliance'
    entity_id UUID NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(document_id, module, entity_id)
);

-- Create OCR text table
CREATE TABLE IF NOT EXISTS ocr_text (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    language VARCHAR(10) DEFAULT 'ru',
    confidence DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create document permissions table (расширение RBAC для документов)
CREATE TABLE IF NOT EXISTS document_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    subject_type VARCHAR(20) NOT NULL, -- 'user', 'role'
    subject_id UUID NOT NULL,
    object_type VARCHAR(20) NOT NULL, -- 'document', 'folder'
    object_id UUID NOT NULL,
    permission VARCHAR(50) NOT NULL, -- 'view', 'edit', 'delete', 'share'
    granted_by UUID NOT NULL REFERENCES users(id),
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(subject_type, subject_id, object_type, object_id, permission)
);

-- Create document versions table (версионирование)
CREATE TABLE IF NOT EXISTS document_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_description TEXT,
    UNIQUE(document_id, version_number)
);

-- Create document audit log table
CREATE TABLE IF NOT EXISTS document_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'created', 'updated', 'deleted', 'viewed', 'downloaded', 'shared'
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_folders_tenant_id ON folders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_folders_parent_id ON folders(parent_id);
CREATE INDEX IF NOT EXISTS idx_folders_owner_id ON folders(owner_id);
CREATE INDEX IF NOT EXISTS idx_folders_is_active ON folders(is_active);

CREATE INDEX IF NOT EXISTS idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_documents_folder_id ON documents(folder_id);
CREATE INDEX IF NOT EXISTS idx_documents_owner_id ON documents(owner_id);
CREATE INDEX IF NOT EXISTS idx_documents_created_by ON documents(created_by);
CREATE INDEX IF NOT EXISTS idx_documents_mime_type ON documents(mime_type);
CREATE INDEX IF NOT EXISTS idx_documents_file_hash ON documents(file_hash);
CREATE INDEX IF NOT EXISTS idx_documents_is_active ON documents(is_active);
CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents(created_at);

CREATE INDEX IF NOT EXISTS idx_document_tags_document_id ON document_tags(document_id);
CREATE INDEX IF NOT EXISTS idx_document_tags_tag ON document_tags(tag);

CREATE INDEX IF NOT EXISTS idx_document_links_document_id ON document_links(document_id);
CREATE INDEX IF NOT EXISTS idx_document_links_module ON document_links(module);
CREATE INDEX IF NOT EXISTS idx_document_links_entity_id ON document_links(entity_id);

CREATE INDEX IF NOT EXISTS idx_ocr_text_document_id ON ocr_text(document_id);
CREATE INDEX IF NOT EXISTS idx_ocr_text_language ON ocr_text(language);

CREATE INDEX IF NOT EXISTS idx_document_permissions_tenant_id ON document_permissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_document_permissions_subject ON document_permissions(subject_type, subject_id);
CREATE INDEX IF NOT EXISTS idx_document_permissions_object ON document_permissions(object_type, object_id);
CREATE INDEX IF NOT EXISTS idx_document_permissions_permission ON document_permissions(permission);
CREATE INDEX IF NOT EXISTS idx_document_permissions_is_active ON document_permissions(is_active);

CREATE INDEX IF NOT EXISTS idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX IF NOT EXISTS idx_document_versions_version_number ON document_versions(version_number);

CREATE INDEX IF NOT EXISTS idx_document_audit_log_tenant_id ON document_audit_log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_document_audit_log_document_id ON document_audit_log(document_id);
CREATE INDEX IF NOT EXISTS idx_document_audit_log_user_id ON document_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_document_audit_log_action ON document_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_document_audit_log_created_at ON document_audit_log(created_at);

-- Insert document-related permissions
INSERT INTO permissions (code, module, description) VALUES
('documents.view', 'documents', 'View documents'),
('documents.create', 'documents', 'Create documents'),
('documents.edit', 'documents', 'Edit documents'),
('documents.delete', 'documents', 'Delete documents'),
('documents.download', 'documents', 'Download documents'),
('documents.share', 'documents', 'Share documents'),
('documents.upload', 'documents', 'Upload documents'),
('folders.view', 'documents', 'View folders'),
('folders.create', 'documents', 'Create folders'),
('folders.edit', 'documents', 'Edit folders'),
('folders.delete', 'documents', 'Delete folders'),
('documents.ocr', 'documents', 'Perform OCR on documents'),
('documents.audit', 'documents', 'View document audit log'),
('documents.permissions', 'documents', 'Manage document permissions')
ON CONFLICT (code) DO NOTHING;

-- Create default root folder for each tenant
INSERT INTO folders (tenant_id, name, description, owner_id, created_by)
SELECT 
    t.id as tenant_id,
    'Root' as name,
    'Root folder for documents' as description,
    u.id as owner_id,
    u.id as created_by
FROM tenants t
CROSS JOIN (
    SELECT id FROM users WHERE tenant_id = t.id LIMIT 1
) u
WHERE NOT EXISTS (
    SELECT 1 FROM folders f WHERE f.tenant_id = t.id AND f.parent_id IS NULL
);

