-- Migration 016: Remove duplicate English permissions
-- Удаляем дублирующие английские права, оставляем только русские

-- Удаляем английские права, которые дублируют русские
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'users.view', 'users.create', 'users.edit', 'users.delete',
        'roles.view', 'roles.create', 'roles.edit', 'roles.delete',
        'assets.view', 'assets.create', 'assets.edit', 'assets.delete',
        'risks.view', 'risks.create', 'risks.edit', 'risks.delete',
        'docs.view', 'docs.create', 'docs.edit', 'docs.approve',
        'incidents.view', 'incidents.create', 'incidents.edit',
        'training.view', 'training.assign', 'training.pass_quiz',
        'reports.view', 'audit.view'
    )
);

-- Удаляем сами английские права
DELETE FROM permissions WHERE code IN (
    'users.view', 'users.create', 'users.edit', 'users.delete',
    'roles.view', 'roles.create', 'roles.edit', 'roles.delete',
    'assets.view', 'assets.create', 'assets.edit', 'assets.delete',
    'risks.view', 'risks.create', 'risks.edit', 'risks.delete',
    'docs.view', 'docs.create', 'docs.edit', 'docs.approve',
    'incidents.view', 'incidents.create', 'incidents.edit',
    'training.view', 'training.assign', 'training.pass_quiz',
    'reports.view', 'audit.view'
);

-- Удаляем дополнительные английские права из других миграций
DELETE FROM role_permissions WHERE permission_id IN (
    SELECT id FROM permissions WHERE code IN (
        'assets.export', 'assets.inventory'
    )
);

DELETE FROM permissions WHERE code IN (
    'assets.export', 'assets.inventory'
);
