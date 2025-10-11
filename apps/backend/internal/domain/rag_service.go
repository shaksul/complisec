package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"risknexus/backend/internal/repo"
	"time"
)

type RAGService struct {
	repo            *repo.RAGRepo
	docRepo         *repo.DocumentRepo
	graphRAGURL     string
	docProcessorURL string
}

func NewRAGService(r *repo.RAGRepo, dr *repo.DocumentRepo) *RAGService {
	return &RAGService{
		repo:            r,
		docRepo:         dr,
		graphRAGURL:     os.Getenv("GRAPHRAG_URL"),
		docProcessorURL: os.Getenv("DOC_PROCESSOR_URL"),
	}
}

// IndexDocument - индексация одного документа
func (s *RAGService) IndexDocument(ctx context.Context, tenantID, documentID string) error {
	// 1. Получить документ из БД
	doc, err := s.docRepo.GetDocumentByID(ctx, documentID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// 2. Создать запись индексации
	s.repo.CreateIndexRecord(ctx, repo.IndexedDocument{
		TenantID:   tenantID,
		DocumentID: documentID,
		Status:     "processing",
	})

	// 3. Конвертировать документ в текст
	filePath := doc.FilePath
	// Преобразуем путь для doc-processor контейнера
	if filePath[0] != '/' {
		filePath = "/" + filePath
	}
	text, err := s.convertToText(ctx, filePath, true) // use_ocr = true
	if err != nil {
		s.repo.UpdateIndexError(ctx, documentID, fmt.Sprintf("conversion error: %v", err))
		return err
	}

	// 4. Отправить в GraphRAG
	graphRAGResp, err := s.sendToGraphRAG(ctx, documentID, doc.Title, text, map[string]interface{}{
		"tenant_id": tenantID,
		"doc_type":  doc.Type,
		"version":   doc.Version,
	})
	if err != nil {
		s.repo.UpdateIndexError(ctx, documentID, fmt.Sprintf("graphrag error: %v", err))
		return err
	}

	// 5. Обновить статус (безопасное извлечение значений)
	var graphRAGID *string
	if val, ok := graphRAGResp["doc_id"].(string); ok && val != "" {
		graphRAGID = &val
	}

	chunksCount := 0
	if val, ok := graphRAGResp["chunks_count"].(float64); ok {
		chunksCount = int(val)
	}

	entitiesCount := 0
	if val, ok := graphRAGResp["entities_count"].(float64); ok {
		entitiesCount = int(val)
	}

	relationshipsCount := 0
	if val, ok := graphRAGResp["relationships_count"].(float64); ok {
		relationshipsCount = int(val)
	}

	return s.repo.UpdateIndexStatus(ctx, documentID, "indexed", graphRAGID, chunksCount, entitiesCount, relationshipsCount)
}

// convertToText - конвертация через doc-processor
func (s *RAGService) convertToText(ctx context.Context, filePath string, useOCR bool) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"file_path": filePath,
		"use_ocr":   useOCR,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.docProcessorURL+"/convert", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("doc-processor error (status %d): %s", resp.StatusCode, body)
	}

	var result struct {
		Text      string `json:"text"`
		CharCount int    `json:"char_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Text, nil
}

// sendToGraphRAG - отправка документа в GraphRAG
func (s *RAGService) sendToGraphRAG(ctx context.Context, docID, title, content string, metadata map[string]interface{}) (map[string]interface{}, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"id":       docID,
		"title":    title,
		"content":  content,
		"metadata": metadata,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.graphRAGURL+"/index", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("graphrag error (status %d): %s", resp.StatusCode, body)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// Query - поиск с RAG
func (s *RAGService) Query(ctx context.Context, tenantID, userID, query string, useGraph bool, topK int) (*RAGQueryResult, error) {
	startTime := time.Now()

	reqBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"use_graph": useGraph,
		"top_k":     topK,
		"filter": map[string]string{
			"tenant_id": tenantID,
		},
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.graphRAGURL+"/query", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("graphrag query error (status %d): %s", resp.StatusCode, body)
	}

	var graphRAGResponse struct {
		Sources []RAGSource `json:"sources"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&graphRAGResponse); err != nil {
		return nil, err
	}

	// Логируем запрос
	responseTime := int(time.Since(startTime).Milliseconds())
	s.repo.LogQuery(ctx, tenantID, userID, query, useGraph, len(graphRAGResponse.Sources), responseTime)

	// Возвращаем только источники, без ответа (ответ генерирует AI провайдер)
	return &RAGQueryResult{
		Answer:  "", // Пустой, ответ генерирует AI провайдер
		Sources: graphRAGResponse.Sources,
	}, nil
}

func (s *RAGService) GetIndexedDocuments(ctx context.Context, tenantID string) ([]repo.IndexedDocument, error) {
	return s.repo.GetIndexedDocuments(ctx, tenantID)
}

func (s *RAGService) GetAllDocuments(ctx context.Context, tenantID string) ([]repo.Document, error) {
	// Вызываем ListAllDocuments с пустыми фильтрами
	return s.docRepo.ListAllDocuments(ctx, tenantID, map[string]interface{}{})
}

type RAGQueryResult struct {
	Answer  string      `json:"answer"`
	Sources []RAGSource `json:"sources"`
}

type RAGSource struct {
	DocumentID string  `json:"document_id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
}
