-- Insert permissions according to the documentation
-- Группировка по модулям: Документы, Активы, Инциденты, Обучение, ИИ, Пользователи, Роли

-- Документы
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'document.read', 'Документы', 'Чтение документов'),
(uuid_generate_v4(), 'document.upload', 'Документы', 'Загрузка документов'),
(uuid_generate_v4(), 'document.edit', 'Документы', 'Редактирование документов'),
(uuid_generate_v4(), 'document.delete', 'Документы', 'Удаление документов'),
(uuid_generate_v4(), 'document.approve', 'Документы', 'Утверждение документов'),
(uuid_generate_v4(), 'document.publish', 'Документы', 'Публикация документов');

-- Активы
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'asset.view', 'Активы', 'Просмотр активов'),
(uuid_generate_v4(), 'asset.create', 'Активы', 'Создание активов'),
(uuid_generate_v4(), 'asset.edit', 'Активы', 'Редактирование активов'),
(uuid_generate_v4(), 'asset.delete', 'Активы', 'Удаление активов'),
(uuid_generate_v4(), 'asset.assign', 'Активы', 'Назначение активов');

-- Риски
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'risk.view', 'Риски', 'Просмотр рисков'),
(uuid_generate_v4(), 'risk.create', 'Риски', 'Создание рисков'),
(uuid_generate_v4(), 'risk.edit', 'Риски', 'Редактирование рисков'),
(uuid_generate_v4(), 'risk.delete', 'Риски', 'Удаление рисков'),
(uuid_generate_v4(), 'risk.assess', 'Риски', 'Оценка рисков'),
(uuid_generate_v4(), 'risk.mitigate', 'Риски', 'Управление рисками');

-- Инциденты
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'incident.view', 'Инциденты', 'Просмотр инцидентов'),
(uuid_generate_v4(), 'incident.create', 'Инциденты', 'Создание инцидентов'),
(uuid_generate_v4(), 'incident.edit', 'Инциденты', 'Редактирование инцидентов'),
(uuid_generate_v4(), 'incident.close', 'Инциденты', 'Закрытие инцидентов'),
(uuid_generate_v4(), 'incident.assign', 'Инциденты', 'Назначение инцидентов');

-- Обучение
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'training.view', 'Обучение', 'Просмотр обучения'),
(uuid_generate_v4(), 'training.assign', 'Обучение', 'Назначение обучения'),
(uuid_generate_v4(), 'training.create', 'Обучение', 'Создание курсов'),
(uuid_generate_v4(), 'training.edit', 'Обучение', 'Редактирование курсов'),
(uuid_generate_v4(), 'training.view_progress', 'Обучение', 'Просмотр прогресса');

-- Соответствие
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'compliance.view', 'Соответствие', 'Просмотр соответствия'),
(uuid_generate_v4(), 'compliance.manage', 'Соответствие', 'Управление соответствием'),
(uuid_generate_v4(), 'compliance.audit', 'Соответствие', 'Проведение аудитов');

-- ИИ
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'ai.providers.view', 'ИИ', 'Просмотр провайдеров ИИ'),
(uuid_generate_v4(), 'ai.providers.manage', 'ИИ', 'Управление провайдерами ИИ'),
(uuid_generate_v4(), 'ai.queries.view', 'ИИ', 'Просмотр запросов ИИ'),
(uuid_generate_v4(), 'ai.queries.create', 'ИИ', 'Создание запросов ИИ');

-- Пользователи
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'users.view', 'Пользователи', 'Просмотр пользователей'),
(uuid_generate_v4(), 'users.create', 'Пользователи', 'Создание пользователей'),
(uuid_generate_v4(), 'users.edit', 'Пользователи', 'Редактирование пользователей'),
(uuid_generate_v4(), 'users.delete', 'Пользователи', 'Удаление пользователей'),
(uuid_generate_v4(), 'users.manage', 'Пользователи', 'Управление пользователями');

-- Роли
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'roles.view', 'Роли', 'Просмотр ролей'),
(uuid_generate_v4(), 'roles.create', 'Роли', 'Создание ролей'),
(uuid_generate_v4(), 'roles.edit', 'Роли', 'Редактирование ролей'),
(uuid_generate_v4(), 'roles.delete', 'Роли', 'Удаление ролей');

-- Аудит
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'audit.view', 'Аудит', 'Просмотр журнала аудита'),
(uuid_generate_v4(), 'audit.export', 'Аудит', 'Экспорт журнала аудита');

-- Дашборд
INSERT INTO permissions (id, code, module, description) VALUES 
(uuid_generate_v4(), 'dashboard.view', 'Дашборд', 'Просмотр дашборда'),
(uuid_generate_v4(), 'dashboard.analytics', 'Дашборд', 'Просмотр аналитики');
