package http

import (
	"log"
	"strconv"
	"time"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RiskHandler struct {
	riskService *domain.RiskService
	validator   *validator.Validate
}

func NewRiskHandler(riskService *domain.RiskService) *RiskHandler {
	return &RiskHandler{
		riskService: riskService,
		validator:   validator.New(),
	}
}

func (h *RiskHandler) Register(r fiber.Router) {
	risks := r.Group("/risks")
	risks.Get("/", RequirePermission("risks.view"), h.listRisks)
	risks.Post("/", RequirePermission("risks.create"), h.createRisk)
	risks.Get("/:id", RequirePermission("risks.view"), h.getRisk)
	risks.Put("/:id", RequirePermission("risks.edit"), h.updateRisk)
	risks.Patch("/:id", RequirePermission("risks.edit"), h.updateRisk)
	risks.Delete("/:id", RequirePermission("risks.delete"), h.deleteRisk)
	risks.Get("/asset/:asset_id", RequirePermission("risks.view"), h.getRisksByAsset)
	risks.Get("/export", RequirePermission("risks.view"), h.exportRisks)

	// Risk related entities endpoints
	riskID := risks.Group("/:risk_id")

	// History
	riskID.Get("/history", RequirePermission("risks.view"), h.getRiskHistory)

	// Comments
	riskID.Get("/comments", RequirePermission("risks.view"), h.getRiskComments)
	riskID.Post("/comments", RequirePermission("risks.comment"), h.addRiskComment)

	// Attachments
	riskID.Get("/attachments", RequirePermission("risks.view"), h.getRiskAttachments)
	riskID.Post("/attachments", RequirePermission("risks.edit"), h.addRiskAttachment)
	riskID.Delete("/attachments/:attachment_id", RequirePermission("risks.edit"), h.deleteRiskAttachment)

	// Controls
	riskID.Get("/controls", RequirePermission("risks.view"), h.getRiskControls)
	riskID.Post("/controls", RequirePermission("risks.edit"), h.addRiskControl)
	riskID.Put("/controls/:control_id", RequirePermission("risks.edit"), h.updateRiskControl)
	riskID.Delete("/controls/:control_id", RequirePermission("risks.edit"), h.deleteRiskControl)

	// Tags
	riskID.Get("/tags", RequirePermission("risks.view"), h.getRiskTags)
	riskID.Post("/tags", RequirePermission("risks.edit"), h.addRiskTag)
	riskID.Delete("/tags/:tag_name", RequirePermission("risks.edit"), h.deleteRiskTag)

	// Documents
	riskID.Get("/documents", RequirePermission("risks.view"), h.getRiskDocuments)
	riskID.Post("/documents/upload", RequirePermission("risks.edit"), h.uploadRiskDocument)
	riskID.Post("/documents/link", RequirePermission("risks.edit"), h.linkRiskDocument)
	riskID.Delete("/documents/:document_id", RequirePermission("risks.edit"), h.deleteRiskDocument)
	riskID.Delete("/documents/:document_id/unlink", RequirePermission("risks.edit"), h.unlinkRiskDocument)
}

// convertToRiskResponse - преобразует Risk в RiskResponse с автоматическим расчетом уровня
func (h *RiskHandler) convertToRiskResponse(risk *repo.Risk) dto.RiskResponse {
	var levelLabel *string
	if risk.Likelihood != nil && risk.Impact != nil {
		_, label := dto.CalculateRiskLevel(*risk.Likelihood, *risk.Impact)
		levelLabel = &label
	}

	return dto.RiskResponse{
		ID:          risk.ID,
		TenantID:    risk.TenantID,
		Title:       risk.Title,
		Description: risk.Description,
		Category:    risk.Category,
		Likelihood:  risk.Likelihood,
		Impact:      risk.Impact,
		Level:       risk.Level,
		Status:      risk.Status,
		OwnerUserID: risk.OwnerUserID,
		AssetID:     risk.AssetID,
		Methodology: risk.Methodology,
		Strategy:    risk.Strategy,
		DueDate:     risk.DueDate,
		CreatedAt:   risk.CreatedAt,
		UpdatedAt:   risk.UpdatedAt,
		LevelLabel:  levelLabel,
	}
}

func (h *RiskHandler) listRisks(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	// Parse filters
	filters := make(map[string]interface{})
	if assetID := c.Query("asset_id"); assetID != "" {
		filters["asset_id"] = assetID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if level := c.Query("level"); level != "" {
		switch level {
		case dto.RiskLevelLabelLow:
			filters["level_range"] = []int{1, 2}
		case dto.RiskLevelLabelMedium:
			filters["level_range"] = []int{3, 4}
		case dto.RiskLevelLabelHigh:
			filters["level_range"] = []int{5, 6}
		case dto.RiskLevelLabelCritical:
			filters["level_range"] = []int{7, 16}
		default:
			if levelValue, err := strconv.Atoi(level); err == nil {
				filters["level_exact"] = levelValue
			}
		}
	}
	if ownerUserID := c.Query("owner_user_id"); ownerUserID != "" {
		filters["owner_user_id"] = ownerUserID
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if methodology := c.Query("methodology"); methodology != "" {
		filters["methodology"] = methodology
	}
	if strategy := c.Query("strategy"); strategy != "" {
		filters["strategy"] = strategy
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	// Parse sorting
	sortField := c.Query("sort_field", "level")
	sortDirection := c.Query("sort_direction", "desc")

	log.Printf("DEBUG: RiskHandler.listRisks tenant=%s user=%s page=%d pageSize=%d filters=%v",
		tenantID, userID, page, pageSize, filters)

	risks, err := h.riskService.ListRisks(c.Context(), tenantID, filters, sortField, sortDirection)
	if err != nil {
		log.Printf("ERROR: RiskHandler.listRisks service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.listRisks returned %d risks", len(risks))

	// Apply pagination manually
	start := (page - 1) * pageSize
	end := start + pageSize
	total := len(risks)

	var paginatedRisks []interface{}
	if start < total {
		if end > total {
			end = total
		}
		paginatedRisks = make([]interface{}, end-start)
		for i, risk := range risks[start:end] {
			response := h.convertToRiskResponse(&risk)
			paginatedRisks[i] = response
		}
	}

	pagination := dto.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      int64(total),
		TotalPages: (total + pageSize - 1) / pageSize,
	}

	return c.JSON(dto.PaginatedResponse{
		Data:       paginatedRisks,
		Pagination: pagination,
	})
}

func (h *RiskHandler) createRisk(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.createRisk tenant=%s user=%s", tenantID, userID)

	var req dto.CreateRiskRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.createRisk invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: RiskHandler.createRisk parsed request: %+v", req)

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.createRisk validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	// Parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			log.Printf("ERROR: RiskHandler.createRisk invalid due_date: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid due_date format. Use YYYY-MM-DD"})
		}
		dueDate = &parsed
	}

	risk, err := h.riskService.CreateRisk(c.Context(), tenantID, req.Title, req.Description, req.Category, req.Likelihood, req.Impact, req.OwnerUserID, req.AssetID, req.Methodology, req.Strategy, dueDate)
	if err != nil {
		log.Printf("ERROR: RiskHandler.createRisk service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.createRisk success id=%s", risk.ID)
	response := h.convertToRiskResponse(risk)
	return c.Status(201).JSON(fiber.Map{"data": response})
}

func (h *RiskHandler) getRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRisk id=%s user=%s", id, userID)

	risk, err := h.riskService.GetRisk(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRisk service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if risk == nil {
		log.Printf("WARN: RiskHandler.getRisk not found id=%s", id)
		return c.Status(404).JSON(fiber.Map{"error": "Risk not found"})
	}

	response := h.convertToRiskResponse(risk)
	return c.JSON(fiber.Map{"data": response})
}

func (h *RiskHandler) updateRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.updateRisk id=%s user=%s", id, userID)

	var req dto.UpdateRiskRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.updateRisk invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.updateRisk validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	// Get current risk first
	currentRisk, err := h.riskService.GetRisk(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: RiskHandler.updateRisk get current risk error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get current risk"})
	}
	if currentRisk == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Risk not found"})
	}

	// Use provided values or keep current ones
	title := currentRisk.Title
	if req.Title != nil {
		title = *req.Title
	}

	likelihood := *currentRisk.Likelihood
	if req.Likelihood != nil {
		likelihood = *req.Likelihood
	}

	impact := *currentRisk.Impact
	if req.Impact != nil {
		impact = *req.Impact
	}

	// Parse due_date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			log.Printf("ERROR: RiskHandler.updateRisk invalid due_date: %v", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid due_date format. Use YYYY-MM-DD"})
		}
		dueDate = &parsed
	}

	err = h.riskService.UpdateRisk(c.Context(), id, title, req.Description, req.Category, likelihood, impact, req.OwnerUserID, req.AssetID, req.Methodology, req.Strategy, dueDate)
	if err != nil {
		log.Printf("ERROR: RiskHandler.updateRisk service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.updateRisk success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Risk updated successfully"})
}

func (h *RiskHandler) deleteRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.deleteRisk id=%s user=%s", id, userID)

	err := h.riskService.DeleteRisk(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: RiskHandler.deleteRisk service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.deleteRisk success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Risk deleted successfully"})
}

func (h *RiskHandler) getRisksByAsset(c *fiber.Ctx) error {
	assetID := c.Params("asset_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRisksByAsset assetID=%s user=%s", assetID, userID)

	// Get risks filtered by asset_id directly
	tenantID := c.Locals("tenant_id").(string)
	filters := map[string]interface{}{
		"asset_id": assetID,
	}
	
	risks, err := h.riskService.ListRisks(c.Context(), tenantID, filters, "created_at", "desc")
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRisksByAsset service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var riskResponses []interface{}
	for _, risk := range risks {
		response := h.convertToRiskResponse(&risk)
		riskResponses = append(riskResponses, response)
	}

	log.Printf("DEBUG: RiskHandler.getRisksByAsset returned %d risks for asset %s", len(riskResponses), assetID)
	return c.JSON(fiber.Map{"data": riskResponses})
}

// Risk History endpoints
func (h *RiskHandler) getRiskHistory(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRiskHistory riskID=%s user=%s", riskID, userID)

	history, err := h.riskService.GetHistory(c.Context(), riskID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskHistory service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var historyResponses []dto.RiskHistoryResponse
	for _, h := range history {
		historyResponses = append(historyResponses, dto.RiskHistoryResponse{
			ID:            h.ID,
			RiskID:        h.RiskID,
			FieldChanged:  h.FieldChanged,
			OldValue:      h.OldValue,
			NewValue:      h.NewValue,
			ChangeReason:  h.ChangeReason,
			ChangedBy:     h.ChangedBy,
			ChangedAt:     h.ChangedAt,
			ChangedByName: h.ChangedByName,
		})
	}

	return c.JSON(fiber.Map{"data": historyResponses})
}

// Risk Comments endpoints
func (h *RiskHandler) getRiskComments(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)
	includeInternal := c.QueryBool("include_internal", false)

	log.Printf("DEBUG: RiskHandler.getRiskComments riskID=%s user=%s includeInternal=%v", riskID, userID, includeInternal)

	comments, err := h.riskService.GetComments(c.Context(), riskID, includeInternal)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskComments service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var commentResponses []dto.RiskCommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, dto.RiskCommentResponse{
			ID:         comment.ID,
			RiskID:     comment.RiskID,
			UserID:     comment.UserID,
			Comment:    comment.Comment,
			IsInternal: comment.IsInternal,
			UserName:   comment.UserName,
			CreatedAt:  comment.CreatedAt,
			UpdatedAt:  comment.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{"data": commentResponses})
}

func (h *RiskHandler) addRiskComment(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.addRiskComment riskID=%s user=%s", riskID, userID)

	var req dto.RiskCommentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskComment invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskComment validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	isInternal := false
	if req.IsInternal != nil {
		isInternal = *req.IsInternal
	}

	err := h.riskService.AddComment(c.Context(), riskID, userID, req.Comment, isInternal)
	if err != nil {
		log.Printf("ERROR: RiskHandler.addRiskComment service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.addRiskComment success riskID=%s", riskID)
	return c.Status(201).JSON(fiber.Map{"message": "Comment added successfully"})
}

// Risk Attachments endpoints
func (h *RiskHandler) getRiskAttachments(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRiskAttachments riskID=%s user=%s", riskID, userID)

	attachments, err := h.riskService.GetAttachments(c.Context(), riskID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskAttachments service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var attachmentResponses []dto.RiskAttachmentResponse
	for _, attachment := range attachments {
		attachmentResponses = append(attachmentResponses, dto.RiskAttachmentResponse{
			ID:             attachment.ID,
			RiskID:         attachment.RiskID,
			FileName:       attachment.FileName,
			FilePath:       attachment.FilePath,
			FileSize:       attachment.FileSize,
			MimeType:       attachment.MimeType,
			FileHash:       attachment.FileHash,
			Description:    attachment.Description,
			UploadedBy:     attachment.UploadedBy,
			UploadedAt:     attachment.UploadedAt,
			UploadedByName: attachment.UploadedByName,
		})
	}

	return c.JSON(fiber.Map{"data": attachmentResponses})
}

func (h *RiskHandler) addRiskAttachment(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.addRiskAttachment riskID=%s user=%s", riskID, userID)

	var req dto.RiskAttachmentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskAttachment invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskAttachment validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.riskService.AddAttachment(c.Context(), riskID, req.FileName, req.FilePath, req.FileSize, req.MimeType, req.FileHash, req.Description, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.addRiskAttachment service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.addRiskAttachment success riskID=%s", riskID)
	return c.Status(201).JSON(fiber.Map{"message": "Attachment added successfully"})
}

func (h *RiskHandler) deleteRiskAttachment(c *fiber.Ctx) error {
	attachmentID := c.Params("attachment_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.deleteRiskAttachment attachmentID=%s user=%s", attachmentID, userID)

	err := h.riskService.DeleteAttachment(c.Context(), attachmentID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.deleteRiskAttachment service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.deleteRiskAttachment success attachmentID=%s", attachmentID)
	return c.Status(200).JSON(fiber.Map{"message": "Attachment deleted successfully"})
}

// Risk Controls endpoints
func (h *RiskHandler) getRiskControls(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRiskControls riskID=%s user=%s", riskID, userID)

	controls, err := h.riskService.GetControls(c.Context(), riskID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskControls service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var controlResponses []dto.RiskControlResponse
	for _, control := range controls {
		controlResponses = append(controlResponses, dto.RiskControlResponse{
			ID:                   control.ID,
			RiskID:               control.RiskID,
			ControlID:            control.ControlID,
			ControlName:          control.ControlName,
			ControlType:          control.ControlType,
			ImplementationStatus: control.ImplementationStatus,
			Effectiveness:        control.Effectiveness,
			Description:          control.Description,
			CreatedBy:            control.CreatedBy,
			CreatedAt:            control.CreatedAt,
			UpdatedAt:            control.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{"data": controlResponses})
}

func (h *RiskHandler) addRiskControl(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.addRiskControl riskID=%s user=%s", riskID, userID)

	var req dto.RiskControlRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskControl invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskControl validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.riskService.AddControl(c.Context(), riskID, req.ControlID, req.ControlName, req.ControlType, req.ImplementationStatus, req.Effectiveness, req.Description, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.addRiskControl service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.addRiskControl success riskID=%s", riskID)
	return c.Status(201).JSON(fiber.Map{"message": "Control added successfully"})
}

func (h *RiskHandler) updateRiskControl(c *fiber.Ctx) error {
	controlID := c.Params("control_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.updateRiskControl controlID=%s user=%s", controlID, userID)

	var req dto.RiskControlRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.updateRiskControl invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.updateRiskControl validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.riskService.UpdateControl(c.Context(), controlID, req.ControlName, req.ControlType, req.ImplementationStatus, req.Effectiveness, req.Description)
	if err != nil {
		log.Printf("ERROR: RiskHandler.updateRiskControl service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.updateRiskControl success controlID=%s", controlID)
	return c.Status(200).JSON(fiber.Map{"message": "Control updated successfully"})
}

func (h *RiskHandler) deleteRiskControl(c *fiber.Ctx) error {
	controlID := c.Params("control_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.deleteRiskControl controlID=%s user=%s", controlID, userID)

	err := h.riskService.DeleteControl(c.Context(), controlID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.deleteRiskControl service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.deleteRiskControl success controlID=%s", controlID)
	return c.Status(200).JSON(fiber.Map{"message": "Control deleted successfully"})
}

// Risk Tags endpoints
func (h *RiskHandler) getRiskTags(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRiskTags riskID=%s user=%s", riskID, userID)

	tags, err := h.riskService.GetTags(c.Context(), riskID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskTags service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert to response format
	var tagResponses []dto.RiskTagResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, dto.RiskTagResponse{
			ID:        tag.ID,
			RiskID:    tag.RiskID,
			TagName:   tag.TagName,
			TagColor:  tag.TagColor,
			CreatedBy: tag.CreatedBy,
			CreatedAt: tag.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{"data": tagResponses})
}

func (h *RiskHandler) addRiskTag(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.addRiskTag riskID=%s user=%s", riskID, userID)

	var req dto.RiskTagRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskTag invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.addRiskTag validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	tagColor := dto.TagColorDefault
	if req.TagColor != nil {
		tagColor = *req.TagColor
	}

	err := h.riskService.AddTag(c.Context(), riskID, req.TagName, tagColor, &userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.addRiskTag service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.addRiskTag success riskID=%s", riskID)
	return c.Status(201).JSON(fiber.Map{"message": "Tag added successfully"})
}

func (h *RiskHandler) deleteRiskTag(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	tagName := c.Params("tag_name")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.deleteRiskTag riskID=%s tagName=%s user=%s", riskID, tagName, userID)

	err := h.riskService.DeleteTagByName(c.Context(), riskID, tagName)
	if err != nil {
		log.Printf("ERROR: RiskHandler.deleteRiskTag service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.deleteRiskTag success riskID=%s tagName=%s", riskID, tagName)
	return c.Status(200).JSON(fiber.Map{"message": "Tag deleted successfully"})
}

// Risk Export endpoint
func (h *RiskHandler) exportRisks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.exportRisks user=%s", userID)

	// TODO: Implement export functionality
	// For now, return a placeholder response
	return c.Status(501).JSON(fiber.Map{"error": "Export functionality not yet implemented"})
}

// Risk Documents endpoints
func (h *RiskHandler) getRiskDocuments(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.getRiskDocuments riskID=%s tenant=%s user=%s", riskID, tenantID, userID)

	documents, err := h.riskService.GetRiskDocuments(c.Context(), riskID, tenantID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.getRiskDocuments service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": documents})
}

func (h *RiskHandler) uploadRiskDocument(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.uploadRiskDocument riskID=%s tenant=%s user=%s", riskID, tenantID, userID)

	// Получаем файл из multipart form
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: RiskHandler.uploadRiskDocument no file: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: RiskHandler.uploadRiskDocument file open error: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// Парсим дополнительные параметры
	req := dto.UploadDocumentDTO{
		Name:        c.FormValue("name", file.Filename),
		Description: stringPtr(c.FormValue("description")),
		FolderID:    stringPtr(c.FormValue("folder_id")),
		Tags:        []string{}, // TODO: парсить теги из формы
		EnableOCR:   c.FormValue("enable_ocr") == "true",
		Metadata:    stringPtr(c.FormValue("metadata")),
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.uploadRiskDocument validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	document, err := h.riskService.UploadRiskDocument(c.Context(), riskID, tenantID, src, file, req, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.uploadRiskDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.uploadRiskDocument success riskID=%s documentID=%s", riskID, document.ID)
	return c.Status(201).JSON(fiber.Map{"data": document})
}

func (h *RiskHandler) linkRiskDocument(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.linkRiskDocument riskID=%s tenant=%s user=%s", riskID, tenantID, userID)

	var req struct {
		DocumentID string `json:"document_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: RiskHandler.linkRiskDocument invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: RiskHandler.linkRiskDocument validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.riskService.LinkExistingDocument(c.Context(), riskID, req.DocumentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.linkRiskDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.linkRiskDocument success riskID=%s documentID=%s", riskID, req.DocumentID)
	return c.Status(201).JSON(fiber.Map{"message": "Document linked successfully"})
}

func (h *RiskHandler) deleteRiskDocument(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	documentID := c.Params("document_id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.deleteRiskDocument riskID=%s documentID=%s tenant=%s user=%s", riskID, documentID, tenantID, userID)

	err := h.riskService.DeleteRiskDocument(c.Context(), riskID, documentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.deleteRiskDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.deleteRiskDocument success riskID=%s documentID=%s", riskID, documentID)
	return c.Status(200).JSON(fiber.Map{"message": "Document deleted successfully"})
}

func (h *RiskHandler) unlinkRiskDocument(c *fiber.Ctx) error {
	riskID := c.Params("risk_id")
	documentID := c.Params("document_id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: RiskHandler.unlinkRiskDocument riskID=%s documentID=%s tenant=%s user=%s", riskID, documentID, tenantID, userID)

	err := h.riskService.UnlinkDocument(c.Context(), riskID, documentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: RiskHandler.unlinkRiskDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: RiskHandler.unlinkRiskDocument success riskID=%s documentID=%s", riskID, documentID)
	return c.Status(200).JSON(fiber.Map{"message": "Document unlinked successfully"})
}

// Helper function
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
