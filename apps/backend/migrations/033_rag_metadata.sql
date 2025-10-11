-- Минимальная схема - только метаданные индексации
-- Сам граф и векторы живут в GraphRAG + Qdrant

CREATE TABLE rag_indexed_documents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  document_id UUID NOT NULL REFERENCES documents(id),
  
  -- Статус индексации
  status VARCHAR(20) DEFAULT 'pending', -- pending, processing, indexed, failed, retrying
  error_message TEXT,
  retry_count INT DEFAULT 0,
  max_retries INT DEFAULT 3,
  
  -- Метаданные GraphRAG
  graphrag_doc_id TEXT, -- ID в GraphRAG (может быть NULL если GraphRAG не вернул)
  chunks_count INT DEFAULT 0,
  entities_count INT DEFAULT 0,
  relationships_count INT DEFAULT 0,
  
  -- Временные метки
  indexed_at TIMESTAMP,
  last_retry_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  
  UNIQUE(tenant_id, document_id)
);

CREATE INDEX idx_rag_indexed_tenant ON rag_indexed_documents(tenant_id);
CREATE INDEX idx_rag_indexed_status ON rag_indexed_documents(status);
CREATE INDEX idx_rag_indexed_document ON rag_indexed_documents(document_id);
CREATE INDEX idx_rag_indexed_retry ON rag_indexed_documents(status, retry_count) WHERE status = 'retrying';

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_rag_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_rag_updated_at
BEFORE UPDATE ON rag_indexed_documents
FOR EACH ROW
EXECUTE FUNCTION update_rag_updated_at();

-- Лог RAG-запросов для аудита
CREATE TABLE rag_query_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  user_id UUID REFERENCES users(id),
  
  query TEXT NOT NULL,
  use_graph BOOLEAN DEFAULT true,
  
  -- Результаты
  sources_count INT,
  response_time_ms INT,
  
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_rag_query_tenant ON rag_query_log(tenant_id);
CREATE INDEX idx_rag_query_created ON rag_query_log(created_at);

-- Добавить permissions для RAG
INSERT INTO permissions (code, module, description) VALUES
('rag.view', 'rag', 'View RAG management'),
('rag.index', 'rag', 'Index documents to RAG'),
('rag.query', 'rag', 'Query RAG system')
ON CONFLICT (code) DO NOTHING;

-- Назначить RAG permissions роли Admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id 
FROM permissions 
WHERE code IN ('rag.view', 'rag.index', 'rag.query')
ON CONFLICT DO NOTHING;

