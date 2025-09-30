package http

import (
	"context"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/gofiber/fiber/v2"
)

type ComplianceHandler struct {
	service *domain.ComplianceService
}

func NewComplianceHandler(s *domain.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{service: s}
}

func (h *ComplianceHandler) Register(r fiber.Router) {
	// Standards
	r.Get("/compliance/standards", h.listStandards)
	r.Post("/compliance/standards", h.createStandard)

	// Requirements
	r.Get("/compliance/standards/:id/requirements", h.listRequirements)
	r.Post("/compliance/requirements", h.createRequirement)

	// Assessments
	r.Get("/compliance/assessments", h.listAssessments)
	r.Post("/compliance/assessments", h.createAssessment)
	r.Put("/compliance/assessments/:id", h.updateAssessment)

	// Gaps
	r.Get("/compliance/assessments/:id/gaps", h.listGaps)
	r.Post("/compliance/gaps", h.createGap)
	r.Put("/compliance/gaps/:id", h.updateGap)
}

func (h *ComplianceHandler) listStandards(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	items, err := h.service.ListStandards(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *ComplianceHandler) createStandard(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var dto dto.CreateComplianceStandardDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	standard := repo.ComplianceStandard{
		TenantID:    tenantID,
		Name:        dto.Name,
		Code:        dto.Code,
		Description: dto.Description,
		Version:     dto.Version,
		IsActive:    true,
	}

	if err := h.service.CreateStandard(context.Background(), standard); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *ComplianceHandler) listRequirements(c *fiber.Ctx) error {
	standardID := c.Params("id")
	items, err := h.service.ListRequirements(context.Background(), standardID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *ComplianceHandler) createRequirement(c *fiber.Ctx) error {
	var dto dto.CreateComplianceRequirementDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	requirement := repo.ComplianceRequirement{
		StandardID:  dto.StandardID,
		Code:        dto.Code,
		Title:       dto.Title,
		Description: dto.Description,
		Category:    dto.Category,
		Priority:    dto.Priority,
		IsMandatory: dto.IsMandatory,
	}

	if err := h.service.CreateRequirement(context.Background(), requirement); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *ComplianceHandler) listAssessments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	items, err := h.service.ListAssessments(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *ComplianceHandler) createAssessment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var dto dto.CreateComplianceAssessmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	assessment := repo.ComplianceAssessment{
		TenantID:       tenantID,
		RequirementID:  dto.RequirementID,
		Status:         dto.Status,
		Evidence:       dto.Evidence,
		AssessorID:     dto.AssessorID,
		AssessedAt:     dto.AssessedAt,
		NextReviewDate: dto.NextReviewDate,
		Notes:          dto.Notes,
	}

	if err := h.service.CreateAssessment(context.Background(), assessment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *ComplianceHandler) updateAssessment(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dto.UpdateComplianceAssessmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	assessment := repo.ComplianceAssessment{
		Status:         dto.Status,
		Evidence:       dto.Evidence,
		AssessorID:     dto.AssessorID,
		AssessedAt:     dto.AssessedAt,
		NextReviewDate: dto.NextReviewDate,
		Notes:          dto.Notes,
	}

	if err := h.service.UpdateAssessment(context.Background(), id, assessment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *ComplianceHandler) listGaps(c *fiber.Ctx) error {
	assessmentID := c.Params("id")
	items, err := h.service.ListGaps(context.Background(), assessmentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": items})
}

func (h *ComplianceHandler) createGap(c *fiber.Ctx) error {
	var dto dto.CreateComplianceGapDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	gap := repo.ComplianceGap{
		AssessmentID:    dto.AssessmentID,
		Title:           dto.Title,
		Description:     dto.Description,
		Severity:        dto.Severity,
		Status:          dto.Status,
		RemediationPlan: dto.RemediationPlan,
		TargetDate:      dto.TargetDate,
		ResponsibleID:   dto.ResponsibleID,
	}

	if err := h.service.CreateGap(context.Background(), gap); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}

func (h *ComplianceHandler) updateGap(c *fiber.Ctx) error {
	id := c.Params("id")
	var dto dto.UpdateComplianceGapDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad input"})
	}

	gap := repo.ComplianceGap{
		Status:          dto.Status,
		RemediationPlan: dto.RemediationPlan,
		TargetDate:      dto.TargetDate,
		ResponsibleID:   dto.ResponsibleID,
	}

	if err := h.service.UpdateGap(context.Background(), id, gap); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "ok"})
}
