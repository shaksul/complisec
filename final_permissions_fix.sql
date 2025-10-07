-- Пересоздание данных прав с правильной кодировкой
-- Удаляем все существующие права и создаем заново

-- Удаляем все права
DELETE FROM role_permissions;
DELETE FROM permissions;

-- Создаем права заново с правильной кодировкой
INSERT INTO permissions (id, code, description, module) VALUES
-- AI модуль
(gen_random_uuid(), 'ai.providers.view', 'Просмотр AI-провайдеров', 'ИИ'),
(gen_random_uuid(), 'ai.query.view', 'Просмотр AI-аналитики', 'ИИ'),

-- Активы
(gen_random_uuid(), 'assets.create', 'Создание активов', 'Активы'),
(gen_random_uuid(), 'assets.view', 'Просмотр активов', 'Активы'),
(gen_random_uuid(), 'assets.edit', 'Редактирование активов', 'Активы'),
(gen_random_uuid(), 'assets.delete', 'Удаление активов', 'Активы'),

-- Документы
(gen_random_uuid(), 'document.create', 'Создание документов', 'Документы'),
(gen_random_uuid(), 'document.read', 'Чтение документов', 'Документы'),
(gen_random_uuid(), 'document.edit', 'Редактирование документов', 'Документы'),
(gen_random_uuid(), 'document.delete', 'Удаление документов', 'Документы'),

-- Риски
(gen_random_uuid(), 'risks.create', 'Создание рисков', 'Риски'),
(gen_random_uuid(), 'risks.view', 'Просмотр рисков', 'Риски'),
(gen_random_uuid(), 'risks.edit', 'Редактирование рисков', 'Риски'),
(gen_random_uuid(), 'risks.delete', 'Удаление рисков', 'Риски'),

-- Инциденты
(gen_random_uuid(), 'incidents.create', 'Создание инцидентов', 'Инциденты'),
(gen_random_uuid(), 'incidents.view', 'Просмотр инцидентов', 'Инциденты'),
(gen_random_uuid(), 'incidents.edit', 'Редактирование инцидентов', 'Инциденты'),
(gen_random_uuid(), 'incidents.delete', 'Удаление инцидентов', 'Инциденты'),

-- Обучение
(gen_random_uuid(), 'training.create', 'Создание курсов обучения', 'Обучение'),
(gen_random_uuid(), 'training.view', 'Просмотр обучения', 'Обучение'),
(gen_random_uuid(), 'training.edit', 'Редактирование обучения', 'Обучение'),
(gen_random_uuid(), 'training.delete', 'Удаление обучения', 'Обучение'),

-- Пользователи
(gen_random_uuid(), 'users.view', 'Просмотр пользователей', 'Пользователи'),
(gen_random_uuid(), 'users.manage', 'Управление пользователями', 'Пользователи'),

-- Роли
(gen_random_uuid(), 'roles.view', 'Просмотр ролей', 'Роли'),
(gen_random_uuid(), 'roles.manage', 'Управление ролями', 'Роли'),

-- Организации
(gen_random_uuid(), 'organizations.view', 'Просмотр организаций', 'Организации'),
(gen_random_uuid(), 'organizations.manage', 'Управление организациями', 'Организации');

-- Проверяем результат
SELECT 'CREATED PERMISSIONS:' as status, code, description, module 
FROM permissions 
WHERE code IN ('ai.providers.view', 'ai.query.view')
ORDER BY module, code;
