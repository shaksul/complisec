package http

import (
	"context"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type AssetHandler struct {
	assetService *domain.AssetService
}

func NewAssetHandler(assetService *domain.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

func (h *AssetHandler) Register(r fiber.Router) {
	assets := r.Group("/assets")
	assets.Get("/", h.listAssets)
	assets.Post("/", RequirePermission("assets.create"), h.createAsset)
	assets.Get("/:id", h.getAsset)
	assets.Put("/:id", RequirePermission("assets.edit"), h.updateAsset)
	assets.Delete("/:id", RequirePermission("assets.delete"), h.deleteAsset)
}

func (h *AssetHandler) listAssets(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	assets, err := h.assetService.ListAssets(context.Background(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": assets})
}

func (h *AssetHandler) createAsset(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	var req struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		InvCode  *string `json:"inv_code"`
		OwnerID  *string `json:"owner_id"`
		Location *string `json:"location"`
		Software *string `json:"software"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	asset, err := h.assetService.CreateAsset(context.Background(), tenantID, req.Name, req.Type, req.InvCode, req.OwnerID, req.Location, req.Software)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": asset})
}

func (h *AssetHandler) getAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	asset, err := h.assetService.GetAsset(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if asset == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.JSON(fiber.Map{"data": asset})
}

func (h *AssetHandler) updateAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		InvCode  *string `json:"inv_code"`
		OwnerID  *string `json:"owner_id"`
		Location *string `json:"location"`
		Software *string `json:"software"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.assetService.UpdateAsset(context.Background(), id, req.Name, req.Type, req.InvCode, req.OwnerID, req.Location, req.Software)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Asset updated successfully"})
}

func (h *AssetHandler) deleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.assetService.DeleteAsset(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": "Asset deleted successfully"})
}
