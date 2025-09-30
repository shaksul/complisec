package http

import (
	"context"
	"time"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type TrainingHandler struct {
	trainingService *domain.TrainingService
}

func NewTrainingHandler(trainingService *domain.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingService: trainingService}
}

func (h *TrainingHandler) Register(r fiber.Router) {
	training := r.Group("/training")
	training.Get("/materials", h.listMaterials)
	training.Post("/materials", RequirePermission("training.create"), h.createMaterial)
	training.Post("/assign", RequirePermission("training.assign"), h.createAssignment)
	training.Get("/assignments", h.getUserAssignments)
	training.Post("/assignments/:id/complete", h.completeAssignment)
}

func (h *TrainingHandler) listMaterials(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	materials, err := h.trainingService.ListMaterials(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": materials})
}

func (h *TrainingHandler) createMaterial(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		Type        string  `json:"type"`
		URI         string  `json:"uri"`
		CreatedBy   *string `json:"created_by"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	material, err := h.trainingService.CreateMaterial(context.Background(), tenantID, req.Title, req.Type, req.URI, req.Description, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": material})
}

func (h *TrainingHandler) createAssignment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req struct {
		MaterialID string  `json:"material_id"`
		UserID     string  `json:"user_id"`
		DueAt      *string `json:"due_at"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var dueAt *time.Time
	if req.DueAt != nil {
		if parsed, err := time.Parse(time.RFC3339, *req.DueAt); err == nil {
			dueAt = &parsed
		}
	}

	assignment, err := h.trainingService.CreateAssignment(context.Background(), tenantID, req.MaterialID, req.UserID, dueAt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": assignment})
}

func (h *TrainingHandler) getUserAssignments(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	assignments, err := h.trainingService.GetUserAssignments(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": assignments})
}

func (h *TrainingHandler) completeAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.trainingService.CompleteAssignment(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Assignment completed successfully"})
}
