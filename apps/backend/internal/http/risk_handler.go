package http

import (
	"context"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type RiskHandler struct {
	riskService *domain.RiskService
}

func NewRiskHandler(riskService *domain.RiskService) *RiskHandler {
	return &RiskHandler{riskService: riskService}
}

func (h *RiskHandler) Register(r fiber.Router) {
	risks := r.Group("/risks")
	risks.Get("/", h.listRisks)
	risks.Post("/", RequirePermission("risks.create"), h.createRisk)
	risks.Get("/:id", h.getRisk)
	risks.Put("/:id", RequirePermission("risks.edit"), h.updateRisk)
	risks.Patch("/:id/status", RequirePermission("risks.edit"), h.updateRiskStatus)
	risks.Delete("/:id", RequirePermission("risks.delete"), h.deleteRisk)
}

func (h *RiskHandler) listRisks(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	risks, err := h.riskService.ListRisks(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": risks})
}

func (h *RiskHandler) createRisk(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		Likelihood  int     `json:"likelihood"`
		Impact      int     `json:"impact"`
		OwnerID     *string `json:"owner_id"`
		AssetID     *string `json:"asset_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	risk, err := h.riskService.CreateRisk(context.Background(), tenantID, req.Title, req.Description, req.Category, req.Likelihood, req.Impact, req.OwnerID, req.AssetID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": risk})
}

func (h *RiskHandler) getRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	risk, err := h.riskService.GetRisk(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if risk == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Risk not found"})
	}

	return c.JSON(fiber.Map{"data": risk})
}

func (h *RiskHandler) updateRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		Likelihood  int     `json:"likelihood"`
		Impact      int     `json:"impact"`
		OwnerID     *string `json:"owner_id"`
		AssetID     *string `json:"asset_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.riskService.UpdateRisk(context.Background(), id, req.Title, req.Description, req.Category, req.Likelihood, req.Impact, req.OwnerID, req.AssetID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Risk updated successfully"})
}

func (h *RiskHandler) updateRiskStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.riskService.UpdateRiskStatus(context.Background(), id, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Risk status updated successfully"})
}

func (h *RiskHandler) deleteRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.riskService.DeleteRisk(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Risk deleted successfully"})
}
