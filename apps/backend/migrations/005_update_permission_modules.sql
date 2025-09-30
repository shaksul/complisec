-- Update existing permissions to use Russian module names according to documentation
-- Обновляем существующие права для использования русских названий модулей

-- Assets -> Активы
UPDATE permissions SET module = 'Активы' WHERE module = 'assets';

-- Documents -> Документы  
UPDATE permissions SET module = 'Документы' WHERE module = 'documents';

-- Risks -> Риски
UPDATE permissions SET module = 'Риски' WHERE module = 'risks';

-- Incidents -> Инциденты
UPDATE permissions SET module = 'Инциденты' WHERE module = 'incidents';

-- Training -> Обучение
UPDATE permissions SET module = 'Обучение' WHERE module = 'training';

-- Compliance -> Соответствие
UPDATE permissions SET module = 'Соответствие' WHERE module = 'compliance';

-- AI -> ИИ
UPDATE permissions SET module = 'ИИ' WHERE module = 'ai';

-- Users -> Пользователи
UPDATE permissions SET module = 'Пользователи' WHERE module = 'users';

-- Roles -> Роли
UPDATE permissions SET module = 'Роли' WHERE module = 'roles';

-- Audit -> Аудит
UPDATE permissions SET module = 'Аудит' WHERE module = 'audit';

-- Dashboard -> Дашборд
UPDATE permissions SET module = 'Дашборд' WHERE module = 'dashboard';
