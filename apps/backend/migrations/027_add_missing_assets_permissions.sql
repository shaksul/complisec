-- Migration: Add missing assets permissions
-- Date: 2025-10-09

-- Add missing document-related permissions
INSERT INTO permissions (code, module, description)
SELECT 'assets.documents:create', 'assets', 'Create/upload asset documents'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.documents:create');

INSERT INTO permissions (code, module, description)
SELECT 'assets.documents:link', 'assets', 'Link existing documents to assets'
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE code = 'assets.documents:link');

-- Grant all assets permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
  AND p.code IN (
    'assets.view',
    'assets.create',
    'assets.edit',
    'assets.delete',
    'assets.export',
    'assets.inventory',
    'assets.documents:create',
    'assets.documents:link'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- Grant all assets permissions to Manager role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Manager'
  AND p.code IN (
    'assets.view',
    'assets.create',
    'assets.edit',
    'assets.export',
    'assets.inventory',
    'assets.documents:create',
    'assets.documents:link'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

-- Grant view permissions to User role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'User'
  AND p.code IN (
    'assets.view'
  )
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );

