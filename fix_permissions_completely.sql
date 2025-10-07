-- Полное исправление данных permissions
-- Очищаем все связанные таблицы и пересоздаем данные

-- Удаляем все связанные данные
DELETE FROM role_permissions;
DELETE FROM permissions;

-- Перезапускаем sequence для id (если есть)
-- SELECT setval('permissions_id_seq', 1, false);

-- Вставляем правильные данные из миграции с правильной кодировкой
-- Документы
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'document.read', 'Документы', 'Чтение документов'),
(gen_random_uuid(), 'document.upload', 'Документы', 'Загрузка документов'),
(gen_random_uuid(), 'document.edit', 'Документы', 'Редактирование документов'),
(gen_random_uuid(), 'document.delete', 'Документы', 'Удаление документов'),
(gen_random_uuid(), 'document.approve', 'Документы', 'Утверждение документов'),
(gen_random_uuid(), 'document.publish', 'Документы', 'Публикация документов');

-- Активы
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'asset.view', 'Активы', 'Просмотр активов'),
(gen_random_uuid(), 'asset.create', 'Активы', 'Создание активов'),
(gen_random_uuid(), 'asset.edit', 'Активы', 'Редактирование активов'),
(gen_random_uuid(), 'asset.delete', 'Активы', 'Удаление активов'),
(gen_random_uuid(), 'asset.assign', 'Активы', 'Назначение активов');

-- Риски
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'risk.view', 'Риски', 'Просмотр рисков'),
(gen_random_uuid(), 'risk.create', 'Риски', 'Создание рисков'),
(gen_random_uuid(), 'risk.edit', 'Риски', 'Редактирование рисков'),
(gen_random_uuid(), 'risk.delete', 'Риски', 'Удаление рисков'),
(gen_random_uuid(), 'risk.assess', 'Риски', 'Оценка рисков'),
(gen_random_uuid(), 'risk.mitigate', 'Риски', 'Управление рисками');

-- Инциденты
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'incident.view', 'Инциденты', 'Просмотр инцидентов'),
(gen_random_uuid(), 'incident.create', 'Инциденты', 'Создание инцидентов'),
(gen_random_uuid(), 'incident.edit', 'Инциденты', 'Редактирование инцидентов'),
(gen_random_uuid(), 'incident.close', 'Инциденты', 'Закрытие инцидентов'),
(gen_random_uuid(), 'incident.assign', 'Инциденты', 'Назначение инцидентов');

-- Обучение
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'training.view', 'Обучение', 'Просмотр обучения'),
(gen_random_uuid(), 'training.assign', 'Обучение', 'Назначение обучения'),
(gen_random_uuid(), 'training.create', 'Обучение', 'Создание курсов'),
(gen_random_uuid(), 'training.edit', 'Обучение', 'Редактирование курсов'),
(gen_random_uuid(), 'training.view_progress', 'Обучение', 'Просмотр прогресса');

-- Соответствие
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'compliance.view', 'Соответствие', 'Просмотр соответствия'),
(gen_random_uuid(), 'compliance.manage', 'Соответствие', 'Управление соответствием'),
(gen_random_uuid(), 'compliance.audit', 'Соответствие', 'Проведение аудитов');

-- ИИ (это ключевые права, которые были повреждены)
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'ai.providers.view', 'ИИ', 'Просмотр провайдеров ИИ'),
(gen_random_uuid(), 'ai.providers.manage', 'ИИ', 'Управление провайдерами ИИ'),
(gen_random_uuid(), 'ai.queries.view', 'ИИ', 'Просмотр запросов ИИ'),
(gen_random_uuid(), 'ai.queries.create', 'ИИ', 'Создание запросов ИИ');

-- Пользователи
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'users.view', 'Пользователи', 'Просмотр пользователей'),
(gen_random_uuid(), 'users.create', 'Пользователи', 'Создание пользователей'),
(gen_random_uuid(), 'users.edit', 'Пользователи', 'Редактирование пользователей'),
(gen_random_uuid(), 'users.delete', 'Пользователи', 'Удаление пользователей'),
(gen_random_uuid(), 'users.manage', 'Пользователи', 'Управление пользователями');

-- Роли
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'roles.view', 'Роли', 'Просмотр ролей'),
(gen_random_uuid(), 'roles.create', 'Роли', 'Создание ролей'),
(gen_random_uuid(), 'roles.edit', 'Роли', 'Редактирование ролей'),
(gen_random_uuid(), 'roles.delete', 'Роли', 'Удаление ролей');

-- Аудит
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'audit.view', 'Аудит', 'Просмотр журнала аудита'),
(gen_random_uuid(), 'audit.export', 'Аудит', 'Экспорт журнала аудита');

-- Дашборд
INSERT INTO permissions (id, code, module, description) VALUES 
(gen_random_uuid(), 'dashboard.view', 'Дашборд', 'Просмотр дашборда'),
(gen_random_uuid(), 'dashboard.analytics', 'Дашборд', 'Просмотр аналитики');

-- Проверяем результат для AI прав
SELECT 'FIXED AI PERMISSIONS:' as status, code, description, module 
FROM permissions 
WHERE code LIKE 'ai.%'
ORDER BY code;
