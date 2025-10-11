-- Migration: Add missing permissions for navigation
-- Date: 2025-10-09
-- Description: Add alias permissions to match frontend navigation requirements

-- Add missing view permissions as aliases
INSERT INTO permissions (code, module, description)
SELECT 'users.view', 'users', 'View users'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'users.view');

INSERT INTO permissions (code, module, description)
SELECT 'roles.view', 'roles', 'View roles'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'roles.view');

INSERT INTO permissions (code, module, description)
SELECT 'training.view', 'training', 'View training'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'training.view');

INSERT INTO permissions (code, module, description)
SELECT 'incidents.view', 'incidents', 'View incidents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'incidents.view');

-- Grant all new permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
  AND p.code IN (
    'users.view',
    'roles.view',
    'training.view',
    'incidents.view',
    'asset.view',
    'risk.view',
    'document.read',
    'incident.view',
    'users.manage',
    'training.view_progress'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );




