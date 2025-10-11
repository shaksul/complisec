package http

import (
	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	service *domain.AIService
}

func NewAIHandler(s *domain.AIService) *AIHandler {
	return &AIHandler{service: s}
}

func (h *AIHandler) Register(r fiber.Router) {
	r.Get("/ai/providers", h.listProviders)
	r.Get("/ai/providers/:id", h.getProvider)
	r.Post("/ai/providers", h.createProvider)
	r.Put("/ai/providers/:id", h.updateProvider)
	r.Delete("/ai/providers/:id", h.deleteProvider)
	r.Post("/ai/query", h.query)
}

func (h *AIHandler) listProviders(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	items, err := h.service.List(c.Context(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *AIHandler) getProvider(c *fiber.Ctx) error {
	id := c.Params("id")
	provider, err := h.service.Get(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if provider == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Провайдер не найден"})
	}
	return c.JSON(fiber.Map{"data": provider})
}

func (h *AIHandler) createProvider(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var dto dto.CreateAIProviderDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	// Если модели не указаны, используем дефолтные
	models := dto.Models
	if len(models) == 0 {
		models = []string{"llama3.2"}
	}
	defaultModel := dto.DefaultModel
	if defaultModel == "" {
		defaultModel = "llama3.2"
	}

	p := repo.AIProvider{
		TenantID:       tenantID,
		Name:           dto.Name,
		BaseURL:        dto.BaseURL,
		APIKey:         &dto.APIKey,
		Roles:          dto.Roles,
		Models:         models,
		DefaultModel:   defaultModel,
		PromptTemplate: &dto.PromptTemplate,
		IsActive:       true,
	}
	if err := h.service.Create(c.Context(), p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *AIHandler) updateProvider(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dto.UpdateAIProviderDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	// Если модели не указаны, используем дефолтные
	models := dto.Models
	if len(models) == 0 {
		models = []string{"llama3.2"}
	}
	defaultModel := dto.DefaultModel
	if defaultModel == "" {
		defaultModel = "llama3.2"
	}

	p := repo.AIProvider{
		ID:             id,
		Name:           dto.Name,
		BaseURL:        dto.BaseURL,
		APIKey:         &dto.APIKey,
		Roles:          dto.Roles,
		Models:         models,
		DefaultModel:   defaultModel,
		PromptTemplate: &dto.PromptTemplate,
		IsActive:       dto.IsActive,
	}
	if err := h.service.Update(c.Context(), p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *AIHandler) deleteProvider(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *AIHandler) query(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, _ := c.Locals("user_id").(string)

	var req struct {
		ProviderID string `json:"provider_id"`
		Role       string `json:"role"`
		Input      string `json:"input"`
		UseRAG     bool   `json:"use_rag"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	// Если use_rag = true, используем QueryWithRAG
	if req.UseRAG {
		result, err := h.service.QueryWithRAG(c.Context(), tenantID, userID, req.ProviderID, req.Role, req.Input, true)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	}

	// Обычный запрос без RAG
	prov, err := h.service.Get(c.Context(), req.ProviderID)
	if err != nil || prov == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Provider not found"})
	}

	out, err := h.service.Query(c.Context(), *prov, req.Role, req.Input, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"output": out})
}
