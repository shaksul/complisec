-- Простой скрипт для исправления UTF-8 кодировки
-- Устанавливаем UTF-8 кодировку для сессии
SET client_encoding = 'UTF8';

-- Проверяем текущую кодировку
SELECT 
    'Current encoding check:' as status,
    current_setting('client_encoding') as client_encoding,
    pg_encoding_to_char(encoding) as database_encoding
FROM pg_database 
WHERE datname = 'complisec';

-- Проверяем данные в таблице permissions
SELECT 
    'Permissions check:' as status,
    code,
    description
FROM permissions 
WHERE description LIKE '%документ%' 
   OR description LIKE '%актив%'
   OR description LIKE '%риск%'
LIMIT 10;

-- Проверяем данные в таблице roles
SELECT 
    'Roles check:' as status,
    name,
    description
FROM roles 
WHERE name LIKE '%админ%' 
   OR name LIKE '%пользователь%'
   OR description LIKE '%роль%'
LIMIT 10;
