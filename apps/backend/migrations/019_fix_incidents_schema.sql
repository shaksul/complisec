-- Fix incidents table structure
-- Add missing columns to incidents table

-- Add category column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS category VARCHAR(50) NOT NULL DEFAULT 'other' CHECK (category IN ('technical_failure', 'data_breach', 'unauthorized_access', 'physical', 'malware', 'social_engineering', 'other'));

-- Add criticality column (rename from severity)
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS criticality VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (criticality IN ('low', 'medium', 'high', 'critical'));

-- Add source column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS source VARCHAR(50) NOT NULL DEFAULT 'user_report' CHECK (source IN ('user_report', 'automatic_agent', 'admin_manual', 'monitoring', 'siem'));

-- Add reported_by column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS reported_by UUID REFERENCES users(id);

-- Add detected_at column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS detected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Add resolved_at column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMP NULL;

-- Add closed_at column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS closed_at TIMESTAMP NULL;

-- Add deleted_at column
ALTER TABLE incidents ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Copy data from severity to criticality if severity exists
UPDATE incidents SET criticality = severity WHERE severity IS NOT NULL;

-- Drop old columns that are no longer needed
ALTER TABLE incidents DROP COLUMN IF EXISTS severity;
ALTER TABLE incidents DROP COLUMN IF EXISTS asset_id;
ALTER TABLE incidents DROP COLUMN IF EXISTS risk_id;
ALTER TABLE incidents DROP COLUMN IF EXISTS created_by;

-- Create incident_assets table for many-to-many relationship
CREATE TABLE IF NOT EXISTS incident_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(incident_id, asset_id)
);

-- Create incident_risks table for many-to-many relationship
CREATE TABLE IF NOT EXISTS incident_risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(incident_id, risk_id)
);

-- Create incident_comments table for timeline
CREATE TABLE IF NOT EXISTS incident_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create incident_attachments table
CREATE TABLE IF NOT EXISTS incident_attachments (
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
CREATE TABLE IF NOT EXISTS incident_actions (
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
CREATE TABLE IF NOT EXISTS incident_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    metric_type VARCHAR(50) NOT NULL CHECK (metric_type IN ('mttd', 'mttr', 'resolution_time', 'response_time')),
    value_minutes INTEGER NOT NULL,
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
CREATE INDEX IF NOT EXISTS idx_incidents_criticality ON incidents(criticality);
CREATE INDEX IF NOT EXISTS idx_incidents_category ON incidents(category);
CREATE INDEX IF NOT EXISTS idx_incidents_assigned_to ON incidents(assigned_to);
CREATE INDEX IF NOT EXISTS idx_incidents_reported_by ON incidents(reported_by);
CREATE INDEX IF NOT EXISTS idx_incidents_detected_at ON incidents(detected_at);
CREATE INDEX IF NOT EXISTS idx_incidents_created_at ON incidents(created_at);

CREATE INDEX IF NOT EXISTS idx_incident_assets_incident_id ON incident_assets(incident_id);
CREATE INDEX IF NOT EXISTS idx_incident_assets_asset_id ON incident_assets(asset_id);
CREATE INDEX IF NOT EXISTS idx_incident_risks_incident_id ON incident_risks(incident_id);
CREATE INDEX IF NOT EXISTS idx_incident_risks_risk_id ON incident_risks(risk_id);
CREATE INDEX IF NOT EXISTS idx_incident_comments_incident_id ON incident_comments(incident_id);
CREATE INDEX IF NOT EXISTS idx_incident_comments_created_at ON incident_comments(created_at);
CREATE INDEX IF NOT EXISTS idx_incident_attachments_incident_id ON incident_attachments(incident_id);
CREATE INDEX IF NOT EXISTS idx_incident_actions_incident_id ON incident_actions(incident_id);
CREATE INDEX IF NOT EXISTS idx_incident_actions_assigned_to ON incident_actions(assigned_to);
CREATE INDEX IF NOT EXISTS idx_incident_actions_status ON incident_actions(status);
CREATE INDEX IF NOT EXISTS idx_incident_metrics_incident_id ON incident_metrics(incident_id);













