package http

import (
	"context"
	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type RAGHandler struct {
	service *domain.RAGService
}

func NewRAGHandler(s *domain.RAGService) *RAGHandler {
	return &RAGHandler{service: s}
}

func (h *RAGHandler) Register(r fiber.Router) {
	r.Post("/rag/index/:document_id", h.indexDocument)
	r.Post("/rag/index-all", h.indexAllDocuments)
	r.Get("/rag/indexed", h.getIndexedDocuments)
	r.Post("/rag/query", h.query)
}

func (h *RAGHandler) indexDocument(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("document_id")

	// Создаем новый context для фоновой задачи
	ctx := context.Background()

	// Запускаем индексацию в фоне
	go func() {
		if err := h.service.IndexDocument(ctx, tenantID, documentID); err != nil {
			// Логируем ошибку (в production использовать logger)
			println("IndexDocument error:", err.Error())
		}
	}()

	return c.JSON(fiber.Map{"status": "indexing_started"})
}

func (h *RAGHandler) indexAllDocuments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Получаем все документы
	docs, err := h.service.GetAllDocuments(c.Context(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Создаем новый context для фоновой задачи
	ctx := context.Background()

	// Запускаем индексацию всех документов в фоне
	go func() {
		for _, doc := range docs {
			if err := h.service.IndexDocument(ctx, tenantID, doc.ID); err != nil {
				println("IndexAllDocuments error for doc", doc.ID, ":", err.Error())
			}
		}
	}()

	return c.JSON(fiber.Map{
		"status":     "indexing_started",
		"total_docs": len(docs),
		"message":    "Индексация запущена в фоне",
	})
}

func (h *RAGHandler) getIndexedDocuments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	docs, err := h.service.GetIndexedDocuments(c.Context(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": docs})
}

func (h *RAGHandler) query(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		userID = "" // Optional
	}

	var req struct {
		Query    string `json:"query"`
		UseGraph bool   `json:"use_graph"`
		TopK     int    `json:"top_k"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.TopK == 0 {
		req.TopK = 5
	}

	result, err := h.service.Query(c.Context(), tenantID, userID, req.Query, req.UseGraph, req.TopK)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
