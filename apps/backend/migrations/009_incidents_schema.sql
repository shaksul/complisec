-- Incidents module migration
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('technical_failure', 'data_breach', 'unauthorized_access', 'physical', 'malware', 'social_engineering', 'other')),
    status VARCHAR(20) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'assigned', 'in_progress', 'resolved', 'closed')),
    criticality VARCHAR(20) NOT NULL CHECK (criticality IN ('low', 'medium', 'high', 'critical')),
    source VARCHAR(50) NOT NULL CHECK (source IN ('user_report', 'automatic_agent', 'admin_manual', 'monitoring', 'siem')),
    reported_by UUID NOT NULL REFERENCES users(id),
    assigned_to UUID REFERENCES users(id),
    detected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP NULL,
    closed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create incident_assets table for many-to-many relationship
CREATE TABLE incident_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(incident_id, asset_id)
);

-- Create incident_risks table for many-to-many relationship
CREATE TABLE incident_risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(incident_id, risk_id)
);

-- Create incident_comments table for timeline
CREATE TABLE incident_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create incident_attachments table
CREATE TABLE incident_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create incident_actions table for corrective actions
CREATE TABLE incident_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL CHECK (action_type IN ('investigation', 'containment', 'eradication', 'recovery', 'prevention')),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assigned_to UUID REFERENCES users(id),
    due_date TIMESTAMP,
    completed_at TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create incident_metrics table for MTTR/MTTD tracking
CREATE TABLE incident_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL CHECK (metric_type IN ('mttd', 'mttr', 'resolution_time', 'response_time')),
    value_minutes INTEGER NOT NULL,
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_incidents_tenant_id ON incidents(tenant_id);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_criticality ON incidents(criticality);
CREATE INDEX idx_incidents_category ON incidents(category);
CREATE INDEX idx_incidents_assigned_to ON incidents(assigned_to);
CREATE INDEX idx_incidents_reported_by ON incidents(reported_by);
CREATE INDEX idx_incidents_detected_at ON incidents(detected_at);
CREATE INDEX idx_incidents_created_at ON incidents(created_at);

CREATE INDEX idx_incident_assets_incident_id ON incident_assets(incident_id);
CREATE INDEX idx_incident_assets_asset_id ON incident_assets(asset_id);
CREATE INDEX idx_incident_risks_incident_id ON incident_risks(incident_id);
CREATE INDEX idx_incident_risks_risk_id ON incident_risks(risk_id);
CREATE INDEX idx_incident_comments_incident_id ON incident_comments(incident_id);
CREATE INDEX idx_incident_comments_created_at ON incident_comments(created_at);
CREATE INDEX idx_incident_attachments_incident_id ON incident_attachments(incident_id);
CREATE INDEX idx_incident_actions_incident_id ON incident_actions(incident_id);
CREATE INDEX idx_incident_actions_assigned_to ON incident_actions(assigned_to);
CREATE INDEX idx_incident_actions_status ON incident_actions(status);
CREATE INDEX idx_incident_metrics_incident_id ON incident_metrics(incident_id);

-- Add permissions for incidents module
INSERT INTO permissions (code, module, description) 
SELECT 'incidents.view', 'incidents', 'View incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.view');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.create', 'incidents', 'Create incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.create');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.edit', 'incidents', 'Edit incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.edit');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.delete', 'incidents', 'Delete incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.delete');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.assign', 'incidents', 'Assign incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.assign');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.resolve', 'incidents', 'Resolve incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.resolve');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.close', 'incidents', 'Close incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.close');

INSERT INTO permissions (code, module, description) 
SELECT 'incidents.report', 'incidents', 'View incident reports'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.report');

-- Assign incidents permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions WHERE module = 'incidents';

-- Assign basic permissions to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000002', id FROM permissions WHERE code IN ('incidents.view', 'incidents.create');
