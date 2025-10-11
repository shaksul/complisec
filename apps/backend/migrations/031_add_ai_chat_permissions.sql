-- Добавление прав для AI чата
INSERT INTO permissions (code, module, description) VALUES
('ai.chat.use', 'ai', 'Использование AI чата')
ON CONFLICT (code) DO NOTHING;

-- Выдача прав администратору
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Admin' AND p.code = 'ai.chat.use'
ON CONFLICT (role_id, permission_id) DO NOTHING;

