package http

import (
	"strconv"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TenantHandler struct {
	tenantService *domain.TenantService
	validator     *validator.Validate
}

func NewTenantHandler(tenantService *domain.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		validator:     validator.New(),
	}
}

func (h *TenantHandler) Register(r fiber.Router) {
	tenants := r.Group("/tenants")
	tenants.Get("/", RequirePermission("tenants.view"), h.listTenants)
	tenants.Post("/", RequirePermission("tenants.create"), h.createTenant)
	tenants.Get("/:id", RequirePermission("tenants.view"), h.getTenant)
	tenants.Put("/:id", RequirePermission("tenants.edit"), h.updateTenant)
	tenants.Delete("/:id", RequirePermission("tenants.delete"), h.deleteTenant)
	tenants.Get("/domain/:domain", RequirePermission("tenants.view"), h.getTenantByDomain)
}

func (h *TenantHandler) listTenants(c *fiber.Ctx) error {
	// Получаем tenant_id и роли из контекста
	tenantID := c.Locals("tenant_id").(string)
	roles := c.Locals("roles").([]string)

	// Параметры пагинации
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Проверяем, является ли пользователь администратором
	isAdmin := false
	for _, role := range roles {
		if role == "Admin" {
			isAdmin = true
			break
		}
	}

	var tenants []dto.TenantResponse
	var total int64

	if isAdmin {
		// Администраторы видят все организации
		allTenants, err := h.tenantService.ListAllTenants(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tenants = allTenants.Data
		total = allTenants.Pagination.Total
	} else {
		// Обычные пользователи видят только свою организацию
		result, err := h.tenantService.GetTenant(c.Context(), tenantID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tenants = []dto.TenantResponse{*result}
		total = 1
	}

	return c.JSON(dto.PaginatedResponse{
		Data: tenants,
		Pagination: dto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		},
	})
}

func (h *TenantHandler) createTenant(c *fiber.Ctx) error {
	var req dto.CreateTenantDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	// Получаем ID пользователя из контекста
	userID := c.Locals("user_id").(string)

	// Создаем организацию
	tenant, err := h.tenantService.CreateTenant(c.Context(), req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(tenant)
}

func (h *TenantHandler) getTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant ID is required"})
	}

	tenant, err := h.tenantService.GetTenant(c.Context(), id)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Tenant not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tenant)
}

func (h *TenantHandler) getTenantByDomain(c *fiber.Ctx) error {
	domain := c.Params("domain")
	if domain == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Domain is required"})
	}

	tenant, err := h.tenantService.GetTenantByDomain(c.Context(), domain)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Tenant not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tenant)
}

func (h *TenantHandler) updateTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant ID is required"})
	}

	var req dto.UpdateTenantDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Валидация
	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	// Получаем ID пользователя из контекста
	userID := c.Locals("user_id").(string)

	// Обновляем организацию
	tenant, err := h.tenantService.UpdateTenant(c.Context(), id, req, userID)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Tenant not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tenant)
}

func (h *TenantHandler) deleteTenant(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Tenant ID is required"})
	}

	// Получаем ID пользователя из контекста
	userID := c.Locals("user_id").(string)

	// Удаляем организацию
	err := h.tenantService.DeleteTenant(c.Context(), id, userID)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Tenant not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(204).Send(nil)
}
