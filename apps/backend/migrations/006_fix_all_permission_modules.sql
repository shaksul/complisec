-- Fix all permission modules to use Russian names
-- Исправляем все модули прав на русские названия

-- Check what we have first
SELECT 'Before update:' as status, module, COUNT(*) as count FROM permissions GROUP BY module ORDER BY module;

-- Update all modules to Russian names
UPDATE permissions SET module = 'Документы' WHERE module IN ('docs', 'documents', 'document');
UPDATE permissions SET module = 'Активы' WHERE module IN ('assets', 'asset');
UPDATE permissions SET module = 'Риски' WHERE module IN ('risks', 'risk');
UPDATE permissions SET module = 'Инциденты' WHERE module IN ('incidents', 'incident');
UPDATE permissions SET module = 'Обучение' WHERE module IN ('training', 'trainings');
UPDATE permissions SET module = 'Соответствие' WHERE module IN ('compliance');
UPDATE permissions SET module = 'ИИ' WHERE module IN ('ai', 'artificial_intelligence');
UPDATE permissions SET module = 'Пользователи' WHERE module IN ('users', 'user');
UPDATE permissions SET module = 'Роли' WHERE module IN ('roles', 'role');
UPDATE permissions SET module = 'Аудит' WHERE module IN ('audit', 'audits');
UPDATE permissions SET module = 'Дашборд' WHERE module IN ('dashboard', 'dashboards');

-- Show results
SELECT 'After update:' as status, module, COUNT(*) as count FROM permissions GROUP BY module ORDER BY module;
