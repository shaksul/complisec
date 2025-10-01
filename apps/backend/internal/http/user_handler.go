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

type UserHandler struct {
	userService *domain.UserService
	roleService *domain.RoleService
	validator   *validator.Validate
}

func NewUserHandler(userService *domain.UserService, roleService *domain.RoleService) *UserHandler {
	return &UserHandler{
		userService: userService,
		roleService: roleService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) Register(r fiber.Router) {
	users := r.Group("/users")
	users.Get("/", RequirePermission("users.view"), h.listUsers)
	users.Post("/", RequirePermission("users.create"), h.createUser)
	users.Get("/:id", RequirePermission("users.view"), h.getUser)
	users.Put("/:id", RequirePermission("users.edit"), h.updateUser)
	users.Delete("/:id", RequirePermission("users.delete"), h.deleteUser)
	users.Get("/:id/roles", RequirePermission("users.view"), h.getUserRoles)
	users.Post("/:id/roles", RequirePermission("users.edit"), h.assignRoleToUser)
	users.Delete("/:id/roles/:role_id", RequirePermission("users.edit"), h.removeRoleFromUser)

	// Role routes are handled by the dedicated RoleHandler

	permissions := r.Group("/permissions")
	permissions.Get("/", RequirePermission("roles.view"), h.listPermissions)
}

func (h *UserHandler) listUsers(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Параметры пагинации
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	fmt.Printf("DEBUG: ListUsers called with tenantID: %s, page: %d, pageSize: %d\n", tenantID, page, pageSize)

	users, total, err := h.userService.ListUsersPaginated(context.Background(), tenantID, page, pageSize)
	if err != nil {
		fmt.Printf("DEBUG: ListUsers error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("DEBUG: ListUsers returned %d users of %d total\n", len(users), total)

	pagination := dto.NewPaginationResponse(page, pageSize, total)

	return c.JSON(dto.PaginatedResponse{
		Data:       users,
		Pagination: pagination,
	})
}

func (h *UserHandler) createUser(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	log.Printf("DEBUG: user_handler.createUser tenant=%s", tenantID)

	// Log raw body
	body := c.Body()
	log.Printf("DEBUG: user_handler.createUser raw body: %s", string(body))

	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: user_handler.createUser invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("DEBUG: user_handler.createUser parsed request: %+v", req)

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: user_handler.createUser validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	log.Printf("DEBUG: user_handler.createUser email=%s roles=%v", req.Email, req.RoleIDs)
	user, err := h.userService.CreateUser(context.Background(), tenantID, req.Email, req.Password, req.FirstName, req.LastName, req.RoleIDs)
	if err != nil {
		log.Printf("ERROR: user_handler.createUser service error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("DEBUG: user_handler.createUser success id=%s", user.ID)
	return c.Status(201).JSON(fiber.Map{"data": user})
}

func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUser(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"data": user})
}

func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	if err := h.userService.UpdateUser(c.Context(), id, req.FirstName, req.LastName, req.IsActive, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "User updated successfully"})
}

func (h *UserHandler) listPermissions(c *fiber.Ctx) error {
	permissions, err := h.userService.GetPermissions(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": permissions})
}

func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	_ = c.Params("id")

	// В реальной системе здесь должна быть логика удаления пользователя
	// Пока что просто возвращаем успех
	return c.JSON(fiber.Map{"data": "User deleted successfully"})
}

func (h *UserHandler) getUserRoles(c *fiber.Ctx) error {
	id := c.Params("id")
	userWithRoles, err := h.userService.GetUserWithRoles(context.Background(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if userWithRoles == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"data": userWithRoles.Roles})
}

func (h *UserHandler) assignRoleToUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req dto.UserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}
	if req.UserID != userID {
		return c.Status(400).JSON(fiber.Map{"error": "User ID mismatch"})
	}

	if err := h.roleService.AssignRoleToUser(c.Context(), req.UserID, req.RoleID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "Role assigned successfully"})
}

func (h *UserHandler) removeRoleFromUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	roleID := c.Params("role_id")

	if err := h.roleService.RemoveRoleFromUser(c.Context(), userID, roleID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "Role removed successfully"})
}
