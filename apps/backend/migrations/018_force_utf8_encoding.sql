-- Принудительная установка UTF-8 кодировки для базы данных
-- Эта миграция гарантирует, что все новые подключения будут использовать UTF-8

-- Устанавливаем UTF-8 кодировку для текущей сессии
SET client_encoding = 'UTF8';

-- Устанавливаем UTF-8 кодировку для базы данных по умолчанию
ALTER DATABASE complisec SET client_encoding = 'UTF8';

-- Проверяем текущую кодировку
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

-- Проверяем, что все таблицы с русским текстом используют правильную кодировку
SELECT 
    'Table encoding check:' as status,
    schemaname,
    tablename,
    pg_encoding_to_char(encoding) as table_encoding
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public' 
  AND c.relkind = 'r'
  AND tablename IN ('permissions', 'roles', 'users');
