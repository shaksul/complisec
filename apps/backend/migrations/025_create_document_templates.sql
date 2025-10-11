-- Migration: Create document templates and inventory number rules tables
-- Description: Tables for managing document templates and automatic inventory number generation

-- Table for document templates
CREATE TABLE IF NOT EXISTS document_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_type VARCHAR(50) NOT NULL, -- 'passport_pc', 'passport_monitor', 'passport_device', 'transfer_act', 'writeoff_act', 'repair_log'
    content TEXT NOT NULL, -- HTML template with variables {{field_name}}
    is_system BOOLEAN DEFAULT false, -- System-provided template (read-only for users)
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT chk_template_type CHECK (
        template_type IN (
            'passport_pc', 
            'passport_monitor', 
            'passport_device', 
            'transfer_act', 
            'writeoff_act',
            'repair_log',
            'other'
        )
    )
);

-- Table for inventory number generation rules
CREATE TABLE IF NOT EXISTS inventory_number_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    asset_type VARCHAR(50) NOT NULL, -- Asset type from assets.type or custom classification
    asset_class VARCHAR(50), -- Optional: specific class within type
    pattern VARCHAR(255) NOT NULL, -- Pattern: 'РСП-{{year}}-{{sequence:0000}}'
    current_sequence INTEGER DEFAULT 0,
    prefix VARCHAR(50),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uq_inventory_rules_tenant_type UNIQUE(tenant_id, asset_type, asset_class)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_document_templates_tenant ON document_templates(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_document_templates_type ON document_templates(template_type) WHERE deleted_at IS NULL AND is_active = true;
CREATE INDEX IF NOT EXISTS idx_document_templates_system ON document_templates(is_system) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_rules_tenant ON inventory_number_rules(tenant_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_inventory_rules_type ON inventory_number_rules(asset_type) WHERE is_active = true;

-- Add comments
COMMENT ON TABLE document_templates IS 'Document templates for generating asset passports and other documents';
COMMENT ON COLUMN document_templates.template_type IS 'Type of template (passport_pc, passport_monitor, etc.)';
COMMENT ON COLUMN document_templates.content IS 'HTML content with placeholders like {{asset_name}}, {{serial_number}}';
COMMENT ON COLUMN document_templates.is_system IS 'System templates are pre-installed and can only be copied, not modified directly';

COMMENT ON TABLE inventory_number_rules IS 'Rules for automatic generation of inventory numbers';
COMMENT ON COLUMN inventory_number_rules.pattern IS 'Pattern with variables: {{type_code}}, {{year}}, {{month}}, {{sequence}}';
COMMENT ON COLUMN inventory_number_rules.current_sequence IS 'Current sequence number for auto-increment';

-- Insert default inventory number rules for the default tenant
-- This will be executed only if there's a default tenant
DO $$
DECLARE
    default_tenant_id UUID;
    admin_user_id UUID;
BEGIN
    -- Get the first tenant (assuming it's the default one)
    SELECT id INTO default_tenant_id FROM tenants LIMIT 1;
    
    -- Get the first admin user (any user will do for template creation)
    SELECT id INTO admin_user_id FROM users LIMIT 1;
    
    IF default_tenant_id IS NOT NULL AND admin_user_id IS NOT NULL THEN
        -- Insert default rules based on the asset classification from the document
        INSERT INTO inventory_number_rules (tenant_id, asset_type, pattern, prefix, description) VALUES
        (default_tenant_id, 'hardware', 'РСП-{{year}}-{{sequence:0000}}', 'РСП', 'Рабочая станция пользователя'),
        (default_tenant_id, 'server', 'СО-{{year}}-{{sequence:0000}}', 'СО', 'Серверное оборудование'),
        (default_tenant_id, 'network', 'СР-{{year}}-{{sequence:0000}}', 'СР', 'Сетевой ресурс'),
        (default_tenant_id, 'storage', 'СХД-{{year}}-{{sequence:0000}}', 'СХД', 'Система хранения данных'),
        (default_tenant_id, 'mobile', 'МУ-{{year}}-{{sequence:0000}}', 'МУ', 'Мобильное устройство')
        ON CONFLICT (tenant_id, asset_type, asset_class) DO NOTHING;
    END IF;
END $$;


