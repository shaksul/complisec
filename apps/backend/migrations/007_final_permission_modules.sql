-- Final fix for all permission modules
-- Окончательное исправление всех модулей прав

-- Update remaining modules
UPDATE permissions SET module = 'Отчеты' WHERE module = 'reports';
UPDATE permissions SET module = 'Общие' WHERE module IS NULL OR module = '';

-- Show final results
SELECT module, COUNT(*) as count FROM permissions GROUP BY module ORDER BY module;
