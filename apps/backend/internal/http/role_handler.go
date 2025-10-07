package http

import (
	"log"
	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService *domain.RoleService
	validator   *validator.Validate
}

func NewRoleHandler(roleService *domain.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validator:   validator.New(),
	}
}

func (h *RoleHandler) Register(r fiber.Router) {
	log.Printf("DEBUG: RoleHandler.Register called")
	roles := r.Group("/roles")
	roles.Get("/", RequirePermission("roles.view"), h.listRoles)
	roles.Post("/", RequirePermission("roles.create"), h.createRole)
	// Удаляем тестовый эндпоинт - он не должен быть в продакшене
	// roles.Get("/test", h.testRoleHandler)
	roles.Get("/:id", RequirePermission("roles.view"), h.getRole)
	roles.Put("/:id", RequirePermission("roles.edit"), h.updateRole)
	roles.Delete("/:id", RequirePermission("roles.delete"), h.deleteRole)
	roles.Get("/:id/users", RequirePermission("roles.view"), h.getRoleUsers)
	log.Printf("DEBUG: RoleHandler.Register completed")
}

func (h *RoleHandler) listRoles(c *fiber.Ctx) error {
	log.Printf("DEBUG: listRoles called - THIS IS WORKING!")
	log.Printf("DEBUG: Request method: %s, path: %s", c.Method(), c.Path())
	tenantID := c.Locals("tenant_id").(string)
	log.Printf("DEBUG: tenantID: %s", tenantID)

	roles, err := h.roleService.ListRoles(c.Context(), tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Преобразуем repo.Role в dto.RoleResponse
	var roleResponses []dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})
	}

	return c.JSON(fiber.Map{"data": roleResponses})
}

func (h *RoleHandler) createRole(c *fiber.Ctx) error {
	log.Printf("DEBUG: role_handler.createRole called - THIS IS WORKING!")
	tenantID := c.Locals("tenant_id").(string)

	// Log raw body
	body := c.Body()
	log.Printf("DEBUG: role_handler.createRole raw body: %s", string(body))

	var req dto.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: role_handler.createRole invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: role_handler.createRole parsed request: %+v", req)

	// Filter out null/empty values from permission_ids
	var validPermissionIDs []string
	for _, id := range req.PermissionIDs {
		if id != "" && id != "null" {
			validPermissionIDs = append(validPermissionIDs, id)
		}
	}
	req.PermissionIDs = validPermissionIDs

	// Validate the request
	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: role_handler.createRole validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	role, err := h.roleService.CreateRole(c.Context(), tenantID, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		log.Printf("ERROR: role_handler.createRole service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Преобразуем repo.Role в dto.RoleResponse
	roleResponse := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	return c.Status(201).JSON(fiber.Map{"data": roleResponse})
}

func (h *RoleHandler) getRole(c *fiber.Ctx) error {
	roleID := c.Params("id")

	role, err := h.roleService.GetRoleWithPermissions(c.Context(), roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if role == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.JSON(fiber.Map{"data": role})
}

func (h *RoleHandler) updateRole(c *fiber.Ctx) error {
	log.Printf("DEBUG: updateRole method called!")
	roleID := c.Params("id")
	log.Printf("DEBUG: updateRole called with roleID: %s", roleID)
	log.Printf("DEBUG: updateRole method: %s, path: %s", c.Method(), c.Path())
	log.Printf("DEBUG: updateRole headers: %v", c.GetReqHeaders())

	var req dto.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: updateRole body parsing failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: updateRole parsed request: %+v", req)
	log.Printf("DEBUG: updateRole permission_ids count: %d", len(req.PermissionIDs))

	err := h.roleService.UpdateRole(c.Context(), roleID, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		log.Printf("ERROR: updateRole service call failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: updateRole completed successfully")
	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

func (h *RoleHandler) deleteRole(c *fiber.Ctx) error {
	roleID := c.Params("id")

	err := h.roleService.DeleteRole(c.Context(), roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role deleted successfully"})
}

func (h *RoleHandler) getRoleUsers(c *fiber.Ctx) error {
	roleID := c.Params("id")

	users, err := h.roleService.GetUsersByRole(c.Context(), roleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": users})
}
