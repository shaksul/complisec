package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	assetService domain.AssetServiceInterface
	validator    *validator.Validate
}

func NewAssetHandler(assetService domain.AssetServiceInterface) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
		validator:    validator.New(),
	}
}

func (h *AssetHandler) Register(r fiber.Router) {
	assets := r.Group("/assets")
	assets.Get("/", RequirePermission("assets.view"), h.listAssets)
	assets.Post("/", RequirePermission("assets.create"), h.createAsset)
	assets.Get("/export", RequirePermission("assets.export"), h.exportAssets)
	assets.Post("/inventory", RequirePermission("assets.inventory"), h.performInventory)
	assets.Get("/:id", RequirePermission("assets.view"), h.getAsset)
	assets.Put("/:id", RequirePermission("assets.edit"), h.updateAsset)
	assets.Delete("/:id", RequirePermission("assets.delete"), h.deleteAsset)
	assets.Get("/:id/details", RequirePermission("assets.view"), h.getAssetDetails)
	assets.Get("/:id/documents", RequirePermission("assets.view"), h.getAssetDocuments)
	assets.Post("/:id/documents", RequirePermission("assets.documents:create"), h.addAssetDocument)
	assets.Post("/:id/documents/upload", RequirePermission("assets.documents:create"), h.uploadAssetDocument)
	assets.Post("/:id/documents/link", RequirePermission("assets.documents:link"), h.linkAssetDocument)
	// Document storage endpoints (должен быть ПЕРЕД /documents/:docId)
	assets.Get("/documents/storage", RequirePermission("assets.view"), h.getDocumentStorage)
	assets.Delete("/documents/:docId", RequirePermission("assets.edit"), h.deleteAssetDocument)
	assets.Get("/documents/:docId", RequirePermission("assets.view"), h.getAssetDocument)
	assets.Get("/documents/:docId/download", RequirePermission("assets.view"), h.downloadAssetDocument)
	// New centralized document endpoints
	assets.Post("/:id/documents/unlink", RequirePermission("assets.edit"), h.unlinkAssetDocument)
	assets.Get("/:id/software", RequirePermission("assets.view"), h.getAssetSoftware)
	assets.Post("/:id/software", RequirePermission("assets.edit"), h.addAssetSoftware)
	assets.Get("/:id/history", RequirePermission("assets.view"), h.getAssetHistory)
	assets.Get("/:id/history/filtered", RequirePermission("assets.view"), h.getAssetHistoryWithFilters)
	assets.Get("/:id/risks", RequirePermission("assets.view"), h.getAssetRisks)
	assets.Get("/:id/incidents", RequirePermission("assets.view"), h.getAssetIncidents)
	assets.Get("/:id/can-add-risk", RequirePermission("assets.view"), h.canAddRisk)
	assets.Get("/:id/can-add-incident", RequirePermission("assets.view"), h.canAddIncident)
	assets.Get("/inventory/without-owner", RequirePermission("assets.inventory"), h.getAssetsWithoutOwner)
	assets.Get("/inventory/without-passport", RequirePermission("assets.inventory"), h.getAssetsWithoutPassport)
	assets.Get("/inventory/without-criticality", RequirePermission("assets.inventory"), h.getAssetsWithoutCriticality)
	assets.Post("/bulk/update-status", RequirePermission("assets.edit"), h.bulkUpdateStatus)
	assets.Post("/bulk/update-owner", RequirePermission("assets.edit"), h.bulkUpdateOwner)
}

func (h *AssetHandler) listAssets(c *fiber.Ctx) error {
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
	if assetType := c.Query("type"); assetType != "" {
		filters["type"] = assetType
	}
	if class := c.Query("class"); class != "" {
		filters["class"] = class
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if criticality := c.Query("criticality"); criticality != "" {
		filters["criticality"] = criticality
	}
	if ownerID := c.Query("owner_id"); ownerID != "" {
		filters["owner_id"] = ownerID
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	log.Printf("DEBUG: AssetHandler.listAssets tenant=%s user=%s page=%d pageSize=%d filters=%v",
		tenantID, userID, page, pageSize, filters)

	assets, total, err := h.assetService.ListAssetsPaginated(c.Context(), tenantID, page, pageSize, filters)
	if err != nil {
		log.Printf("ERROR: AssetHandler.listAssets service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.listAssets returned %d assets of %d total", len(assets), total)

	pagination := dto.NewPaginationResponse(page, pageSize, total)

	return c.JSON(dto.PaginatedResponse{
		Data:       assets,
		Pagination: pagination,
	})
}

func (h *AssetHandler) createAsset(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.createAsset tenant=%s user=%s", tenantID, userID)

	var req dto.CreateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.createAsset invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: AssetHandler.createAsset parsed request: %+v", req)

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.createAsset validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	asset, err := h.assetService.CreateAsset(c.Context(), tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.createAsset service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.createAsset success id=%s", asset.ID)
	return c.Status(201).JSON(fiber.Map{"data": asset})
}

func (h *AssetHandler) getAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAsset id=%s user=%s", id, userID)

	asset, err := h.assetService.GetAsset(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAsset service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if asset == nil {
		log.Printf("WARN: AssetHandler.getAsset not found id=%s", id)
		return c.Status(404).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.JSON(fiber.Map{"data": asset})
}

func (h *AssetHandler) getAssetDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetDetails id=%s user=%s", id, userID)

	asset, err := h.assetService.GetAssetWithDetails(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetDetails service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if asset == nil {
		log.Printf("WARN: AssetHandler.getAssetDetails not found id=%s", id)
		return c.Status(404).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.JSON(fiber.Map{"data": asset})
}

func (h *AssetHandler) updateAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.updateAsset id=%s user=%s", id, userID)

	var req dto.UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.updateAsset invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.updateAsset validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.UpdateAsset(c.Context(), id, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.updateAsset service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.updateAsset success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Asset updated successfully"})
}

func (h *AssetHandler) deleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.deleteAsset id=%s user=%s", id, userID)

	err := h.assetService.DeleteAsset(c.Context(), id, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.deleteAsset service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.deleteAsset success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Asset deleted successfully"})
}

func (h *AssetHandler) getAssetDocuments(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetDocuments id=%s user=%s", id, userID)

	documents, err := h.assetService.GetAssetDocumentsFromStorage(c.Context(), id, tenantID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetDocuments service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": documents})
}

func (h *AssetHandler) addAssetDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.addAssetDocument id=%s user=%s", id, userID)

	var req dto.AssetDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.addAssetDocument invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.addAssetDocument validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.AddDocument(c.Context(), id, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.addAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.addAssetDocument success id=%s", id)
	return c.Status(201).JSON(fiber.Map{"message": "Document added successfully"})
}

func (h *AssetHandler) uploadAssetDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.uploadAssetDocument id=%s user=%s", id, userID)

	// Parse multipart form
	_, err := c.MultipartForm()
	if err != nil {
		log.Printf("ERROR: AssetHandler.uploadAssetDocument invalid form: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid multipart form"})
	}

	// Get document type
	documentType := c.FormValue("document_type")
	if documentType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Document type is required"})
	}

	// Get title (optional)
	title := c.FormValue("title")
	if title == "" {
		title = documentType
	}

	// Get description (optional)
	_ = c.FormValue("description")

	// Get tags (optional)
	tags := []string{"#активы", fmt.Sprintf("#%s", documentType)}
	if tagsStr := c.FormValue("tags"); tagsStr != "" {
		// Parse tags from JSON array
		if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
			log.Printf("WARN: AssetHandler.uploadAssetDocument invalid tags format: %v", err)
		}
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: AssetHandler.uploadAssetDocument no file: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "File is required"})
	}

	// Validate file size (50MB limit)
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if file.Size > maxFileSize {
		return c.Status(400).JSON(fiber.Map{"error": "File size exceeds 50MB limit"})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf":    true,
		"image/jpeg":         true,
		"image/png":          true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
		"text/plain": true, // Add support for TXT files
	}

	// Get content type from file header
	fileHeader, err := file.Open()
	if err != nil {
		log.Printf("ERROR: AssetHandler.uploadAssetDocument file open: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Cannot open file"})
	}
	defer fileHeader.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		log.Printf("ERROR: AssetHandler.uploadAssetDocument file read: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Cannot read file"})
	}

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		return c.Status(400).JSON(fiber.Map{"error": "Unsupported file type"})
	}

	// Reset file pointer
	fileHeader.Seek(0, 0)

	// Upload document to centralized storage
	document, err := h.assetService.UploadAssetDocument(c.Context(), id, tenantID, fileHeader, file, dto.AssetDocumentUploadRequest{
		DocumentType: documentType,
		Title:        title,
	}, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.uploadAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.uploadAssetDocument success id=%s", id)
	return c.Status(201).JSON(fiber.Map{"data": document})
}

func (h *AssetHandler) linkAssetDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.linkAssetDocument id=%s user=%s", id, userID)

	var req struct {
		DocumentID string `json:"document_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.linkAssetDocument invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.linkAssetDocument validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.LinkExistingDocumentToAsset(c.Context(), id, req.DocumentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.linkAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.linkAssetDocument success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Document linked successfully"})
}

func (h *AssetHandler) downloadAssetDocument(c *fiber.Ctx) error {
	docID := c.Params("docId")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.downloadAssetDocument docID=%s user=%s", docID, userID)

	filePath, fileName, err := h.assetService.GetDocumentDownloadPath(c.Context(), docID, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.downloadAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Download(filePath, fileName)
}

func (h *AssetHandler) getDocumentStorage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.getDocumentStorage user=%s", userID)

	// Parse query parameters
	query := c.Query("query", "")
	docType := c.Query("type", "")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 25)

	req := dto.DocumentStorageRequest{
		Query:    query,
		Type:     docType,
		Page:     page,
		PageSize: pageSize,
	}

	documents, total, err := h.assetService.GetDocumentStorage(c.Context(), tenantID, req)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getDocumentStorage service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": documents,
		"pagination": fiber.Map{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *AssetHandler) getAssetSoftware(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetSoftware id=%s user=%s", id, userID)

	software, err := h.assetService.GetAssetSoftware(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetSoftware service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": software})
}

func (h *AssetHandler) addAssetSoftware(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.addAssetSoftware id=%s user=%s", id, userID)

	var req dto.AssetSoftwareRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.addAssetSoftware invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.addAssetSoftware validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.AddSoftware(c.Context(), id, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.addAssetSoftware service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.addAssetSoftware success id=%s", id)
	return c.Status(201).JSON(fiber.Map{"message": "Software added successfully"})
}

func (h *AssetHandler) getAssetHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetHistory id=%s user=%s", id, userID)

	history, err := h.assetService.GetAssetHistory(c.Context(), id)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetHistory service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": history})
}

func (h *AssetHandler) performInventory(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.performInventory tenant=%s user=%s", tenantID, userID)

	var req dto.AssetInventoryRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.performInventory invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.performInventory validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.PerformInventory(c.Context(), tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.performInventory service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.performInventory success")
	return c.Status(200).JSON(fiber.Map{"message": "Inventory performed successfully"})
}

func (h *AssetHandler) deleteAssetDocument(c *fiber.Ctx) error {
	documentID := c.Params("docId")
	assetID := c.Query("asset_id") // Получаем asset_id из query параметров
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.deleteAssetDocument docID=%s assetID=%s user=%s", documentID, assetID, userID)

	if assetID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "asset_id query parameter is required"})
	}

	err := h.assetService.DeleteAssetDocument(c.Context(), assetID, documentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.deleteAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.deleteAssetDocument success docID=%s", documentID)
	return c.Status(200).JSON(fiber.Map{"message": "Document deleted successfully"})
}

func (h *AssetHandler) getAssetDocument(c *fiber.Ctx) error {
	documentID := c.Params("docId")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetDocument docID=%s user=%s", documentID, userID)

	document, err := h.assetService.GetDocumentByID(c.Context(), documentID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if document == nil {
		log.Printf("WARN: AssetHandler.getAssetDocument not found docID=%s", documentID)
		return c.Status(404).JSON(fiber.Map{"error": "Document not found"})
	}

	return c.JSON(fiber.Map{"data": document})
}

func (h *AssetHandler) getAssetHistoryWithFilters(c *fiber.Ctx) error {
	assetID := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetHistoryWithFilters assetID=%s user=%s", assetID, userID)

	// Parse filters
	filters := make(map[string]interface{})
	if changedBy := c.Query("changed_by"); changedBy != "" {
		filters["changed_by"] = changedBy
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	history, err := h.assetService.GetAssetHistoryWithFilters(c.Context(), assetID, filters)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetHistoryWithFilters service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": history})
}

func (h *AssetHandler) getAssetRisks(c *fiber.Ctx) error {
	assetID := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetRisks assetID=%s user=%s", assetID, userID)

	risks, err := h.assetService.GetAssetRisks(c.Context(), assetID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetRisks service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": risks})
}

func (h *AssetHandler) getAssetIncidents(c *fiber.Ctx) error {
	assetID := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetIncidents assetID=%s user=%s", assetID, userID)

	incidents, err := h.assetService.GetAssetIncidents(c.Context(), assetID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetIncidents service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": incidents})
}

func (h *AssetHandler) canAddRisk(c *fiber.Ctx) error {
	assetID := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.canAddRisk assetID=%s user=%s", assetID, userID)

	err := h.assetService.CanAddRisk(c.Context(), assetID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error(), "can_add": false})
	}

	return c.JSON(fiber.Map{"can_add": true})
}

func (h *AssetHandler) canAddIncident(c *fiber.Ctx) error {
	assetID := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.canAddIncident assetID=%s user=%s", assetID, userID)

	err := h.assetService.CanAddIncident(c.Context(), assetID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error(), "can_add": false})
	}

	return c.JSON(fiber.Map{"can_add": true})
}

func (h *AssetHandler) getAssetsWithoutOwner(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetsWithoutOwner tenant=%s user=%s", tenantID, userID)

	assets, err := h.assetService.GetAssetsWithoutOwner(c.Context(), tenantID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetsWithoutOwner service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": assets})
}

func (h *AssetHandler) getAssetsWithoutPassport(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetsWithoutPassport tenant=%s user=%s", tenantID, userID)

	assets, err := h.assetService.GetAssetsWithoutPassport(c.Context(), tenantID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetsWithoutPassport service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": assets})
}

func (h *AssetHandler) getAssetsWithoutCriticality(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetsWithoutCriticality tenant=%s user=%s", tenantID, userID)

	assets, err := h.assetService.GetAssetsWithoutCriticality(c.Context(), tenantID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.getAssetsWithoutCriticality service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": assets})
}

func (h *AssetHandler) bulkUpdateStatus(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.bulkUpdateStatus tenant=%s user=%s", tenantID, userID)

	var req dto.BulkUpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateStatus invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateStatus validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.BulkUpdateStatus(c.Context(), req.AssetIDs, req.Status, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateStatus service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.bulkUpdateStatus success")
	return c.Status(200).JSON(fiber.Map{"message": "Assets updated successfully"})
}

func (h *AssetHandler) bulkUpdateOwner(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.bulkUpdateOwner tenant=%s user=%s", tenantID, userID)

	var req dto.BulkUpdateOwnerRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateOwner invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateOwner validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.BulkUpdateOwner(c.Context(), req.AssetIDs, req.OwnerID, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.bulkUpdateOwner service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.bulkUpdateOwner success")
	return c.Status(200).JSON(fiber.Map{"message": "Assets updated successfully"})
}

func (h *AssetHandler) exportAssets(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.exportAssets tenant=%s user=%s", tenantID, userID)

	// Parse filters (same as listAssets)
	filters := make(map[string]interface{})
	if assetType := c.Query("type"); assetType != "" {
		filters["type"] = assetType
	}
	if class := c.Query("class"); class != "" {
		filters["class"] = class
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if criticality := c.Query("criticality"); criticality != "" {
		filters["criticality"] = criticality
	}
	if ownerID := c.Query("owner_id"); ownerID != "" {
		filters["owner_id"] = ownerID
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	// Get all assets (no pagination for export)
	assets, err := h.assetService.ListAssets(c.Context(), tenantID, filters)
	if err != nil {
		log.Printf("ERROR: AssetHandler.exportAssets service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Set headers for CSV download
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=assets.csv")

	// Simple CSV export (in production, use a proper CSV library)
	csv := "ID,Inventory Number,Name,Type,Class,Owner,Location,Criticality,Confidentiality,Integrity,Availability,Status,Created At\n"
	for _, asset := range assets {
		owner := ""
		if asset.OwnerID != nil {
			owner = *asset.OwnerID
		}
		location := ""
		if asset.Location != nil {
			location = *asset.Location
		}
		csv += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			asset.ID, asset.InventoryNumber, asset.Name, asset.Type, asset.Class,
			owner, location, asset.Criticality, asset.Confidentiality,
			asset.Integrity, asset.Availability, asset.Status, asset.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return c.SendString(csv)
}

func (h *AssetHandler) unlinkAssetDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	tenantID := c.Locals("tenant_id").(string)

	log.Printf("DEBUG: AssetHandler.unlinkAssetDocument id=%s user=%s", id, userID)

	var req struct {
		DocumentID string `json:"document_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AssetHandler.unlinkAssetDocument invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AssetHandler.unlinkAssetDocument validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	err := h.assetService.UnlinkDocumentFromAsset(c.Context(), id, req.DocumentID, tenantID, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.unlinkAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.unlinkAssetDocument success id=%s", id)
	return c.Status(200).JSON(fiber.Map{"message": "Document unlinked successfully"})
}
