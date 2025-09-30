package http

import (
	"strconv"

	"risknexus/backend/internal/repo"

	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	auditRepo *repo.AuditRepo
}

func NewAuditHandler(auditRepo *repo.AuditRepo) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

func (h *AuditHandler) Register(r fiber.Router) {
	audit := r.Group("/audit")
	audit.Get("/", RequirePermission("audit.view"), h.getAuditLogs)
}

func (h *AuditHandler) getAuditLogs(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Параметры пагинации
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Фильтры
	var actorID, action, entity *string
	if actorIDStr := c.Query("actor_id"); actorIDStr != "" {
		actorID = &actorIDStr
	}
	if actionStr := c.Query("action"); actionStr != "" {
		action = &actionStr
	}
	if entityStr := c.Query("entity"); entityStr != "" {
		entity = &entityStr
	}

	logs, err := h.auditRepo.GetAuditLogs(c.Context(), tenantID, limit, offset, actorID, action, entity)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": logs,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"count":  len(logs),
		},
	})
}
