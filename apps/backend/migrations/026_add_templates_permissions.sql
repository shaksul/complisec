-- Add permissions for templates and administration
-- Добавляем права для шаблонов и администрирования

-- Шаблоны документов
INSERT INTO permissions (code, module, description) VALUES
('templates.view', 'Templates', 'View templates'),
('templates.create', 'Templates', 'Create templates'),
('templates.edit', 'Templates', 'Edit templates'),
('templates.delete', 'Templates', 'Delete templates'),
('templates.manage', 'Templates', 'Manage templates')
ON CONFLICT (code) DO NOTHING;

-- Инвентарные номера
INSERT INTO permissions (code, module, description) VALUES
('inventory.view', 'Inventory', 'View inventory rules'),
('inventory.manage', 'Inventory', 'Manage inventory rules')
ON CONFLICT (code) DO NOTHING;

-- Организации
INSERT INTO permissions (code, module, description) VALUES
('organizations.view', 'Organizations', 'View organizations'),
('organizations.create', 'Organizations', 'Create organizations'),
('organizations.edit', 'Organizations', 'Edit organizations'),
('organizations.delete', 'Organizations', 'Delete organizations'),
('organizations.manage', 'Organizations', 'Manage organizations')
ON CONFLICT (code) DO NOTHING;

-- Системное администрирование
INSERT INTO permissions (code, module, description) VALUES
('admin', 'Administration', 'Full system access'),
('admin.settings', 'Administration', 'Manage system settings'),
('admin.system', 'Administration', 'System administration')
ON CONFLICT (code) DO NOTHING;

-- Назначаем новые права роли Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id 
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'Admin' 
AND p.code IN (
    'templates.view', 'templates.create', 'templates.edit', 'templates.delete', 'templates.manage',
    'inventory.view', 'inventory.manage',
    'organizations.view', 'organizations.create', 'organizations.edit', 'organizations.delete', 'organizations.manage',
    'admin', 'admin.settings', 'admin.system'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;