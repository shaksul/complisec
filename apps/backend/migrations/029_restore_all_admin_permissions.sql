-- Migration: Restore all permissions for Admin role
-- Date: 2025-10-09
-- Description: Ensure Admin role has ALL permissions in the system

-- Grant ALL existing permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id AND rp.permission_id = p.id
  );




