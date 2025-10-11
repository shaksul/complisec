-- Добавление поля models для хранения списка доступных моделей
ALTER TABLE ai_providers ADD COLUMN IF NOT EXISTS models text[] DEFAULT ARRAY['llama3.2'];

-- Добавление поля default_model для модели по умолчанию
ALTER TABLE ai_providers ADD COLUMN IF NOT EXISTS default_model varchar(100) DEFAULT 'llama3.2';

-- Обновление существующих провайдеров
UPDATE ai_providers SET models = ARRAY['llama3.2', 'llama3.1', 'mistral', 'gemma2'] WHERE models IS NULL;
UPDATE ai_providers SET default_model = 'llama3.2' WHERE default_model IS NULL;




