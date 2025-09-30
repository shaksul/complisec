package http

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type DocumentWorkflowHandler struct {
	documentService *domain.DocumentService
	workflowService *domain.WorkflowService
}

func NewDocumentWorkflowHandler(documentService *domain.DocumentService, workflowService *domain.WorkflowService) *DocumentWorkflowHandler {
	return &DocumentWorkflowHandler{
		documentService: documentService,
		workflowService: workflowService,
	}
}

// UploadDocumentVersion handles file upload for document versions
func (h *DocumentWorkflowHandler) UploadDocumentVersion(c *fiber.Ctx) error {
	fmt.Printf("DEBUG: UploadDocumentVersion called - START\n")
	documentID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	fmt.Printf("DEBUG: UploadDocumentVersion called for document %s\n", documentID)
	fmt.Printf("DEBUG: tenantID=%s, userID=%s\n", tenantID, userID)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Printf("DEBUG: MultipartForm error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid multipart form"})
	}
	fmt.Printf("DEBUG: MultipartForm parsed successfully\n")

	// Get file
	files := form.File["file"]
	fmt.Printf("DEBUG: Found %d files\n", len(files))
	if len(files) == 0 {
		fmt.Printf("DEBUG: No file provided\n")
		return c.Status(400).JSON(fiber.Map{"error": "No file provided"})
	}

	file := files[0]

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".docx" && ext != ".txt" {
		return c.Status(400).JSON(fiber.Map{"error": "Only PDF, DOCX and TXT files are allowed"})
	}

	// Get OCR option
	enableOCR := false
	if ocrStr := form.Value["enableOCR"]; len(ocrStr) > 0 {
		enableOCR, _ = strconv.ParseBool(ocrStr[0])
	}

	// Read file content
	fmt.Printf("DEBUG: Opening file: %s\n", file.Filename)
	src, err := file.Open()
	if err != nil {
		fmt.Printf("DEBUG: Failed to open file: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	fmt.Printf("DEBUG: Reading file content\n")
	fileContent, err := io.ReadAll(src)
	if err != nil {
		fmt.Printf("DEBUG: Failed to read file: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read file"})
	}
	fmt.Printf("DEBUG: File content read successfully, size: %d bytes\n", len(fileContent))

	// Create version DTO
	fmt.Printf("DEBUG: Creating version DTO with enableOCR: %v\n", enableOCR)
	versionDTO := dto.CreateDocumentVersionDTO{
		EnableOCR: enableOCR,
	}

	// Upload version
	fmt.Printf("DEBUG: Calling UploadDocumentVersion service\n")
	version, err := h.documentService.UploadDocumentVersion(c.Context(), tenantID, userID, documentID, fileContent, file.Filename, versionDTO)
	if err != nil {
		fmt.Printf("DEBUG: UploadDocumentVersion service error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Printf("DEBUG: UploadDocumentVersion service success\n")

	return c.Status(201).JSON(version)
}

// SubmitDocumentForApproval submits document for approval workflow
func (h *DocumentWorkflowHandler) SubmitDocumentForApproval(c *fiber.Ctx) error {
	documentID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.SubmitDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Submit for approval
	workflow, err := h.workflowService.SubmitDocumentForApproval(c.Context(), tenantID, userID, documentID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(workflow)
}

// ApproveDocumentStep handles approval/rejection of a workflow step
func (h *DocumentWorkflowHandler) ApproveDocumentStep(c *fiber.Ctx) error {
	documentID := c.Params("id")
	stepID := c.Params("stepId")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.ApprovalActionDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Process approval action
	err := h.workflowService.ProcessApprovalAction(c.Context(), tenantID, userID, documentID, stepID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Action processed successfully"})
}

// PublishDocument publishes an approved document
func (h *DocumentWorkflowHandler) PublishDocument(c *fiber.Ctx) error {
	documentID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	// Publish document
	err := h.documentService.PublishDocument(c.Context(), tenantID, userID, documentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Document published successfully"})
}

// CreateACKCampaign creates an acknowledgment campaign
func (h *DocumentWorkflowHandler) CreateACKCampaign(c *fiber.Ctx) error {
	documentID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	var req dto.CreateACKCampaignDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Create ACK campaign
	campaign, err := h.workflowService.CreateACKCampaign(c.Context(), tenantID, userID, documentID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(campaign)
}

// Register registers workflow routes
func (h *DocumentWorkflowHandler) Register(r fiber.Router) {
	// Document workflow routes
	documents := r.Group("/documents")
	documents.Post("/:id/versions", h.UploadDocumentVersion)
	documents.Post("/:id/submit", h.SubmitDocumentForApproval)
	documents.Post("/:id/publish", h.PublishDocument)
	documents.Post("/:id/ack-campaigns", h.CreateACKCampaign)
	documents.Post("/:id/approval/:stepId", h.ApproveDocumentStep)
}
