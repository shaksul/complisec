package http

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// DocumentHandler handles HTTP requests for documents
type DocumentHandler struct {
	service   *domain.DocumentService
	validator *validator.Validate
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(service *domain.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Register registers document routes
func (h *DocumentHandler) Register(r fiber.Router) {
	// Document CRUD
	r.Get("/documents", RequirePermission("docs.view"), h.listDocuments)
	r.Post("/documents", h.createDocument)
	r.Get("/documents/:id", h.getDocument)
	r.Put("/documents/:id", h.updateDocument)
	r.Delete("/documents/:id", h.deleteDocument)

	// Document versions
	r.Get("/documents/:id/versions", h.listDocumentVersions)
	r.Post("/documents/:id/versions", h.createDocumentVersion)
	r.Get("/documents/versions/:versionId", h.getDocumentVersion)
	r.Get("/documents/versions/:versionId/download", h.downloadDocumentVersion)
	r.Get("/documents/versions/:versionId/preview", h.previewDocumentVersion)
	r.Get("/documents/versions/:versionId/html", h.getDocumentVersionHTML)

	// Document acknowledgments
	r.Get("/documents/:id/acknowledgments", h.listDocumentAcknowledgment)
	r.Post("/documents/:id/acknowledgments", h.createDocumentAcknowledgment)
	r.Put("/documents/acknowledgments/:ackId", h.updateDocumentAcknowledgment)

	// Document quizzes
	r.Get("/documents/:id/quizzes", h.listDocumentQuizzes)
	r.Post("/documents/:id/quizzes", h.createDocumentQuiz)

	// User-specific endpoints
	r.Get("/users/me/pending-acknowledgments", h.getUserPendingAcknowledgment)
}

// listDocuments retrieves documents with filtering
func (h *DocumentHandler) listDocuments(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Parse query parameters
	filters := dto.DocumentFiltersDTO{
		Page:  1,
		Limit: 20,
	}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if docType := c.Query("type"); docType != "" {
		filters.Type = &docType
	}
	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}
	if search := c.Query("search"); search != "" {
		filters.Search = &search
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	documents, err := h.service.ListDocuments(context.Background(), tenantID, filters)
	if err != nil {
		log.Printf("ERROR: listDocuments failed for tenant %s: %v", tenantID, err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": documents})
}

// createDocument creates a new document
func (h *DocumentHandler) createDocument(c *fiber.Ctx) error {
	fmt.Printf("DEBUG: createDocument handler called\n")
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	fmt.Printf("DEBUG: tenantID=%s, userID=%s\n", tenantID, userID)

	var req dto.CreateDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("DEBUG: BodyParser error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	fmt.Printf("DEBUG: Parsed request: %+v\n", req)

	if err := h.validator.Struct(req); err != nil {
		fmt.Printf("DEBUG: Validation failed: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}
	fmt.Printf("DEBUG: Validation passed\n")

	fmt.Printf("DEBUG: Calling service.CreateDocument\n")
	document, err := h.service.CreateDocument(context.Background(), tenantID, userID, req)
	if err != nil {
		fmt.Printf("DEBUG: Service.CreateDocument error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Printf("DEBUG: Service.CreateDocument success: %+v\n", document)

	return c.Status(201).JSON(fiber.Map{"data": document})
}

// getDocument retrieves a document by ID
func (h *DocumentHandler) getDocument(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	document, err := h.service.GetDocument(context.Background(), documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if document == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
	}

	return c.JSON(fiber.Map{"data": document})
}

// updateDocument updates an existing document
func (h *DocumentHandler) updateDocument(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	var req dto.UpdateDocumentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	document, err := h.service.UpdateDocument(context.Background(), documentID, tenantID, userID, req)
	if err != nil {
		if err.Error() == "document not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": document})
}

// deleteDocument deletes a document
func (h *DocumentHandler) deleteDocument(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	err := h.service.DeleteDocument(context.Background(), documentID, tenantID, userID)
	if err != nil {
		if err.Error() == "document not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(204).Send(nil)
}

// listDocumentVersions retrieves versions for a document
func (h *DocumentHandler) listDocumentVersions(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	versions, err := h.service.ListDocumentVersions(context.Background(), documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": versions})
}

// createDocumentVersion creates a new version of a document
func (h *DocumentHandler) createDocumentVersion(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	fmt.Printf("DEBUG: createDocumentVersion handler called\n")
	fmt.Printf("DEBUG: tenantID=%s, userID=%s, documentID=%s\n", tenantID, userID, documentID)

	// Handle multipart form data
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("DEBUG: FormFile error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "File is required"})
	}
	fmt.Printf("DEBUG: File received: %s, size: %d\n", file.Filename, file.Size)

	// Get form values
	title := c.FormValue("title")
	fmt.Printf("DEBUG: Title from form: '%s'\n", title)
	if title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	enableOCRStr := c.FormValue("enableOCR")
	fmt.Printf("DEBUG: EnableOCR from form: '%s'\n", enableOCRStr)
	enableOCR := enableOCRStr == "true"

	// Create DTO
	req := dto.CreateDocumentVersionDTO{
		Title:     title,
		EnableOCR: enableOCR,
	}
	fmt.Printf("DEBUG: Created DTO: %+v\n", req)

	if err := h.validator.Struct(req); err != nil {
		fmt.Printf("DEBUG: Validation failed: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	fmt.Printf("DEBUG: Calling service.CreateDocumentVersionWithFile\n")
	version, err := h.service.CreateDocumentVersionWithFile(context.Background(), documentID, tenantID, userID, req, file, file.Filename)
	if err != nil {
		fmt.Printf("DEBUG: Service error: %v\n", err)
		if err.Error() == "document not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("DEBUG: Success, returning version: %+v\n", version)
	return c.Status(201).JSON(fiber.Map{"data": version})
}

// getDocumentVersion retrieves a specific version of a document
func (h *DocumentHandler) getDocumentVersion(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	versionID := c.Params("versionId")

	version, err := h.service.GetDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if version == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Document version not found"})
	}

	return c.JSON(fiber.Map{"data": version})
}

// listDocumentAcknowledgment retrieves acknowledgments for a document
func (h *DocumentHandler) listDocumentAcknowledgment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	acknowledgments, err := h.service.ListDocumentAcknowledgment(context.Background(), documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": acknowledgments})
}

// createDocumentAcknowledgment creates an acknowledgment for a document
func (h *DocumentHandler) createDocumentAcknowledgment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	var req dto.CreateDocumentAcknowledgmentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	acknowledgment, err := h.service.CreateDocumentAcknowledgment(context.Background(), documentID, tenantID, userID, req)
	if err != nil {
		if err.Error() == "document not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": acknowledgment})
}

// updateDocumentAcknowledgment updates an acknowledgment
func (h *DocumentHandler) updateDocumentAcknowledgment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	ackID := c.Params("ackId")

	var req dto.UpdateDocumentAcknowledgmentDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	acknowledgment, err := h.service.UpdateDocumentAcknowledgment(context.Background(), ackID, tenantID, userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": acknowledgment})
}

// listDocumentQuizzes retrieves quizzes for a document
func (h *DocumentHandler) listDocumentQuizzes(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	documentID := c.Params("id")

	quizzes, err := h.service.ListDocumentQuizzes(context.Background(), documentID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": quizzes})
}

// createDocumentQuiz creates a quiz question for a document
func (h *DocumentHandler) createDocumentQuiz(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)
	documentID := c.Params("id")

	var req dto.CreateDocumentQuizDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	quiz, err := h.service.CreateDocumentQuiz(context.Background(), documentID, tenantID, userID, req)
	if err != nil {
		if err.Error() == "document not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": quiz})
}

// getUserPendingAcknowledgment retrieves pending acknowledgments for the current user
func (h *DocumentHandler) getUserPendingAcknowledgment(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	acknowledgments, err := h.service.GetUserPendingAcknowledgment(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": acknowledgments})
}

// downloadDocumentVersion downloads a document version file
func (h *DocumentHandler) downloadDocumentVersion(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	versionID := c.Params("versionId")

	fmt.Printf("DEBUG: downloadDocumentVersion called for versionID=%s, tenantID=%s\n", versionID, tenantID)

	// Get version info first
	version, err := h.service.GetDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		fmt.Printf("DEBUG: GetDocumentVersion error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if version == nil {
		fmt.Printf("DEBUG: Version not found\n")
		return c.Status(404).JSON(fiber.Map{"error": "Document version not found"})
	}

	fmt.Printf("DEBUG: Version found: %+v\n", version)

	// Get file content from service
	fileContent, err := h.service.DownloadDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		fmt.Printf("DEBUG: DownloadDocumentVersion error: %v\n", err)
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "Document file not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to download document", "details": err.Error()})
	}

	// Determine filename
	filename := fmt.Sprintf("document_v%d", version.VersionNumber)
	if version.MimeType != nil && *version.MimeType != "" {
		switch *version.MimeType {
		case "application/pdf":
			filename += ".pdf"
		case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
			filename += ".docx"
		case "application/msword":
			filename += ".doc"
		case "text/plain":
			filename += ".txt"
		case "image/jpeg":
			filename += ".jpg"
		case "image/png":
			filename += ".png"
		default:
			filename += ".bin"
		}
	} else {
		filename += ".bin"
	}

	fmt.Printf("DEBUG: Setting headers and sending file: %s, size: %d\n", filename, len(fileContent))

	// Set appropriate headers
	mimeType := "application/octet-stream"
	if version.MimeType != nil {
		mimeType = *version.MimeType
	}
	c.Set("Content-Type", mimeType)

	// Check if this is a preview request (from referer or user agent)
	referer := c.Get("Referer")
	isPreview := strings.Contains(referer, "localhost:3000") || c.Query("preview") == "true"

	if isPreview {
		// For preview, use inline disposition to display in browser
		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	} else {
		// For download, use attachment disposition
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	}
	c.Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	return c.Send(fileContent)
}

// previewDocumentVersion generates a preview URL for a document version
func (h *DocumentHandler) previewDocumentVersion(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	versionID := c.Params("versionId")

	fmt.Printf("DEBUG: previewDocumentVersion called for versionID=%s, tenantID=%s\n", versionID, tenantID)

	// Get version info first
	version, err := h.service.GetDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		fmt.Printf("DEBUG: GetDocumentVersion error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if version == nil {
		fmt.Printf("DEBUG: Version not found\n")
		return c.Status(404).JSON(fiber.Map{"error": "Document version not found"})
	}

	fmt.Printf("DEBUG: Version found for preview: %+v\n", version)

	// For now, we'll return the download URL as preview
	// In a real implementation, you might generate a different URL for preview
	previewURL := fmt.Sprintf("/api/documents/versions/%s/download", versionID)

	return c.JSON(fiber.Map{"url": previewURL})
}

// getDocumentVersionHTML converts a document version to HTML for local viewing
func (h *DocumentHandler) getDocumentVersionHTML(c *fiber.Ctx) error {
	fmt.Printf("DEBUG: getDocumentVersionHTML called for versionID=%s\n", c.Params("versionId"))

	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		fmt.Printf("DEBUG: getDocumentVersionHTML - tenant_id not found in locals\n")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	versionID := c.Params("versionId")

	fmt.Printf("DEBUG: getDocumentVersionHTML called for versionID=%s, tenantID=%s\n", versionID, tenantID)

	// Get version info first
	version, err := h.service.GetDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		fmt.Printf("DEBUG: GetDocumentVersion error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if version == nil {
		fmt.Printf("DEBUG: Version not found\n")
		return c.Status(404).JSON(fiber.Map{"error": "Document version not found"})
	}

	fmt.Printf("DEBUG: Version found for HTML conversion: %+v\n", version)

	// Get file content from service
	fileContent, err := h.service.DownloadDocumentVersion(context.Background(), versionID, tenantID)
	if err != nil {
		fmt.Printf("DEBUG: DownloadDocumentVersion error: %v\n", err)
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return c.Status(404).JSON(fiber.Map{"error": "Document file not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to download document", "details": err.Error()})
	}

	// Convert to HTML based on file type
	mimeType := ""
	if version.MimeType != nil {
		mimeType = *version.MimeType
	}
	htmlContent, err := h.service.ConvertDocumentToHTML(context.Background(), fileContent, mimeType)
	if err != nil {
		fmt.Printf("DEBUG: ConvertDocumentToHTML error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to convert document to HTML", "details": err.Error()})
	}

	// Set appropriate headers
	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("Content-Length", fmt.Sprintf("%d", len(htmlContent)))

	return c.Send(htmlContent)
}
