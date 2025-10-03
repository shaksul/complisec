package http

import (
	"log"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IncidentHandler struct {
	incidentService domain.IncidentServiceInterface
	validator       *validator.Validate
}

func NewIncidentHandler(incidentService domain.IncidentServiceInterface) *IncidentHandler {
	return &IncidentHandler{
		incidentService: incidentService,
		validator:       validator.New(),
	}
}

func (h *IncidentHandler) Register(r fiber.Router) {
	incidents := r.Group("/incidents")
	incidents.Get("/", RequirePermission("incidents.view"), h.listIncidents)
	incidents.Post("/", RequirePermission("incidents.create"), h.createIncident)
	incidents.Get("/metrics", RequirePermission("incidents.report"), h.getIncidentMetrics)
	incidents.Get("/:id", RequirePermission("incidents.view"), h.getIncident)
	incidents.Put("/:id", RequirePermission("incidents.edit"), h.updateIncident)
	incidents.Delete("/:id", RequirePermission("incidents.delete"), h.deleteIncident)
	incidents.Put("/:id/status", RequirePermission("incidents.edit"), h.updateIncidentStatus)
	incidents.Post("/:id/comments", RequirePermission("incidents.edit"), h.addComment)
	incidents.Get("/:id/comments", RequirePermission("incidents.view"), h.getComments)
	incidents.Post("/:id/actions", RequirePermission("incidents.edit"), h.addAction)
	incidents.Get("/:id/actions", RequirePermission("incidents.view"), h.getActions)
	incidents.Put("/:id/actions/:actionId", RequirePermission("incidents.edit"), h.updateAction)
	incidents.Delete("/:id/actions/:actionId", RequirePermission("incidents.delete"), h.deleteAction)
}

func (h *IncidentHandler) listIncidents(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// Parse filters
	req := dto.IncidentListRequest{
		Page:        page,
		PageSize:    pageSize,
		Status:      c.Query("status"),
		Criticality: c.Query("criticality"),
		Category:    c.Query("category"),
		AssetID:     c.Query("asset_id"),
		RiskID:      c.Query("risk_id"),
		AssignedTo:  c.Query("assigned_to"),
		Search:      c.Query("search"),
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.listIncidents validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request parameters",
			"details": err.Error(),
		})
	}

	incidents, total, err := h.incidentService.ListIncidents(c.Context(), tenantID, req)
	if err != nil {
		log.Printf("ERROR: incident_handler.listIncidents ListIncidents: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to list incidents",
		})
	}

	// Convert to response format
	var responses []dto.IncidentResponse
	for _, incident := range incidents {
		response := dto.IncidentResponse{
			ID:          incident.ID,
			TenantID:    incident.TenantID,
			Title:       incident.Title,
			Description: incident.Description,
			Category:    incident.Category,
			Status:      incident.Status,
			Criticality: incident.Criticality,
			Source:      incident.Source,
			ReportedBy:  incident.ReportedBy,
			AssignedTo:  incident.AssignedTo,
			DetectedAt:  incident.DetectedAt,
			ResolvedAt:  incident.ResolvedAt,
			ClosedAt:    incident.ClosedAt,
			CreatedAt:   incident.CreatedAt,
			UpdatedAt:   incident.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return c.JSON(fiber.Map{
		"data": responses,
		"pagination": fiber.Map{
			"page":        req.Page,
			"page_size":   req.PageSize,
			"total":       total,
			"total_pages": (total + req.PageSize - 1) / req.PageSize,
		},
	})
}

func (h *IncidentHandler) createIncident(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.CreateIncidentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.createIncident BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.createIncident validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	incident, err := h.incidentService.CreateIncident(c.Context(), tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.createIncident CreateIncident: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create incident",
		})
	}

	response := dto.IncidentResponse{
		ID:          incident.ID,
		TenantID:    incident.TenantID,
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Status:      incident.Status,
		Criticality: incident.Criticality,
		Source:      incident.Source,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
		DetectedAt:  incident.DetectedAt,
		ResolvedAt:  incident.ResolvedAt,
		ClosedAt:    incident.ClosedAt,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}

	return c.Status(201).JSON(response)
}

func (h *IncidentHandler) getIncident(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	id := c.Params("id")

	incident, err := h.incidentService.GetIncident(c.Context(), id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_handler.getIncident GetIncident: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Incident not found",
		})
	}

	response := dto.IncidentResponse{
		ID:          incident.ID,
		TenantID:    incident.TenantID,
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Status:      incident.Status,
		Criticality: incident.Criticality,
		Source:      incident.Source,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
		DetectedAt:  incident.DetectedAt,
		ResolvedAt:  incident.ResolvedAt,
		ClosedAt:    incident.ClosedAt,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}

	return c.JSON(response)
}

func (h *IncidentHandler) updateIncident(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	var req dto.UpdateIncidentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.updateIncident BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.updateIncident validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	incident, err := h.incidentService.UpdateIncident(c.Context(), id, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.updateIncident UpdateIncident: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update incident",
		})
	}

	response := dto.IncidentResponse{
		ID:          incident.ID,
		TenantID:    incident.TenantID,
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Status:      incident.Status,
		Criticality: incident.Criticality,
		Source:      incident.Source,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
		DetectedAt:  incident.DetectedAt,
		ResolvedAt:  incident.ResolvedAt,
		ClosedAt:    incident.ClosedAt,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}

	return c.JSON(response)
}

func (h *IncidentHandler) deleteIncident(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	id := c.Params("id")

	err := h.incidentService.DeleteIncident(c.Context(), id, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_handler.deleteIncident DeleteIncident: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete incident",
		})
	}

	return c.SendStatus(204)
}

func (h *IncidentHandler) updateIncidentStatus(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	id := c.Params("id")

	var req dto.IncidentStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.updateIncidentStatus BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.updateIncidentStatus validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	incident, err := h.incidentService.UpdateIncidentStatus(c.Context(), id, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.updateIncidentStatus UpdateIncidentStatus: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update incident status",
		})
	}

	response := dto.IncidentResponse{
		ID:          incident.ID,
		TenantID:    incident.TenantID,
		Title:       incident.Title,
		Description: incident.Description,
		Category:    incident.Category,
		Status:      incident.Status,
		Criticality: incident.Criticality,
		Source:      incident.Source,
		ReportedBy:  incident.ReportedBy,
		AssignedTo:  incident.AssignedTo,
		DetectedAt:  incident.DetectedAt,
		ResolvedAt:  incident.ResolvedAt,
		ClosedAt:    incident.ClosedAt,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}

	return c.JSON(response)
}

func (h *IncidentHandler) addComment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	incidentID := c.Params("id")

	var req dto.IncidentCommentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.addComment BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.addComment validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	comment, err := h.incidentService.AddComment(c.Context(), incidentID, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.addComment AddComment: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to add comment",
		})
	}

	response := dto.IncidentCommentResponse{
		ID:         comment.ID,
		IncidentID: comment.IncidentID,
		UserID:     comment.UserID,
		Comment:    comment.Comment,
		IsInternal: comment.IsInternal,
		CreatedAt:  comment.CreatedAt,
	}

	return c.Status(201).JSON(response)
}

func (h *IncidentHandler) getComments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	incidentID := c.Params("id")

	comments, err := h.incidentService.GetComments(c.Context(), incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_handler.getComments GetComments: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get comments",
		})
	}

	var responses []dto.IncidentCommentResponse
	for _, comment := range comments {
		response := dto.IncidentCommentResponse{
			ID:         comment.ID,
			IncidentID: comment.IncidentID,
			UserID:     comment.UserID,
			Comment:    comment.Comment,
			IsInternal: comment.IsInternal,
			CreatedAt:  comment.CreatedAt,
		}
		responses = append(responses, response)
	}

	return c.JSON(responses)
}

func (h *IncidentHandler) addAction(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	incidentID := c.Params("id")

	var req dto.IncidentActionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.addAction BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.addAction validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	action, err := h.incidentService.AddAction(c.Context(), incidentID, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.addAction AddAction: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to add action",
		})
	}

	response := dto.IncidentActionResponse{
		ID:          action.ID,
		IncidentID:  action.IncidentID,
		ActionType:  action.ActionType,
		Title:       action.Title,
		Description: action.Description,
		AssignedTo:  action.AssignedTo,
		DueDate:     action.DueDate,
		CompletedAt: action.CompletedAt,
		Status:      action.Status,
		CreatedBy:   action.CreatedBy,
		CreatedAt:   action.CreatedAt,
		UpdatedAt:   action.UpdatedAt,
	}

	return c.Status(201).JSON(response)
}

func (h *IncidentHandler) getActions(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	incidentID := c.Params("id")

	actions, err := h.incidentService.GetActions(c.Context(), incidentID, tenantID)
	if err != nil {
		log.Printf("ERROR: incident_handler.getActions GetActions: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get actions",
		})
	}

	var responses []dto.IncidentActionResponse
	for _, action := range actions {
		response := dto.IncidentActionResponse{
			ID:          action.ID,
			IncidentID:  action.IncidentID,
			ActionType:  action.ActionType,
			Title:       action.Title,
			Description: action.Description,
			AssignedTo:  action.AssignedTo,
			DueDate:     action.DueDate,
			CompletedAt: action.CompletedAt,
			Status:      action.Status,
			CreatedBy:   action.CreatedBy,
			CreatedAt:   action.CreatedAt,
			UpdatedAt:   action.UpdatedAt,
		}
		responses = append(responses, response)
	}

	return c.JSON(responses)
}

func (h *IncidentHandler) updateAction(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	actionID := c.Params("actionId")

	var req dto.IncidentActionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: incident_handler.updateAction BodyParser: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: incident_handler.updateAction validation: %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
	}

	action, err := h.incidentService.UpdateAction(c.Context(), actionID, tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: incident_handler.updateAction UpdateAction: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update action",
		})
	}

	response := dto.IncidentActionResponse{
		ID:          action.ID,
		IncidentID:  action.IncidentID,
		ActionType:  action.ActionType,
		Title:       action.Title,
		Description: action.Description,
		AssignedTo:  action.AssignedTo,
		DueDate:     action.DueDate,
		CompletedAt: action.CompletedAt,
		Status:      action.Status,
		CreatedBy:   action.CreatedBy,
		CreatedAt:   action.CreatedAt,
		UpdatedAt:   action.UpdatedAt,
	}

	return c.JSON(response)
}

func (h *IncidentHandler) deleteAction(c *fiber.Ctx) error {
	// Note: This would need to be implemented in the service layer
	// For now, return not implemented
	return c.Status(501).JSON(fiber.Map{
		"error": "Delete action not implemented",
	})
}

func (h *IncidentHandler) getIncidentMetrics(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	metrics, err := h.incidentService.GetIncidentMetrics(c.Context(), tenantID)
	if err != nil {
		log.Printf("ERROR: incident_handler.getIncidentMetrics GetIncidentMetrics: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get incident metrics",
		})
	}

	response := dto.IncidentMetricsResponse{
		TotalIncidents:  metrics.TotalIncidents,
		OpenIncidents:   metrics.OpenIncidents,
		ClosedIncidents: metrics.ClosedIncidents,
		AverageMTTR:     metrics.AverageMTTR,
		AverageMTTD:     metrics.AverageMTTD,
		ByCriticality:   metrics.ByCriticality,
		ByCategory:      metrics.ByCategory,
		ByStatus:        metrics.ByStatus,
	}

	return c.JSON(response)
}
