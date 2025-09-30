package http

import (
	"context"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type IncidentHandler struct {
	incidentService *domain.IncidentService
}

func NewIncidentHandler(incidentService *domain.IncidentService) *IncidentHandler {
	return &IncidentHandler{incidentService: incidentService}
}

func (h *IncidentHandler) Register(r fiber.Router) {
	incidents := r.Group("/incidents")
	incidents.Get("/", h.listIncidents)
	incidents.Post("/", RequirePermission("incidents.create"), h.createIncident)
	incidents.Get("/:id", h.getIncident)
	incidents.Put("/:id", RequirePermission("incidents.edit"), h.updateIncident)
	incidents.Patch("/:id/status", RequirePermission("incidents.edit"), h.updateIncidentStatus)
	incidents.Delete("/:id", RequirePermission("incidents.delete"), h.deleteIncident)
}

func (h *IncidentHandler) listIncidents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	incidents, err := h.incidentService.ListIncidents(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": incidents})
}

func (h *IncidentHandler) createIncident(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Severity    string  `json:"severity"`
		AssetID     *string `json:"asset_id"`
		RiskID      *string `json:"risk_id"`
		AssignedTo  *string `json:"assigned_to"`
		CreatedBy   *string `json:"created_by"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	incident, err := h.incidentService.CreateIncident(context.Background(), tenantID, req.Title, req.Severity, req.Description, req.AssetID, req.RiskID, req.AssignedTo, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": incident})
}

func (h *IncidentHandler) getIncident(c *fiber.Ctx) error {
	id := c.Params("id")
	incident, err := h.incidentService.GetIncident(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if incident == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Incident not found"})
	}

	return c.JSON(fiber.Map{"data": incident})
}

func (h *IncidentHandler) updateIncident(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Severity    string  `json:"severity"`
		AssetID     *string `json:"asset_id"`
		RiskID      *string `json:"risk_id"`
		AssignedTo  *string `json:"assigned_to"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.incidentService.UpdateIncident(context.Background(), id, req.Title, req.Severity, req.Description, req.AssetID, req.RiskID, req.AssignedTo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Incident updated successfully"})
}

func (h *IncidentHandler) updateIncidentStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.incidentService.UpdateIncidentStatus(context.Background(), id, req.Status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Incident status updated successfully"})
}

func (h *IncidentHandler) deleteIncident(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.incidentService.DeleteIncident(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Incident deleted successfully"})
}
