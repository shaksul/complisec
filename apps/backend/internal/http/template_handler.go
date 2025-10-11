package http

import (
	"log"
	"net/http"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TemplateHandler handles template-related HTTP requests
type TemplateHandler struct {
	templateService *domain.TemplateService
	validator       *validator.Validate
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(templateService *domain.TemplateService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		validator:       validator.New(),
	}
}

// Register registers template routes
func (h *TemplateHandler) Register(r fiber.Router) {
	templates := r.Group("/admin/templates")
	templates.Get("/", RequirePermission("admin"), h.listTemplates)
	templates.Get("/variables", RequirePermission("admin"), h.getTemplateVariables)
	templates.Get("/:id", RequirePermission("admin"), h.getTemplate)
	templates.Post("/", RequirePermission("admin"), h.createTemplate)
	templates.Put("/:id", RequirePermission("admin"), h.updateTemplate)
	templates.Delete("/:id", RequirePermission("admin"), h.deleteTemplate)
	templates.Post("/initialize-defaults", RequirePermission("admin"), h.initializeDefaultTemplates)

	// Inventory number rules
	inventory := r.Group("/admin/inventory-rules")
	inventory.Get("/", RequirePermission("admin"), h.listInventoryRules)
	inventory.Post("/", RequirePermission("admin"), h.createInventoryRule)
	inventory.Put("/:id", RequirePermission("admin"), h.updateInventoryRule)

	// Asset template operations
	assets := r.Group("/assets")
	assets.Post("/:id/fill-template", RequirePermission("assets.view"), h.fillTemplate)
	assets.Post("/:id/generate-inventory-number", RequirePermission("assets.create"), h.generateInventoryNumber)
}

// listTemplates returns all templates for a tenant
func (h *TemplateHandler) listTemplates(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	filters := make(map[string]interface{})
	if templateType := c.Query("template_type"); templateType != "" {
		filters["template_type"] = templateType
	}
	if isSystem := c.Query("is_system"); isSystem != "" {
		filters["is_system"] = isSystem == "true"
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	templates, err := h.templateService.ListTemplates(c.Context(), tenantID, filters)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.listTemplates: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": templates})
}

// getTemplate returns a single template
func (h *TemplateHandler) getTemplate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	templateID := c.Params("id")

	template, err := h.templateService.GetTemplate(c.Context(), templateID, tenantID)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.getTemplate: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": template})
}

// createTemplate creates a new template
func (h *TemplateHandler) createTemplate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	template, err := h.templateService.CreateTemplate(c.Context(), tenantID, userID, req)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.createTemplate: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": template})
}

// updateTemplate updates an existing template
func (h *TemplateHandler) updateTemplate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	templateID := c.Params("id")

	var req dto.UpdateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.templateService.UpdateTemplate(c.Context(), templateID, tenantID, req)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.updateTemplate: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Template updated successfully"})
}

// deleteTemplate deletes a template
func (h *TemplateHandler) deleteTemplate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	templateID := c.Params("id")

	err := h.templateService.DeleteTemplate(c.Context(), templateID, tenantID)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.deleteTemplate: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Template deleted successfully"})
}

// initializeDefaultTemplates creates default system templates
func (h *TemplateHandler) initializeDefaultTemplates(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	err := h.templateService.InitializeDefaultTemplates(c.Context(), tenantID, userID)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.initializeDefaultTemplates: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Default templates initialized successfully"})
}

// getTemplateVariables returns available template variables
func (h *TemplateHandler) getTemplateVariables(c *fiber.Ctx) error {
	variables := h.templateService.GetAvailableVariables()
	return c.JSON(fiber.Map{"data": variables})
}

// fillTemplate fills a template with asset data
func (h *TemplateHandler) fillTemplate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	assetID := c.Params("id")

	var req dto.FillTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Set asset ID from URL
	req.AssetID = assetID

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	result, err := h.templateService.FillTemplate(c.Context(), tenantID, userID, req)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.fillTemplate: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": result})
}

// listInventoryRules returns all inventory number rules
func (h *TemplateHandler) listInventoryRules(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	rules, err := h.templateService.ListInventoryRules(c.Context(), tenantID)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.listInventoryRules: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": rules})
}

// createInventoryRule creates a new inventory number rule
func (h *TemplateHandler) createInventoryRule(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req dto.CreateInventoryRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	rule, err := h.templateService.CreateInventoryRule(c.Context(), tenantID, req)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.createInventoryRule: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": rule})
}

// updateInventoryRule updates an existing inventory number rule
func (h *TemplateHandler) updateInventoryRule(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	ruleID := c.Params("id")

	var req dto.UpdateInventoryRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.templateService.UpdateInventoryRule(c.Context(), ruleID, tenantID, req)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.updateInventoryRule: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Inventory rule updated successfully"})
}

// generateInventoryNumber generates a new inventory number for an asset
func (h *TemplateHandler) generateInventoryNumber(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	var req dto.GenerateInventoryNumberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	result, err := h.templateService.GenerateInventoryNumber(c.Context(), tenantID, req.AssetType, req.AssetClass)
	if err != nil {
		log.Printf("ERROR: TemplateHandler.generateInventoryNumber: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": result})
}

