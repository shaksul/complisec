-- Скрипт для исправления UTF-8 кодировки в базе данных
-- Этот скрипт исправляет данные, которые были вставлены с неправильной кодировкой

-- Устанавливаем UTF-8 кодировку для сессии
SET client_encoding = 'UTF8';

-- Создаем временную функцию для исправления кодировки
CREATE OR REPLACE FUNCTION fix_utf8_encoding(input_text TEXT)
RETURNS TEXT AS $$
BEGIN
    -- Если текст содержит искаженные символы, пытаемся исправить
    IF input_text ~ '[^\x00-\x7F]' AND input_text !~ '[а-яё]' THEN
        -- Пытаемся исправить кодировку, предполагая что это была Windows-1251
        RETURN convert_from(convert_to(input_text, 'UTF8'), 'Windows-1251');
    ELSE
        RETURN input_text;
    END IF;
EXCEPTION
    WHEN OTHERS THEN
        -- Если не удалось исправить, возвращаем исходный текст
        RETURN input_text;
END;
$$ LANGUAGE plpgsql;

-- Исправляем кодировку в таблице permissions
UPDATE permissions 
SET description = fix_utf8_encoding(description)
WHERE description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]';

-- Исправляем кодировку в таблице roles
UPDATE roles 
SET name = fix_utf8_encoding(name),
    description = fix_utf8_encoding(description)
WHERE (name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') 
   OR (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице users
UPDATE users 
SET first_name = fix_utf8_encoding(first_name),
    last_name = fix_utf8_encoding(last_name)
WHERE (first_name ~ '[^\x00-\x7F]' AND first_name !~ '[а-яё]') 
   OR (last_name ~ '[^\x00-\x7F]' AND last_name !~ '[а-яё]');

-- Исправляем кодировку в таблице assets
UPDATE assets 
SET name = fix_utf8_encoding(name),
    location = fix_utf8_encoding(location),
    software = fix_utf8_encoding(software)
WHERE (name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') 
   OR (location ~ '[^\x00-\x7F]' AND location !~ '[а-яё]')
   OR (software ~ '[^\x00-\x7F]' AND software !~ '[а-яё]');

-- Исправляем кодировку в таблице risks
UPDATE risks 
SET title = fix_utf8_encoding(title),
    description = fix_utf8_encoding(description),
    category = fix_utf8_encoding(category)
WHERE (title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]')
   OR (category ~ '[^\x00-\x7F]' AND category !~ '[а-яё]');

-- Исправляем кодировку в таблице documents
UPDATE documents 
SET title = fix_utf8_encoding(title)
WHERE title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]';

-- Исправляем кодировку в таблице incidents
UPDATE incidents 
SET title = fix_utf8_encoding(title),
    description = fix_utf8_encoding(description)
WHERE (title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице materials
UPDATE materials 
SET title = fix_utf8_encoding(title),
    description = fix_utf8_encoding(description)
WHERE (title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице quiz_questions
UPDATE quiz_questions 
SET text = fix_utf8_encoding(text)
WHERE text ~ '[^\x00-\x7F]' AND text !~ '[а-яё]';

-- Удаляем временную функцию
DROP FUNCTION fix_utf8_encoding(TEXT);

-- Проверяем результат
SELECT 'Encoding fix completed' as status;
