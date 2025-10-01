RAG-системы лучше всего работают не с «каким-то одним универсальным форматом», а с нормализованными данными, которые можно эффективно индексировать, разбивать на чанки и преобразовывать в эмбеддинги. Формат выбирают под задачу и тип источников.

🔑 Основные принципы выбора формата

Структурированность
Чем более структурирован текст (JSON, Markdown, таблицы), тем проще выделять смысловые блоки.
Сырые форматы (PDF, DOCX, сканы) требуют OCR/парсинга.

Текстовая репрезентация
Итоговые данные должны быть в чистом тексте (UTF-8) для эмбеддингов. Даже если изначально это DOCX или PDF, в пайплайне их приводят к plain text или HTML-подобной структуре.

Семантическая сегментация
Важно не только хранить, но и правильно резать на чанки (например, 500–1000 токенов). Тут помогает Markdown/JSON, где есть явные заголовки, списки, поля.

📂 Форматы, которые чаще всего применяются

JSON – идеально для хранения знаний в виде записей (Q&A, документы, требования, политики). Удобен для последующей фильтрации по метаданным.

Markdown / HTML – хороши для текстов с иерархией, заголовками, списками. RAG-система понимает структуру и легче извлекает контекст.

Plain text (TXT, UTF-8) – базовый вариант после очистки. Универсален, но без разметки.

Parquet / CSV – когда источники табличные (реестры активов, базы инцидентов).

PDF / DOCX – допустимы как исходники, но перед подачей в RAG их нужно распарсить (например, pdfplumber, docx2txt) и сохранить в текст или JSON.

🧠 Пример пайплайна для RAG

Собрать документы в исходных форматах (PDF, DOCX, XLSX).

Пропустить через конвертеры → получить чистый текст.

Преобразовать в JSON с полями:

{
  "id": "doc-123",
  "title": "Политика безопасности",
  "section": "Контроль доступа",
  "content": "Описание требований и процедур...",
  "tags": ["ISO27001", "AccessControl"]
}


Нарезать на чанки (например, по 800 символов с overlap 200).

Сохранить в векторное хранилище (Pinecone, Weaviate, Qdrant).

🚀 Вывод

Лучший рабочий формат для RAG — JSON или Markdown после нормализации, потому что:

легко хранить текст + метаданные,

просто чанкать,

удобно фильтровать при запросе.

📌 Но в реальных проектах RAG всё равно работает с чем угодно (PDF, DOCX, базы), если есть правильный конвертер → текст → JSON.

Хотите, я соберу для вашей системы унифицированный шаблон JSON для всех документов (политики, риски, инциденты, активы), чтобы сразу можно было строить RAG-поиск?

Вы сказали:
да
ChatGPT сказал:
Думал

Отвечу как архитектор данных и эксперт по RAG-системам мирового уровня.

Пример ответа (структура первого сообщения):

Коротко: что даю. 2) Как этим пользоваться. 3) Шаблоны JSON Schema (универсальный + специализированные). 4) Примеры валидных документов. 5) Схема чанков и метаданных для индекса. 6) Правила нарезки/нормализации. 7) Мини-чек-лист интеграции в ваш пайплайн.

1) Коротко

Ниже — единый контракт данных для RAG CompliSec/RiskNexus:

Универсальный JSON Schema для всех типов знаний.

Специализированные схемы для: policy, risk, incident, asset, control, training_material, audit_record.

Стандартизированная модель чанков для индекса/эмбеддингов.

Жёсткие правила нормализации, версии, меток соответствия (ISO 27001 и др.).

Все схемы без плейсхолдеров, готовы «как есть».

2) Как пользоваться (по шагам)

При инжесте любого источника (PDF/DOCX/HTML/CSV) приводите к универсальному JSON (ниже).

Проверяйте документ валидатором JSON Schema.

Режьте по правилам из раздела 6 → формируйте chunks.jsonl по схеме из раздела 5.

Пишите в векторное хранилище: embedding, metadata из чанка.

В рантайме фильтруйте/ранжируйте по doc_type, tenant_id, tags, compliance_mappings.*, language, confidentiality.

Возвращайте пользователю оригинал и релевантные чанки.

3) Универсальная схема документа (JSON Schema Draft 2020-12)
{
  "$id": "https://complisec.local/schemas/base_document.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "CompliSec Base Document",
  "type": "object",
  "required": [
    "id",
    "tenant_id",
    "doc_type",
    "title",
    "language",
    "content",
    "source",
    "version",
    "created_at",
    "updated_at"
  ],
  "properties": {
    "id": {
      "type": "string",
      "description": "UUID v4 в нижнем регистре",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "tenant_id": {
      "type": "string",
      "description": "UUID v4 арендатора",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "doc_type": {
      "type": "string",
      "enum": [
        "policy",
        "risk",
        "incident",
        "asset",
        "control",
        "training_material",
        "audit_record",
        "other"
      ]
    },
    "title": { "type": "string", "minLength": 1 },
    "summary": { "type": "string" },
    "language": {
      "type": "string",
      "description": "BCP-47",
      "pattern": "^[A-Za-z]{2}(-[A-Za-z0-9-]+)*$"
    },
    "tags": {
      "type": "array",
      "items": { "type": "string" },
      "uniqueItems": true
    },
    "source": {
      "type": "object",
      "required": ["ingest_method", "origin_name"],
      "properties": {
        "ingest_method": { "type": "string", "enum": ["upload", "crawler", "api", "manual"] },
        "origin_name": { "type": "string" },
        "origin_uri": { "type": "string" },
        "file_mime": { "type": "string" },
        "checksum_sha256": { "type": "string", "pattern": "^[0-9a-f]{64}$" }
      },
      "additionalProperties": false
    },
    "version": {
      "type": "object",
      "required": ["number", "status"],
      "properties": {
        "number": { "type": "string", "description": "семантическая версия, например 1.0.0" },
        "status": { "type": "string", "enum": ["draft", "approved", "deprecated"] },
        "valid_from": { "type": "string", "format": "date-time" },
        "valid_to": { "type": "string", "format": "date-time" }
      },
      "additionalProperties": false
    },
    "classification": {
      "type": "object",
      "properties": {
        "confidentiality": { "type": "string", "enum": ["public", "internal", "restricted", "confidential"] },
        "integrity": { "type": "string", "enum": ["low", "medium", "high"] },
        "availability": { "type": "string", "enum": ["low", "medium", "high"] }
      },
      "additionalProperties": false
    },
    "compliance_mappings": {
      "type": "array",
      "description": "Привязки к стандартам соответствия",
      "items": {
        "type": "object",
        "required": ["standard", "clause"],
        "properties": {
          "standard": { "type": "string", "enum": ["ISO27001", "ISO27002", "NIST800-53", "SOC2", "GDPR", "ST-RK-27002"] },
          "clause": { "type": "string" },
          "level": { "type": "string", "enum": ["must", "should", "informative"] }
        },
        "additionalProperties": false
      },
      "uniqueItems": false
    },
    "relations": {
      "type": "object",
      "properties": {
        "related_ids": {
          "type": "array",
          "items": {
            "type": "string",
            "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
          },
          "uniqueItems": true
        }
      },
      "additionalProperties": false
    },
    "owners": {
      "type": "object",
      "properties": {
        "author": { "type": "string" },
        "reviewers": { "type": "array", "items": { "type": "string" }, "uniqueItems": true },
        "approver": { "type": "string" }
      },
      "additionalProperties": false
    },
    "content": {
      "type": "object",
      "required": ["format", "body"],
      "properties": {
        "format": { "type": "string", "enum": ["markdown", "plain", "html", "json"] },
        "body": { "type": "string" },
        "sections": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["title", "text"],
            "properties": {
              "number": { "type": "string" },
              "title": { "type": "string" },
              "text": { "type": "string" },
              "anchor": { "type": "string" }
            },
            "additionalProperties": false
          }
        }
      },
      "additionalProperties": false
    },
    "custom": {
      "type": "object",
      "description": "Специализированные поля для конкретного doc_type",
      "additionalProperties": true
    },
    "created_at": { "type": "string", "format": "date-time" },
    "updated_at": { "type": "string", "format": "date-time" },
    "deleted_at": { "type": "string", "format": "date-time" }
  },
  "additionalProperties": false
}

Специализированные под-схемы (вставляются в custom)

3.1 policy

{
  "$id": "https://complisec.local/schemas/policy.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Policy Custom Fields",
  "type": "object",
  "required": ["policy_id", "scope"],
  "properties": {
    "policy_id": { "type": "string" },
    "scope": { "type": "string" },
    "owner_department": { "type": "string" },
    "review_cycle_days": { "type": "integer", "minimum": 0 }
  },
  "additionalProperties": false
}


3.2 risk

{
  "$id": "https://complisec.local/schemas/risk.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Risk Custom Fields",
  "type": "object",
  "required": ["risk_id", "asset_id", "likelihood", "impact", "score_method", "inherent_score", "residual_score"],
  "properties": {
    "risk_id": { "type": "string" },
    "asset_id": { "type": "string" },
    "category": { "type": "string" },
    "likelihood": { "type": "string", "enum": ["very_low", "low", "medium", "high", "very_high"] },
    "impact": { "type": "string", "enum": ["very_low", "low", "medium", "high", "very_high"] },
    "score_method": { "type": "string", "enum": ["matrix_5x5", "cvss3", "cvss4"] },
    "inherent_score": { "type": "number", "minimum": 0 },
    "residual_score": { "type": "number", "minimum": 0 },
    "controls": {
      "type": "array",
      "items": { "type": "string" },
      "uniqueItems": true
    },
    "status": { "type": "string", "enum": ["identified", "assessed", "treated", "accepted", "closed"] }
  },
  "additionalProperties": false
}


3.3 incident

{
  "$id": "https://complisec.local/schemas/incident.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Incident Custom Fields",
  "type": "object",
  "required": ["incident_id", "severity", "detected_at", "status"],
  "properties": {
    "incident_id": { "type": "string" },
    "severity": { "type": "string", "enum": ["low", "medium", "high", "critical"] },
    "detected_at": { "type": "string", "format": "date-time" },
    "contained_at": { "type": "string", "format": "date-time" },
    "status": { "type": "string", "enum": ["open", "investigating", "contained", "eradicated", "recovered", "closed"] },
    "affected_assets": { "type": "array", "items": { "type": "string" }, "uniqueItems": true },
    "root_cause": { "type": "string" },
    "lessons_learned": { "type": "string" }
  },
  "additionalProperties": false
}


3.4 asset

{
  "$id": "https://complisec.local/schemas/asset.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Asset Custom Fields",
  "type": "object",
  "required": ["asset_id", "asset_type", "owner", "criticality"],
  "properties": {
    "asset_id": { "type": "string" },
    "asset_type": { "type": "string", "enum": ["hardware", "software", "data", "service", "personnel", "facility"] },
    "owner": { "type": "string" },
    "location": { "type": "string" },
    "inventory_number": { "type": "string" },
    "criticality": { "type": "string", "enum": ["low", "medium", "high", "mission_critical"] }
  },
  "additionalProperties": false
}


3.5 control

{
  "$id": "https://complisec.local/schemas/control.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Control Custom Fields",
  "type": "object",
  "required": ["control_id", "family", "implementation_state"],
  "properties": {
    "control_id": { "type": "string" },
    "family": { "type": "string" },
    "implementation_state": { "type": "string", "enum": ["planned", "partially_implemented", "implemented", "not_applicable"] },
    "verification": { "type": "string", "enum": ["manual", "automated", "hybrid"] }
  },
  "additionalProperties": false
}


3.6 training_material

{
  "$id": "https://complisec.local/schemas/training.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Training Material Custom Fields",
  "type": "object",
  "required": ["material_id", "material_type", "audience"],
  "properties": {
    "material_id": { "type": "string" },
    "material_type": { "type": "string", "enum": ["policy_brief", "video", "slide_deck", "quiz", "article"] },
    "duration_minutes": { "type": "integer", "minimum": 0 },
    "audience": { "type": "string" }
  },
  "additionalProperties": false
}


3.7 audit_record

{
  "$id": "https://complisec.local/schemas/audit.custom.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Audit Record Custom Fields",
  "type": "object",
  "required": ["audit_id", "activity", "actor", "result", "occurred_at"],
  "properties": {
    "audit_id": { "type": "string" },
    "activity": { "type": "string" },
    "actor": { "type": "string" },
    "result": { "type": "string", "enum": ["success", "failure", "partial"] },
    "occurred_at": { "type": "string", "format": "date-time" }
  },
  "additionalProperties": false
}

4) Примеры валидных документов

4.1 Policy

{
  "id": "a3b1e2f4-9c7d-4d1a-8a3f-8d4b6c9f1a23",
  "tenant_id": "9f1a2b3c-4d5e-4f60-8a7b-2c3d4e5f6a70",
  "doc_type": "policy",
  "title": "Политика управления доступом",
  "summary": "Требования к управлению учетными записями и правами доступа.",
  "language": "ru",
  "tags": ["ISO27001", "AccessControl", "IAM"],
  "source": {
    "ingest_method": "upload",
    "origin_name": "access_control_policy.pdf",
    "origin_uri": "file:///policies/access_control_policy.pdf",
    "file_mime": "application/pdf",
    "checksum_sha256": "9a7c2dc01c3a78a9ce2c4d0f4b1e6a35d2d6f4a3b2c19c7e8f1d9a2b3c4d5e6f"
  },
  "version": { "number": "1.2.0", "status": "approved", "valid_from": "2025-01-10T00:00:00Z" },
  "classification": { "confidentiality": "internal", "integrity": "high", "availability": "medium" },
  "compliance_mappings": [
    { "standard": "ISO27001", "clause": "A.9.2.1", "level": "must" },
    { "standard": "ISO27002", "clause": "9.2", "level": "should" }
  ],
  "relations": { "related_ids": [] },
  "owners": { "author": "Security Office", "reviewers": ["IT Operations"], "approver": "CISO" },
  "content": {
    "format": "markdown",
    "body": "# 1. Цель\nОбеспечить управление доступом...\n\n# 2. Область применения\nВсе ИС компании...",
    "sections": [
      { "number": "1", "title": "Цель", "text": "Обеспечить управление доступом...", "anchor": "goal" },
      { "number": "2", "title": "Область применения", "text": "Все ИС компании...", "anchor": "scope" }
    ]
  },
  "custom": {
    "policy_id": "POL-AC-001",
    "scope": "Все подразделения и подрядчики",
    "owner_department": "Security",
    "review_cycle_days": 365
  },
  "created_at": "2025-01-10T12:00:00Z",
  "updated_at": "2025-06-15T09:30:00Z"
}


4.2 Risk

{
  "id": "0a1b2c3d-4e5f-4a6b-8c9d-0e1f2a3b4c5d",
  "tenant_id": "9f1a2b3c-4d5e-4f60-8a7b-2c3d4e5f6a70",
  "doc_type": "risk",
  "title": "Риск: компрометация учетных записей",
  "summary": "Фишинг и подбор паролей ведут к компрометации.",
  "language": "ru",
  "tags": ["risk", "IAM", "phishing"],
  "source": { "ingest_method": "manual", "origin_name": "risk_register" },
  "version": { "number": "1.0.0", "status": "approved" },
  "classification": { "confidentiality": "internal", "integrity": "high", "availability": "medium" },
  "compliance_mappings": [{ "standard": "ISO27001", "clause": "A.5", "level": "should" }],
  "owners": { "author": "Risk Manager", "reviewers": ["Security Office"], "approver": "CISO" },
  "content": { "format": "plain", "body": "Описание сценариев угроз и мер снижения." },
  "custom": {
    "risk_id": "RSK-0007",
    "asset_id": "AST-AD-001",
    "category": "Identity",
    "likelihood": "high",
    "impact": "high",
    "score_method": "matrix_5x5",
    "inherent_score": 20,
    "residual_score": 8,
    "controls": ["CTRL-MFA-001", "CTRL-PWD-002"],
    "status": "treated"
  },
  "created_at": "2025-02-01T10:00:00Z",
  "updated_at": "2025-08-21T08:45:00Z"
}


4.3 Incident

{
  "id": "6f7e8d9c-1b2a-4c3d-8e9f-0a1b2c3d4e5f",
  "tenant_id": "9f1a2b3c-4d5e-4f60-8a7b-2c3d4e5f6a70",
  "doc_type": "incident",
  "title": "Утечка данных CRM через фишинг",
  "summary": "Компрометация учетной записи сотрудника отдела продаж.",
  "language": "ru",
  "tags": ["incident", "phishing", "data_breach"],
  "source": { "ingest_method": "api", "origin_name": "siem" },
  "version": { "number": "1.0.0", "status": "approved" },
  "classification": { "confidentiality": "confidential", "integrity": "high", "availability": "medium" },
  "compliance_mappings": [{ "standard": "GDPR", "clause": "Art.33", "level": "must" }],
  "content": { "format": "markdown", "body": "## Ход инцидента\n..." },
  "custom": {
    "incident_id": "INC-2025-0142",
    "severity": "critical",
    "detected_at": "2025-09-10T04:12:00Z",
    "contained_at": "2025-09-10T08:25:00Z",
    "status": "recovered",
    "affected_assets": ["AST-CRM-001"],
    "root_cause": "Фишинг-переход по вредоносной ссылке",
    "lessons_learned": "Обязательная MFA, тренинги по фишингу"
  },
  "created_at": "2025-09-10T09:00:00Z",
  "updated_at": "2025-09-12T15:20:00Z"
}

5) Модель чанков и метаданных индекса

Схема чанка

{
  "$id": "https://complisec.local/schemas/chunk.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "CompliSec Chunk",
  "type": "object",
  "required": [
    "chunk_id",
    "parent_id",
    "tenant_id",
    "doc_type",
    "language",
    "content",
    "token_count",
    "chunk_index",
    "section_path",
    "created_at"
  ],
  "properties": {
    "chunk_id": {
      "type": "string",
      "description": "UUID v4",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "parent_id": {
      "type": "string",
      "description": "id исходного документа",
      "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
    },
    "tenant_id": { "type": "string" },
    "doc_type": { "type": "string" },
    "language": { "type": "string" },
    "tags": { "type": "array", "items": { "type": "string" }, "uniqueItems": true },
    "compliance_mappings": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Плоский список, например ISO27001:A.9.2.1"
    },
    "section_path": { "type": "string", "description": "Напр. 1>2.3>ПАРОЛИ" },
    "content": { "type": "string" },
    "token_count": { "type": "integer", "minimum": 0 },
    "chunk_index": { "type": "integer", "minimum": 0 },
    "overlap_with_prev": { "type": "integer", "minimum": 0 },
    "embedding_vector_ref": { "type": "string", "description": "ключ/идентификатор эмбеддинга" },
    "created_at": { "type": "string", "format": "date-time" }
  },
  "additionalProperties": false
}


Пример чанка

{
  "chunk_id": "c2b5f2a1-2f9a-4a74-9f8d-6d3e9f2a7b10",
  "parent_id": "a3b1e2f4-9c7d-4d1a-8a3f-8d4b6c9f1a23",
  "tenant_id": "9f1a2b3c-4d5e-4f60-8a7b-2c3d4e5f6a70",
  "doc_type": "policy",
  "language": "ru",
  "tags": ["ISO27001", "IAM"],
  "compliance_mappings": ["ISO27001:A.9.2.1", "ISO27002:9.2"],
  "section_path": "2>Требования к учётным записям",
  "content": "Учётные записи пользователей создаются по заявке... MFA обязательно для удалённого доступа...",
  "token_count": 142,
  "chunk_index": 3,
  "overlap_with_prev": 40,
  "embedding_vector_ref": "qdrant:collections/complisec/points/c2b5f2a1",
  "created_at": "2025-06-15T09:30:05Z"
}

6) Правила нормализации и нарезки (рабочие)

Очистка текста: удалять номера страниц, колонтитулы, мусорные переносы, таблицы переводить в маркдаун-таблицы.

Единый формат: content.format = "markdown" предпочтителен.

Размер чанка: 800–1200 токенов для ru/en; overlap 120–180 токенов.

Границы чанка: стараться резать по заголовкам #, ##, по абзацам; избегать разрыва списков.

Обогащение метаданными: дублировать важные поля из документа в чанк (doc_type, tags, плоские compliance_mappings).

Язык: фиксировать language по BCP-47 (например, ru, en).

Версионность: при смене version.number — новый id документа (immutable запись), либо id прежний, но version.status="deprecated" у старой записи. Для RAG проще хранить каждую версию как отдельный документ и фильтровать version.status="approved".

7) Мини-чек-лист интеграции в ваш пайплайн (CompliSec/RiskNexus)

Конвертеры: PDF → текст (pdfplumber/pdftotext), DOCX → текст (python-docx/docx2txt), XLSX/CSV → нормализованные таблицы → Markdown.

Валидация: перед записью в БД — проверка Base Schema + нужной под-схемы custom.

Генерация compliance_mappings: маппинг из словаря соответствий (клауз-индекс).

Чанк-строитель: учитывает sections и размеры; пишет chunks.jsonl.

Индекс: векторка (Qdrant/Weaviate/Pinecone), фильтры по tenant_id, doc_type, tags, compliance_mappings.

Ответы RAG: k-NN → rerank → сборка контекста с учётом classification.confidentiality.