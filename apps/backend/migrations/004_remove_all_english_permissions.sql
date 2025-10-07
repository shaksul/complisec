-- Migration 004: Remove all English permissions, keep only Russian ones
-- Удаляем все английские права, оставляем только русские

-- Удаляем связи ролей с английскими правами
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE module IN (
        'assets', 'docs', 'incidents', 'risks', 'training', 'users', 'roles', 'audit', 'reports'
    )
);

-- Удаляем все английские права
DELETE FROM permissions WHERE module IN (
    'assets', 'docs', 'incidents', 'risks', 'training', 'users', 'roles', 'audit', 'reports'
);
