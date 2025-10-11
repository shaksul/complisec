-- Миграция для связывания существующих документов с централизованным хранилищем
-- Эта миграция создает связи между существующими документами в модулях и централизованным хранилищем

-- Создаем временную таблицу для хранения информации о миграции
CREATE TABLE IF NOT EXISTS document_migration_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_module VARCHAR(50) NOT NULL,
    original_entity_id UUID NOT NULL,
    original_document_id UUID NOT NULL,
    new_document_id UUID NOT NULL,
    migration_status VARCHAR(20) DEFAULT 'pending',
    migration_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_document_migration_log_original ON document_migration_log(original_module, original_entity_id);
CREATE INDEX IF NOT EXISTS idx_document_migration_log_new ON document_migration_log(new_document_id);

-- Функция для создания связей документов активов с централизованным хранилищем
CREATE OR REPLACE FUNCTION migrate_asset_documents_to_centralized()
RETURNS INTEGER AS $$
DECLARE
    asset_doc RECORD;
    new_doc_id UUID;
    migrated_count INTEGER := 0;
BEGIN
    -- Проходим по всем документам активов
    FOR asset_doc IN 
        SELECT 
            ad.id,
            ad.asset_id,
            ad.document_type,
            ad.file_path,
            ad.created_by,
            ad.created_at,
            a.tenant_id,
            a.name as asset_name
        FROM asset_documents ad
        JOIN assets a ON ad.asset_id = a.id
        WHERE ad.id NOT IN (
            SELECT original_document_id 
            FROM document_migration_log 
            WHERE original_module = 'assets'
        )
    LOOP
        -- Создаем новый документ в централизованном хранилище
        INSERT INTO documents (
            id,
            tenant_id,
            name,
            original_name,
            description,
            file_path,
            file_size,
            mime_type,
            file_hash,
            folder_id,
            owner_id,
            created_by,
            created_at,
            updated_at,
            is_active,
            version,
            metadata
        ) VALUES (
            uuid_generate_v4(),
            asset_doc.tenant_id,
            asset_doc.asset_name || ' - ' || asset_doc.document_type,
            asset_doc.asset_name || ' - ' || asset_doc.document_type,
            'Migrated from asset documents',
            asset_doc.file_path,
            0, -- Размер неизвестен
            'application/octet-stream', -- MIME type по умолчанию
            '', -- Hash пустой (файл не проверялся)
            NULL,
            asset_doc.created_by,
            asset_doc.created_by,
            asset_doc.created_at,
            asset_doc.created_at,
            true,
            1,
            ('{"migrated_from":"asset_documents","original_asset_id":"' || asset_doc.asset_id || '","original_document_id":"' || asset_doc.id || '","document_type":"' || asset_doc.document_type || '"}')::jsonb
        ) RETURNING id INTO new_doc_id;

        -- Создаем связь с активом
        INSERT INTO document_links (
            id,
            document_id,
            module,
            entity_id,
            created_by,
            created_at
        ) VALUES (
            uuid_generate_v4(),
            new_doc_id,
            'assets',
            asset_doc.asset_id,
            asset_doc.created_by,
            asset_doc.created_at
        );

        -- Добавляем теги
        INSERT INTO document_tags (document_id, tag) VALUES (new_doc_id, '#активы');
        INSERT INTO document_tags (document_id, tag) VALUES (new_doc_id, '#' || asset_doc.document_type);
        INSERT INTO document_tags (document_id, tag) VALUES (new_doc_id, '#мигрировано');

        -- Записываем в лог миграции
        INSERT INTO document_migration_log (
            original_module,
            original_entity_id,
            original_document_id,
            new_document_id,
            migration_status,
            migration_notes
        ) VALUES (
            'assets',
            asset_doc.asset_id,
            asset_doc.id,
            new_doc_id,
            'completed',
            'Successfully migrated asset document to centralized storage'
        );

        migrated_count := migrated_count + 1;
    END LOOP;

    RETURN migrated_count;
END;
$$ LANGUAGE plpgsql;

-- Функция для создания связей документов рисков с централизованным хранилищем
CREATE OR REPLACE FUNCTION migrate_risk_attachments_to_centralized()
RETURNS INTEGER AS $$
DECLARE
    risk_attachment RECORD;
    new_doc_id UUID;
    migrated_count INTEGER := 0;
BEGIN
    -- Проходим по всем вложениям рисков
    FOR risk_attachment IN 
        SELECT 
            ra.id,
            ra.risk_id,
            ra.file_name,
            ra.file_path,
            ra.file_size,
            ra.mime_type,
            ra.description,
            ra.uploaded_by,
            ra.uploaded_at,
            r.tenant_id
        FROM risk_attachments ra
        JOIN risks r ON ra.risk_id = r.id
        WHERE ra.id NOT IN (
            SELECT original_document_id 
            FROM document_migration_log 
            WHERE original_module = 'risks'
        )
    LOOP
        -- Создаем новый документ в централизованном хранилище
        INSERT INTO documents (
            id,
            tenant_id,
            name,
            original_name,
            description,
            file_path,
            file_size,
            mime_type,
            file_hash,
            folder_id,
            owner_id,
            created_by,
            created_at,
            updated_at,
            is_active,
            version,
            metadata
        ) VALUES (
            uuid_generate_v4(),
            risk_attachment.tenant_id,
            risk_attachment.file_name,
            risk_attachment.file_name,
            COALESCE(risk_attachment.description, 'Migrated from risk attachments'),
            risk_attachment.file_path,
            risk_attachment.file_size,
            risk_attachment.mime_type,
            COALESCE(risk_attachment.file_hash, ''), -- Используем существующий hash или пустую строку
            NULL,
            risk_attachment.uploaded_by,
            risk_attachment.uploaded_by,
            risk_attachment.uploaded_at,
            risk_attachment.uploaded_at,
            true,
            1,
            ('{"migrated_from":"risk_attachments","original_risk_id":"' || risk_attachment.risk_id || '","original_attachment_id":"' || risk_attachment.id || '"}')::jsonb
        ) RETURNING id INTO new_doc_id;

        -- Создаем связь с риском
        INSERT INTO document_links (
            id,
            document_id,
            module,
            entity_id,
            created_by,
            created_at
        ) VALUES (
            uuid_generate_v4(),
            new_doc_id,
            'risks',
            risk_attachment.risk_id,
            risk_attachment.uploaded_by,
            risk_attachment.uploaded_at
        );

        -- Добавляем теги
        INSERT INTO document_tags (document_id, tag) VALUES (new_doc_id, '#риски');
        INSERT INTO document_tags (document_id, tag) VALUES (new_doc_id, '#мигрировано');

        -- Записываем в лог миграции
        INSERT INTO document_migration_log (
            original_module,
            original_entity_id,
            original_document_id,
            new_document_id,
            migration_status,
            migration_notes
        ) VALUES (
            'risks',
            risk_attachment.risk_id,
            risk_attachment.id,
            new_doc_id,
            'completed',
            'Successfully migrated risk attachment to centralized storage'
        );

        migrated_count := migrated_count + 1;
    END LOOP;

    RETURN migrated_count;
END;
$$ LANGUAGE plpgsql;

-- Создаем папки для организации документов по модулям
INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, created_at, updated_at, is_active, metadata)
SELECT 
    uuid_generate_v4(),
    t.id,
    'Активы',
    'Документы активов',
    NULL,
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    NOW(),
    NOW(),
    true,
    '{"module": "assets", "created_by_migration": true}'::jsonb
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM folders 
    WHERE name = 'Активы' AND tenant_id = t.id
) AND EXISTS (
    SELECT 1 FROM users WHERE tenant_id = t.id
);

INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, created_at, updated_at, is_active, metadata)
SELECT 
    uuid_generate_v4(),
    t.id,
    'Риски',
    'Документы рисков',
    NULL,
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    NOW(),
    NOW(),
    true,
    '{"module": "risks", "created_by_migration": true}'::jsonb
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM folders 
    WHERE name = 'Риски' AND tenant_id = t.id
) AND EXISTS (
    SELECT 1 FROM users WHERE tenant_id = t.id
);

INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, created_at, updated_at, is_active, metadata)
SELECT 
    uuid_generate_v4(),
    t.id,
    'Инциденты',
    'Документы инцидентов',
    NULL,
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    NOW(),
    NOW(),
    true,
    '{"module": "incidents", "created_by_migration": true}'::jsonb
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM folders 
    WHERE name = 'Инциденты' AND tenant_id = t.id
) AND EXISTS (
    SELECT 1 FROM users WHERE tenant_id = t.id
);

INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, created_at, updated_at, is_active, metadata)
SELECT 
    uuid_generate_v4(),
    t.id,
    'Обучение',
    'Документы обучения',
    NULL,
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    NOW(),
    NOW(),
    true,
    '{"module": "training", "created_by_migration": true}'::jsonb
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM folders 
    WHERE name = 'Обучение' AND tenant_id = t.id
) AND EXISTS (
    SELECT 1 FROM users WHERE tenant_id = t.id
);

-- Создаем папку для мигрированных документов
INSERT INTO folders (id, tenant_id, name, description, parent_id, owner_id, created_by, created_at, updated_at, is_active, metadata)
SELECT 
    uuid_generate_v4(),
    t.id,
    'Мигрированные документы',
    'Документы, перенесенные из старых модулей',
    NULL,
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    (SELECT id FROM users WHERE tenant_id = t.id ORDER BY created_at LIMIT 1),
    NOW(),
    NOW(),
    true,
    '{"migration_folder": true, "created_by_migration": true}'::jsonb
FROM tenants t
WHERE NOT EXISTS (
    SELECT 1 FROM folders 
    WHERE name = 'Мигрированные документы' AND tenant_id = t.id
) AND EXISTS (
    SELECT 1 FROM users WHERE tenant_id = t.id
);

-- Выполняем миграцию документов активов
SELECT migrate_asset_documents_to_centralized() as migrated_asset_documents;

-- Выполняем миграцию вложений рисков
SELECT migrate_risk_attachments_to_centralized() as migrated_risk_attachments;

-- Обновляем статистику миграции
INSERT INTO document_migration_log (
    original_module,
    original_entity_id,
    original_document_id,
    new_document_id,
    migration_status,
    migration_notes
) VALUES (
    'system',
    '00000000-0000-0000-0000-000000000000',
    '00000000-0000-0000-0000-000000000000',
    '00000000-0000-0000-0000-000000000000',
    'completed',
    'Migration completed successfully at ' || NOW()
);

-- Создаем представление для мониторинга миграции
CREATE OR REPLACE VIEW document_migration_status AS
SELECT 
    original_module,
    COUNT(*) as total_documents,
    COUNT(CASE WHEN migration_status = 'completed' THEN 1 END) as migrated_documents,
    COUNT(CASE WHEN migration_status = 'pending' THEN 1 END) as pending_documents,
    COUNT(CASE WHEN migration_status = 'failed' THEN 1 END) as failed_documents,
    MIN(created_at) as migration_started,
    MAX(updated_at) as last_migration_activity
FROM document_migration_log
WHERE original_module != 'system'
GROUP BY original_module;

-- Добавляем комментарии
COMMENT ON TABLE document_migration_log IS 'Лог миграции документов из модулей в централизованное хранилище';
COMMENT ON VIEW document_migration_status IS 'Статистика миграции документов по модулям';
COMMENT ON FUNCTION migrate_asset_documents_to_centralized() IS 'Мигрирует документы активов в централизованное хранилище';
COMMENT ON FUNCTION migrate_risk_attachments_to_centralized() IS 'Мигрирует вложения рисков в централизованное хранилище';
