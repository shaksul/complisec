-- Исправление коррупции кодировки в базе данных
-- Этот скрипт исправляет поврежденные русские тексты в таблице permissions

-- Проверяем текущее состояние поврежденных данных
SELECT 'BEFORE FIX:' as status, code, description 
FROM permissions 
WHERE description LIKE '%РџСЂРѕСЃРјРѕС‚СЂ%' 
   OR description LIKE '%РџСЂРѕСЃРјРѕС‚СЂ%'
   OR description LIKE '%РІРѕР·РјРѕР¶РЅРѕСЃС‚%'
   OR description LIKE '%РЅР°СЃС‚СЂРѕР№Рє%';

-- Исправляем поврежденные описания прав для AI модуля
UPDATE permissions 
SET description = 'Просмотр AI-провайдеров'
WHERE code = 'ai.providers.view' 
  AND description LIKE '%РџСЂРѕСЃРјРѕС‚СЂ%';

UPDATE permissions 
SET description = 'Просмотр AI-аналитики'
WHERE code = 'ai.query.view' 
  AND description LIKE '%РџСЂРѕСЃРјРѕС‚СЂ%';

-- Исправляем другие возможные поврежденные описания
UPDATE permissions 
SET description = 'Создание активов'
WHERE description LIKE '%Создание%' 
  AND description LIKE '%Р°РєС‚РёРІ%';

UPDATE permissions 
SET description = 'Удаление активов'
WHERE description LIKE '%Удаление%' 
  AND description LIKE '%Р°РєС‚РёРІ%';

-- Проверяем результат исправления
SELECT 'AFTER FIX:' as status, code, description 
FROM permissions 
WHERE code IN ('ai.providers.view', 'ai.query.view');

-- Общая проверка на наличие других поврежденных данных
SELECT 'OTHER CORRUPTED DATA:' as status, code, description 
FROM permissions 
WHERE description ~ '[Р-Я]{2,}'
   OR description LIKE '%Рџ%'
   OR description LIKE '%СЂРѕ%'
   OR description LIKE '%РІРѕ%';

-- Проверяем все права с русскими описаниями для выявления других проблем
SELECT 'ALL RUSSIAN DESCRIPTIONS:' as status, code, module, description 
FROM permissions 
WHERE description ~ '[а-яё]'
ORDER BY module, code;
