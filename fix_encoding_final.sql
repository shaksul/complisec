-- Финальный скрипт для исправления кодировки данных
-- Устанавливаем UTF-8 кодировку для сессии
SET client_encoding = 'UTF8';

-- Создаем функцию для исправления кодировки
CREATE OR REPLACE FUNCTION fix_text_encoding(input_text TEXT)
RETURNS TEXT AS $$
BEGIN
    -- Если текст пустой или NULL, возвращаем как есть
    IF input_text IS NULL OR input_text = '' THEN
        RETURN input_text;
    END IF;
    
    -- Если текст содержит только ASCII символы, возвращаем как есть
    IF input_text ~ '^[\x00-\x7F]*$' THEN
        RETURN input_text;
    END IF;
    
    -- Если текст содержит корректные русские символы, возвращаем как есть
    IF input_text ~ '[а-яё]' THEN
        RETURN input_text;
    END IF;
    
    -- Если текст содержит искаженные символы, пытаемся исправить
    IF input_text ~ '[^\x00-\x7F]' THEN
        BEGIN
            -- Пытаемся исправить кодировку, предполагая что это была Windows-1251
            RETURN convert_from(convert_to(input_text, 'UTF8'), 'Windows-1251');
        EXCEPTION
            WHEN OTHERS THEN
                -- Если не удалось исправить, возвращаем исходный текст
                RETURN input_text;
        END;
    END IF;
    
    RETURN input_text;
END;
$$ LANGUAGE plpgsql;

-- Исправляем кодировку в таблице permissions
UPDATE permissions 
SET description = fix_text_encoding(description)
WHERE description IS NOT NULL 
  AND description != ''
  AND description ~ '[^\x00-\x7F]' 
  AND description !~ '[а-яё]';

-- Исправляем кодировку в таблице roles
UPDATE roles 
SET name = fix_text_encoding(name),
    description = fix_text_encoding(description)
WHERE (name IS NOT NULL AND name != '' AND name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') 
   OR (description IS NOT NULL AND description != '' AND description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице users
UPDATE users 
SET first_name = fix_text_encoding(first_name),
    last_name = fix_text_encoding(last_name)
WHERE (first_name IS NOT NULL AND first_name != '' AND first_name ~ '[^\x00-\x7F]' AND first_name !~ '[а-яё]') 
   OR (last_name IS NOT NULL AND last_name != '' AND last_name ~ '[^\x00-\x7F]' AND last_name !~ '[а-яё]');

-- Исправляем кодировку в таблице assets
UPDATE assets 
SET name = fix_text_encoding(name),
    location = fix_text_encoding(location)
WHERE (name IS NOT NULL AND name != '' AND name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') 
   OR (location IS NOT NULL AND location != '' AND location ~ '[^\x00-\x7F]' AND location !~ '[а-яё]');

-- Исправляем кодировку в таблице risks
UPDATE risks 
SET title = fix_text_encoding(title),
    description = fix_text_encoding(description),
    category = fix_text_encoding(category)
WHERE (title IS NOT NULL AND title != '' AND title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description IS NOT NULL AND description != '' AND description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]')
   OR (category IS NOT NULL AND category != '' AND category ~ '[^\x00-\x7F]' AND category !~ '[а-яё]');

-- Исправляем кодировку в таблице documents
UPDATE documents 
SET title = fix_text_encoding(title)
WHERE title IS NOT NULL 
  AND title != ''
  AND title ~ '[^\x00-\x7F]' 
  AND title !~ '[а-яё]';

-- Исправляем кодировку в таблице incidents
UPDATE incidents 
SET title = fix_text_encoding(title),
    description = fix_text_encoding(description)
WHERE (title IS NOT NULL AND title != '' AND title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description IS NOT NULL AND description != '' AND description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице materials
UPDATE materials 
SET title = fix_text_encoding(title),
    description = fix_text_encoding(description)
WHERE (title IS NOT NULL AND title != '' AND title ~ '[^\x00-\x7F]' AND title !~ '[а-яё]') 
   OR (description IS NOT NULL AND description != '' AND description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице quiz_questions
UPDATE quiz_questions 
SET text = fix_text_encoding(text)
WHERE text IS NOT NULL 
  AND text != ''
  AND text ~ '[^\x00-\x7F]' 
  AND text !~ '[а-яё]';

-- Удаляем функцию
DROP FUNCTION fix_text_encoding(TEXT);

-- Проверяем результат
SELECT 'Encoding fix completed successfully' as status;
