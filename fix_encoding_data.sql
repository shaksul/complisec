-- Скрипт для исправления кодировки данных в базе
-- Устанавливаем UTF-8 кодировку для сессии
SET client_encoding = 'UTF8';

-- Создаем функцию для исправления кодировки
CREATE OR REPLACE FUNCTION fix_encoding(text)
RETURNS text AS $$
BEGIN
    -- Если текст содержит искаженные символы, пытаемся исправить
    IF $1 ~ '[^\x00-\x7F]' AND $1 !~ '[а-яё]' THEN
        -- Пытаемся исправить кодировку, предполагая что это была Windows-1251
        BEGIN
            RETURN convert_from(convert_to($1, 'UTF8'), 'Windows-1251');
        EXCEPTION
            WHEN OTHERS THEN
                RETURN $1;
        END;
    ELSE
        RETURN $1;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Исправляем кодировку в таблице permissions
UPDATE permissions 
SET description = fix_encoding(description)
WHERE description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]';

-- Исправляем кодировку в таблице roles
UPDATE roles 
SET name = fix_encoding(name),
    description = fix_encoding(description)
WHERE (name ~ '[^\x00-\x7F]' AND name !~ '[а-яё]') 
   OR (description ~ '[^\x00-\x7F]' AND description !~ '[а-яё]');

-- Исправляем кодировку в таблице users
UPDATE users 
SET first_name = fix_encoding(first_name),
    last_name = fix_encoding(last_name)
WHERE (first_name ~ '[^\x00-\x7F]' AND first_name !~ '[а-яё]') 
   OR (last_name ~ '[^\x00-\x7F]' AND last_name !~ '[а-яё]');

-- Удаляем функцию
DROP FUNCTION fix_encoding(text);

-- Проверяем результат
SELECT 'Encoding fix completed' as status;
