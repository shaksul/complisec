-- Migration: Add missing permissions for navigation
-- Date: 2025-10-09

-- Add missing permissions
INSERT INTO permissions (code, module, description)
SELECT 'compliance.view', 'compliance', 'View compliance'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'compliance.view');

INSERT INTO permissions (code, module, description)
SELECT 'ai.providers.view', 'ai', 'View AI providers'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'ai.providers.view');

INSERT INTO permissions (code, module, description)
SELECT 'ai.query.view', 'ai', 'View AI analytics'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'ai.query.view');

INSERT INTO permissions (code, module, description)
SELECT 'organizations.manage', 'organizations', 'Manage organizations'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'organizations.manage');

INSERT INTO permissions (code, module, description)
SELECT 'admin', 'system', 'Full admin access'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'admin');

-- Grant all new permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
  AND p.code IN (
    'compliance.view',
    'ai.providers.view',
    'ai.query.view',
    'organizations.manage',
    'admin'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );




