-- Assets module migration
-- Drop existing assets table if it exists (it was created in 001_initial_schema.sql but needs to be updated)
DROP TABLE IF EXISTS assets CASCADE;

-- Create assets table with full specification
CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    inventory_number VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('server', 'workstation', 'application', 'database', 'document', 'network_device', 'other')),
    class VARCHAR(50) NOT NULL CHECK (class IN ('hardware', 'software', 'data', 'service')),
    owner_id UUID REFERENCES users(id),
    location TEXT,
    criticality VARCHAR(20) NOT NULL CHECK (criticality IN ('low', 'medium', 'high')),
    confidentiality VARCHAR(20) NOT NULL CHECK (confidentiality IN ('low', 'medium', 'high')),
    integrity VARCHAR(20) NOT NULL CHECK (integrity IN ('low', 'medium', 'high')),
    availability VARCHAR(20) NOT NULL CHECK (availability IN ('low', 'medium', 'high')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'in_repair', 'storage', 'decommissioned')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create asset_documents table
CREATE TABLE asset_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL CHECK (document_type IN ('passport', 'transfer_act', 'writeoff_act', 'repair_log', 'other')),
    file_path TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create asset_history table
CREATE TABLE asset_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    field_changed VARCHAR(100) NOT NULL,
    old_value TEXT,
    new_value TEXT NOT NULL,
    changed_by UUID NOT NULL REFERENCES users(id),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create asset_software table
CREATE TABLE asset_software (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    software_name VARCHAR(255) NOT NULL,
    version VARCHAR(100),
    installed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_assets_tenant_id ON assets(tenant_id);
CREATE INDEX idx_assets_owner_id ON assets(owner_id);
CREATE INDEX idx_assets_type ON assets(type);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_criticality ON assets(criticality);
CREATE INDEX idx_assets_inventory_number ON assets(inventory_number);
CREATE INDEX idx_asset_documents_asset_id ON asset_documents(asset_id);
CREATE INDEX idx_asset_history_asset_id ON asset_history(asset_id);
CREATE INDEX idx_asset_history_changed_at ON asset_history(changed_at);
CREATE INDEX idx_asset_software_asset_id ON asset_software(asset_id);

-- Add permissions for assets module (only if they don't exist)
INSERT INTO permissions (code, module, description) 
SELECT 'assets.view', 'assets', 'View assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.view');

INSERT INTO permissions (code, module, description) 
SELECT 'assets.create', 'assets', 'Create assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.create');

INSERT INTO permissions (code, module, description) 
SELECT 'assets.edit', 'assets', 'Edit assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.edit');

INSERT INTO permissions (code, module, description) 
SELECT 'assets.delete', 'assets', 'Delete assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.delete');

INSERT INTO permissions (code, module, description) 
SELECT 'assets.export', 'assets', 'Export assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.export');

INSERT INTO permissions (code, module, description) 
SELECT 'assets.inventory', 'assets', 'Perform inventory'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.inventory');

-- Assign assets permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions WHERE module = 'assets';
