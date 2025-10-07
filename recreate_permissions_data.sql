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
('assets-create', 'assets.create', 'Создание активов', 'Создание активов', 'Активы'),
('assets-view', 'assets.view', 'Просмотр активов', 'Просмотр активов', 'Активы'),
('assets-edit', 'assets.edit', 'Редактирование активов', 'Редактирование активов', 'Активы'),
('assets-delete', 'assets.delete', 'Удаление активов', 'Удаление активов', 'Активы'),

-- Документы
('document-create', 'document.create', 'Создание документов', 'Создание документов', 'Документы'),
('document-read', 'document.read', 'Чтение документов', 'Чтение документов', 'Документы'),
('document-edit', 'document.edit', 'Редактирование документов', 'Редактирование документов', 'Документы'),
('document-delete', 'document.delete', 'Удаление документов', 'Удаление документов', 'Документы'),

-- Риски
('risks-create', 'risks.create', 'Создание рисков', 'Создание рисков', 'Риски'),
('risks-view', 'risks.view', 'Просмотр рисков', 'Просмотр рисков', 'Риски'),
('risks-edit', 'risks.edit', 'Редактирование рисков', 'Редактирование рисков', 'Риски'),
('risks-delete', 'risks.delete', 'Удаление рисков', 'Удаление рисков', 'Риски'),

-- Инциденты
('incidents-create', 'incidents.create', 'Создание инцидентов', 'Создание инцидентов', 'Инциденты'),
('incidents-view', 'incidents.view', 'Просмотр инцидентов', 'Просмотр инцидентов', 'Инциденты'),
('incidents-edit', 'incidents.edit', 'Редактирование инцидентов', 'Редактирование инцидентов', 'Инциденты'),
('incidents-delete', 'incidents.delete', 'Удаление инцидентов', 'Удаление инцидентов', 'Инциденты'),

-- Обучение
('training-create', 'training.create', 'Создание курсов обучения', 'Создание курсов обучения', 'Обучение'),
('training-view', 'training.view', 'Просмотр обучения', 'Просмотр обучения', 'Обучение'),
('training-edit', 'training.edit', 'Редактирование обучения', 'Редактирование обучения', 'Обучение'),
('training-delete', 'training.delete', 'Удаление обучения', 'Удаление обучения', 'Обучение'),

-- Пользователи
('users-view', 'users.view', 'Просмотр пользователей', 'Просмотр пользователей', 'Пользователи'),
('users-manage', 'users.manage', 'Управление пользователями', 'Управление пользователями', 'Пользователи'),

-- Роли
('roles-view', 'roles.view', 'Просмотр ролей', 'Просмотр ролей', 'Роли'),
('roles-manage', 'roles.manage', 'Управление ролями', 'Управление ролями', 'Роли'),

-- Организации
('organizations-view', 'organizations.view', 'Просмотр организаций', 'Просмотр организаций', 'Организации'),
('organizations-manage', 'organizations.manage', 'Управление организациями', 'Управление организациями', 'Организации');

-- Проверяем результат
SELECT 'CREATED PERMISSIONS:' as status, code, description, module 
FROM permissions 
ORDER BY module, code;
