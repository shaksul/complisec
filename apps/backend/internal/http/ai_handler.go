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
	r.Post("/ai/providers", h.createProvider)
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

func (h *AIHandler) createProvider(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var dto dto.CreateAIProviderDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}
	p := repo.AIProvider{TenantID: tenantID, Name: dto.Name, BaseURL: dto.BaseURL, APIKey: &dto.APIKey, Roles: dto.Roles, PromptTemplate: &dto.PromptTemplate, IsActive: true}
	if err := h.service.Create(c.Context(), p); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *AIHandler) query(c *fiber.Ctx) error {
	var req dto.QueryAIRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}
	prov := repo.AIProvider{ID: req.ProviderID, BaseURL: "http://localhost:11434/api/chat"} // заглушка
	out, err := h.service.Query(c.Context(), prov, req.Role, req.Input, req.Context)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dto.QueryAIResponse{Output: out})
}
