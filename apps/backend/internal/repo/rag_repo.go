package repo

import (
	"context"
	"time"
)

type IndexedDocument struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	DocumentID         string     `json:"document_id"`
	Status             string     `json:"status"`
	ErrorMessage       *string    `json:"error_message,omitempty"`
	RetryCount         int        `json:"retry_count"`
	MaxRetries         int        `json:"max_retries"`
	GraphRAGDocID      *string    `json:"graphrag_doc_id,omitempty"`
	ChunksCount        int        `json:"chunks_count"`
	EntitiesCount      int        `json:"entities_count"`
	RelationshipsCount int        `json:"relationships_count"`
	IndexedAt          *time.Time `json:"indexed_at,omitempty"`
	LastRetryAt        *time.Time `json:"last_retry_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type RAGRepo struct {
	db *DB
}

func NewRAGRepo(db *DB) *RAGRepo {
	return &RAGRepo{db: db}
}

func (r *RAGRepo) CreateIndexRecord(ctx context.Context, doc IndexedDocument) error {
	query := `INSERT INTO rag_indexed_documents 
        (tenant_id, document_id, status) 
        VALUES ($1, $2, $3)
        ON CONFLICT (tenant_id, document_id) 
        DO UPDATE SET status = EXCLUDED.status, updated_at = NOW()`
	_, err := r.db.Exec(query, doc.TenantID, doc.DocumentID, doc.Status)
	return err
}

func (r *RAGRepo) UpdateIndexStatus(ctx context.Context, docID string, status string, graphRAGID *string, chunksCount, entitiesCount, relationshipsCount int) error {
	query := `UPDATE rag_indexed_documents 
        SET status = $1, graphrag_doc_id = $2, chunks_count = $3, 
            entities_count = $4, relationships_count = $5, 
            indexed_at = NOW()
        WHERE document_id = $6`
	_, err := r.db.Exec(query, status, graphRAGID, chunksCount, entitiesCount, relationshipsCount, docID)
	return err
}

func (r *RAGRepo) UpdateIndexError(ctx context.Context, docID string, errorMsg string) error {
	query := `UPDATE rag_indexed_documents 
        SET status = 'retrying', error_message = $1, 
            retry_count = retry_count + 1, last_retry_at = NOW()
        WHERE document_id = $2 AND retry_count < max_retries`
	result, err := r.db.Exec(query, errorMsg, docID)
	if err != nil {
		return err
	}

	// Если превышен лимит retry, пометить как failed
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		query = `UPDATE rag_indexed_documents 
            SET status = 'failed', error_message = $1
            WHERE document_id = $2`
		_, err = r.db.Exec(query, errorMsg, docID)
	}

	return err
}

func (r *RAGRepo) GetIndexedDocuments(ctx context.Context, tenantID string) ([]IndexedDocument, error) {
	query := `SELECT id, tenant_id, document_id, status, error_message, 
        retry_count, max_retries, graphrag_doc_id, chunks_count, entities_count, 
        relationships_count, indexed_at, last_retry_at, created_at, updated_at 
        FROM rag_indexed_documents WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []IndexedDocument
	for rows.Next() {
		var doc IndexedDocument
		err := rows.Scan(&doc.ID, &doc.TenantID, &doc.DocumentID, &doc.Status,
			&doc.ErrorMessage, &doc.RetryCount, &doc.MaxRetries, &doc.GraphRAGDocID,
			&doc.ChunksCount, &doc.EntitiesCount, &doc.RelationshipsCount,
			&doc.IndexedAt, &doc.LastRetryAt, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *RAGRepo) GetIndexedDocument(ctx context.Context, documentID string) (*IndexedDocument, error) {
	query := `SELECT id, tenant_id, document_id, status, error_message, 
        retry_count, max_retries, graphrag_doc_id, chunks_count, entities_count, 
        relationships_count, indexed_at, last_retry_at, created_at, updated_at 
        FROM rag_indexed_documents WHERE document_id = $1`

	var doc IndexedDocument
	err := r.db.QueryRow(query, documentID).Scan(&doc.ID, &doc.TenantID, &doc.DocumentID,
		&doc.Status, &doc.ErrorMessage, &doc.RetryCount, &doc.MaxRetries, &doc.GraphRAGDocID,
		&doc.ChunksCount, &doc.EntitiesCount, &doc.RelationshipsCount,
		&doc.IndexedAt, &doc.LastRetryAt, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *RAGRepo) LogQuery(ctx context.Context, tenantID, userID, query string, useGraph bool, sourcesCount, responseTimeMs int) error {
	sql := `INSERT INTO rag_query_log (tenant_id, user_id, query, use_graph, sources_count, response_time_ms) 
        VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(sql, tenantID, userID, query, useGraph, sourcesCount, responseTimeMs)
	return err
}

func (r *RAGRepo) GetDocumentsToRetry(ctx context.Context) ([]IndexedDocument, error) {
	query := `SELECT id, tenant_id, document_id, status, error_message, 
        retry_count, max_retries, graphrag_doc_id, chunks_count, entities_count, 
        relationships_count, indexed_at, last_retry_at, created_at, updated_at 
        FROM rag_indexed_documents 
        WHERE status = 'retrying' 
        AND retry_count < max_retries
        AND (last_retry_at IS NULL OR last_retry_at < NOW() - INTERVAL '5 minutes')
        ORDER BY created_at ASC
        LIMIT 10`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []IndexedDocument
	for rows.Next() {
		var doc IndexedDocument
		err := rows.Scan(&doc.ID, &doc.TenantID, &doc.DocumentID, &doc.Status,
			&doc.ErrorMessage, &doc.RetryCount, &doc.MaxRetries, &doc.GraphRAGDocID,
			&doc.ChunksCount, &doc.EntitiesCount, &doc.RelationshipsCount,
			&doc.IndexedAt, &doc.LastRetryAt, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

