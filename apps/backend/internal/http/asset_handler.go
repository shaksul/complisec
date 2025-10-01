package http

import (
	"context"
	"fmt"
	"log"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	assetService *domain.AssetService
	validator    *validator.Validate
}

func NewAssetHandler(assetService *domain.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
		validator:    validator.New(),
	}
}

func (h *AssetHandler) Register(r fiber.Router) {
	assets := r.Group("/assets")
	assets.Get("/", RequirePermission("assets.view"), h.listAssets)
	assets.Post("/", RequirePermission("assets.create"), h.createAsset)
	assets.Get("/:id", RequirePermission("assets.view"), h.getAsset)
	assets.Put("/:id", RequirePermission("assets.edit"), h.updateAsset)
	assets.Delete("/:id", RequirePermission("assets.delete"), h.deleteAsset)
	assets.Get("/:id/details", RequirePermission("assets.view"), h.getAssetDetails)
	assets.Get("/:id/documents", RequirePermission("assets.view"), h.getAssetDocuments)
	assets.Post("/:id/documents", RequirePermission("assets.edit"), h.addAssetDocument)
	assets.Get("/:id/software", RequirePermission("assets.view"), h.getAssetSoftware)
	assets.Post("/:id/software", RequirePermission("assets.edit"), h.addAssetSoftware)
	assets.Get("/:id/history", RequirePermission("assets.view"), h.getAssetHistory)
	assets.Post("/inventory", RequirePermission("assets.inventory"), h.performInventory)
	assets.Get("/export", RequirePermission("assets.export"), h.exportAssets)
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

	assets, total, err := h.assetService.ListAssetsPaginated(context.Background(), tenantID, page, pageSize, filters)
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

	asset, err := h.assetService.CreateAsset(context.Background(), tenantID, req, userID)
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

	asset, err := h.assetService.GetAsset(context.Background(), id)
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

	asset, err := h.assetService.GetAssetWithDetails(context.Background(), id)
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

	err := h.assetService.UpdateAsset(context.Background(), id, req, userID)
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

	err := h.assetService.DeleteAsset(context.Background(), id, userID)
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

	log.Printf("DEBUG: AssetHandler.getAssetDocuments id=%s user=%s", id, userID)

	documents, err := h.assetService.GetAssetDocuments(context.Background(), id)
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

	err := h.assetService.AddDocument(context.Background(), id, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.addAssetDocument service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.addAssetDocument success id=%s", id)
	return c.Status(201).JSON(fiber.Map{"message": "Document added successfully"})
}

func (h *AssetHandler) getAssetSoftware(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	log.Printf("DEBUG: AssetHandler.getAssetSoftware id=%s user=%s", id, userID)

	software, err := h.assetService.GetAssetSoftware(context.Background(), id)
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

	err := h.assetService.AddSoftware(context.Background(), id, req, userID)
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

	history, err := h.assetService.GetAssetHistory(context.Background(), id)
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

	err := h.assetService.PerformInventory(context.Background(), tenantID, req, userID)
	if err != nil {
		log.Printf("ERROR: AssetHandler.performInventory service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: AssetHandler.performInventory success")
	return c.Status(200).JSON(fiber.Map{"message": "Inventory performed successfully"})
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
	assets, err := h.assetService.ListAssets(context.Background(), tenantID, filters)
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