-- Скрипт для проверки проблем с кодировкой в базе данных
-- Устанавливаем UTF-8 кодировку для сессии
SET client_encoding = 'UTF8';

-- Проверяем кодировку базы данных
SELECT 
    'Database encoding check:' as status,
    pg_encoding_to_char(encoding) as database_encoding,
    datname as database_name
FROM pg_database 
WHERE datname = 'complisec';

-- Проверяем кодировку клиента
SELECT 
    'Client encoding check:' as status,
    current_setting('client_encoding') as client_encoding;

-- Проверяем данные в таблице permissions на предмет искаженных символов
SELECT 
    'Permissions encoding check:' as status,
    code,
    description,
    CASE 
        WHEN description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]' THEN 'POSSIBLE_ENCODING_ISSUE'
        WHEN description ~ '[а-яё]' THEN 'CORRECT_UTF8'
        ELSE 'ASCII_ONLY'
    END as encoding_status
FROM permissions 
WHERE description IS NOT NULL
ORDER BY encoding_status DESC
LIMIT 20;

-- Проверяем данные в таблице roles
SELECT 
    'Roles encoding check:' as status,
    name,
    description,
    CASE 
        WHEN (name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') OR 
             (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]') THEN 'POSSIBLE_ENCODING_ISSUE'
        WHEN (name ~ '[а-яё]') OR (description ~ '[а-яё]') THEN 'CORRECT_UTF8'
        ELSE 'ASCII_ONLY'
    END as encoding_status
FROM roles 
WHERE name IS NOT NULL OR description IS NOT NULL
ORDER BY encoding_status DESC
LIMIT 10;

-- Проверяем данные в таблице users
SELECT 
    'Users encoding check:' as status,
    first_name,
    last_name,
    CASE 
        WHEN (first_name ~ '[^\x00-\x7F]' AND first_name !~ '[а-яё]') OR 
             (last_name ~ '[^\x00-\x7F]' AND last_name !~ '[а-яё]') THEN 'POSSIBLE_ENCODING_ISSUE'
        WHEN (first_name ~ '[а-яё]') OR (last_name ~ '[а-яё]') THEN 'CORRECT_UTF8'
        ELSE 'ASCII_ONLY'
    END as encoding_status
FROM users 
WHERE first_name IS NOT NULL OR last_name IS NOT NULL
ORDER BY encoding_status DESC
LIMIT 10;
