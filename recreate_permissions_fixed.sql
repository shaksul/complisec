-- Пересоздание данных прав с правильной кодировкой
-- Удаляем все существующие права и создаем заново

-- Удаляем все права
DELETE FROM role_permissions;
DELETE FROM permissions;

-- Создаем права заново с правильной кодировкой
INSERT INTO permissions (id, code, description, module) VALUES
-- AI модуль
('ai-providers-view', 'ai.providers.view', 'Просмотр AI-провайдеров', 'ИИ'),
('ai-query-view', 'ai.query.view', 'Просмотр AI-аналитики', 'ИИ'),

-- Активы
('assets-create', 'assets.create', 'Создание активов', 'Активы'),
('assets-view', 'assets.view', 'Просмотр активов', 'Активы'),
('assets-edit', 'assets.edit', 'Редактирование активов', 'Активы'),
('assets-delete', 'assets.delete', 'Удаление активов', 'Активы'),

-- Документы
('document-create', 'document.create', 'Создание документов', 'Документы'),
('document-read', 'document.read', 'Чтение документов', 'Документы'),
('document-edit', 'document.edit', 'Редактирование документов', 'Документы'),
('document-delete', 'document.delete', 'Удаление документов', 'Документы'),

-- Риски
('risks-create', 'risks.create', 'Создание рисков', 'Риски'),
('risks-view', 'risks.view', 'Просмотр рисков', 'Риски'),
('risks-edit', 'risks.edit', 'Редактирование рисков', 'Риски'),
('risks-delete', 'risks.delete', 'Удаление рисков', 'Риски'),

-- Инциденты
('incidents-create', 'incidents.create', 'Создание инцидентов', 'Инциденты'),
('incidents-view', 'incidents.view', 'Просмотр инцидентов', 'Инциденты'),
('incidents-edit', 'incidents.edit', 'Редактирование инцидентов', 'Инциденты'),
('incidents-delete', 'incidents.delete', 'Удаление инцидентов', 'Инциденты'),

-- Обучение
('training-create', 'training.create', 'Создание курсов обучения', 'Обучение'),
('training-view', 'training.view', 'Просмотр обучения', 'Обучение'),
('training-edit', 'training.edit', 'Редактирование обучения', 'Обучение'),
('training-delete', 'training.delete', 'Удаление обучения', 'Обучение'),

-- Пользователи
('users-view', 'users.view', 'Просмотр пользователей', 'Пользователи'),
('users-manage', 'users.manage', 'Управление пользователями', 'Пользователи'),

-- Роли
('roles-view', 'roles.view', 'Просмотр ролей', 'Роли'),
('roles-manage', 'roles.manage', 'Управление ролями', 'Роли'),

-- Организации
('organizations-view', 'organizations.view', 'Просмотр организаций', 'Организации'),
('organizations-manage', 'organizations.manage', 'Управление организациями', 'Организации');

-- Проверяем результат
SELECT 'CREATED PERMISSIONS:' as status, code, description, module 
FROM permissions 
ORDER BY module, code;
