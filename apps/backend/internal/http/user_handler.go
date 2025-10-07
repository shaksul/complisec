package http

import (
	"context"
	"fmt"
	"log"
	"time"

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
	users.Get("/catalog", RequirePermission("users.view"), h.getUserCatalog)
	users.Post("/", RequirePermission("users.create"), h.createUser)
	users.Get("/:id", RequirePermission("users.view"), h.getUser)
	users.Get("/:id/detail", RequirePermission("users.view"), h.getUserDetail)
	users.Get("/:id/activity", RequirePermission("users.view"), h.getUserActivity)
	users.Get("/:id/activity/stats", RequirePermission("users.view"), h.getUserActivityStats)
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

	// Преобразуем в DTO
	var userResponses []dto.UserResponse
	for _, user := range users {
		roles, err := h.userService.GetUserRoles(context.Background(), user.ID)
		var roleNames []string
		if err == nil {
			roleNames = roles
		}

		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			Roles:     roleNames,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	pagination := dto.NewPaginationResponse(page, pageSize, total)

	return c.JSON(dto.PaginatedResponse{
		Data:       userResponses,
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

	// Преобразуем в DTO
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.Status(201).JSON(fiber.Map{"data": userResponse})
}

func (h *UserHandler) getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	user, err := h.userService.GetUserByTenant(context.Background(), id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Преобразуем в DTO
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(fiber.Map{"data": userResponse})
}

func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	if err := h.userService.UpdateUserByTenant(c.Context(), id, tenantID, req.FirstName, req.LastName, req.IsActive, req.RoleIDs); err != nil {
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
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	// Проверяем, что пользователь существует в текущем тенанте
	user, err := h.userService.GetUserByTenant(context.Background(), id, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// В реальной системе здесь должна быть логика удаления пользователя
	// Пока что просто возвращаем успех
	return c.JSON(fiber.Map{"data": "User deleted successfully"})
}

func (h *UserHandler) getUserRoles(c *fiber.Ctx) error {
	id := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	userWithRoles, err := h.userService.GetUserWithRolesByTenant(context.Background(), id, tenantID)
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
	tenantID := c.Locals("tenant_id").(string)

	// Проверяем, что пользователь существует в текущем тенанте
	user, err := h.userService.GetUserByTenant(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

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
	tenantID := c.Locals("tenant_id").(string)

	// Проверяем, что пользователь существует в текущем тенанте
	user, err := h.userService.GetUserByTenant(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := h.roleService.RemoveRoleFromUser(c.Context(), userID, roleID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": "Role removed successfully"})
}

func (h *UserHandler) getUserCatalog(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)

	// Параметры запроса
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	search := c.Query("search", "")
	role := c.Query("role", "")
	sortBy := c.Query("sort_by", "created_at")
	sortDir := c.Query("sort_dir", "desc")

	var isActive *bool
	if c.Query("is_active") != "" {
		active := c.QueryBool("is_active")
		isActive = &active
	}

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 20
	}

	users, total, err := h.userService.SearchUsers(context.Background(), tenantID, search, role, isActive, sortBy, sortDir, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Преобразуем в DTO
	var catalogUsers []dto.UserCatalogResponse
	for _, user := range users {
		roles, err := h.userService.GetUserRoles(context.Background(), user.ID)
		var roleNames []string
		if err == nil {
			roleNames = roles
		}

		catalogUsers = append(catalogUsers, dto.UserCatalogResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			Roles:     roleNames,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	pagination := dto.NewPaginationResponse(page, pageSize, total)

	return c.JSON(dto.PaginatedResponse{
		Data:       catalogUsers,
		Pagination: pagination,
	})
}

func (h *UserHandler) getUserDetail(c *fiber.Ctx) error {
	userID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	user, roles, stats, err := h.userService.GetUserDetailByTenant(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := dto.UserDetailResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Stats: dto.UserStats{
			DocumentsCount: stats["documents_count"],
			RisksCount:     stats["risks_count"],
			IncidentsCount: stats["incidents_count"],
			AssetsCount:    stats["assets_count"],
			SessionsCount:  stats["sessions_count"],
			LoginCount:     stats["login_count"],
			ActivityScore:  stats["activity_score"],
		},
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *UserHandler) getUserActivity(c *fiber.Ctx) error {
	userID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	// Проверяем, что пользователь существует в текущем тенанте
	user, err := h.userService.GetUserByTenant(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Параметры пагинации
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 50)

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 50
	}

	activities, total, err := h.userService.GetUserActivity(context.Background(), userID, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Преобразуем в DTO
	var activityResponses []dto.UserActivity
	for _, activity := range activities {
		activityResponses = append(activityResponses, dto.UserActivity{
			ID:          activity.ID,
			UserID:      activity.UserID,
			Action:      activity.Action,
			Description: activity.Description,
			IPAddress:   activity.IPAddress,
			UserAgent:   activity.UserAgent,
			CreatedAt:   activity.CreatedAt,
			Metadata:    activity.Metadata,
		})
	}

	pagination := dto.NewPaginationResponse(page, pageSize, total)

	return c.JSON(dto.PaginatedResponse{
		Data:       activityResponses,
		Pagination: pagination,
	})
}

func (h *UserHandler) getUserActivityStats(c *fiber.Ctx) error {
	userID := c.Params("id")
	tenantID := c.Locals("tenant_id").(string)

	// Проверяем, что пользователь существует в текущем тенанте
	user, err := h.userService.GetUserByTenant(context.Background(), userID, tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	stats, err := h.userService.GetUserActivityStats(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Преобразуем в DTO
	response := dto.UserActivityStats{
		DailyActivity: convertToDailyActivity(stats["daily_activity"]),
		TopActions:    convertToTopActions(stats["top_actions"]),
		LoginHistory:  convertToLoginHistory(stats["login_history"]),
	}

	return c.JSON(fiber.Map{"data": response})
}

// Helper functions to convert interface{} to specific types
func convertToDailyActivity(data interface{}) []dto.DailyActivity {
	if data == nil {
		return []dto.DailyActivity{}
	}

	items, ok := data.([]map[string]interface{})
	if !ok {
		return []dto.DailyActivity{}
	}

	result := make([]dto.DailyActivity, len(items))
	for i, item := range items {
		result[i] = dto.DailyActivity{
			Date:  getString(item["date"]),
			Count: getInt(item["count"]),
		}
	}
	return result
}

func convertToTopActions(data interface{}) []dto.TopAction {
	if data == nil {
		return []dto.TopAction{}
	}

	items, ok := data.([]map[string]interface{})
	if !ok {
		return []dto.TopAction{}
	}

	result := make([]dto.TopAction, len(items))
	for i, item := range items {
		result[i] = dto.TopAction{
			Action: getString(item["action"]),
			Count:  getInt(item["count"]),
		}
	}
	return result
}

func convertToLoginHistory(data interface{}) []dto.LoginHistory {
	if data == nil {
		return []dto.LoginHistory{}
	}

	items, ok := data.([]map[string]interface{})
	if !ok {
		return []dto.LoginHistory{}
	}

	result := make([]dto.LoginHistory, len(items))
	for i, item := range items {
		result[i] = dto.LoginHistory{
			IPAddress: getString(item["ip_address"]),
			UserAgent: getString(item["user_agent"]),
			CreatedAt: getTime(item["created_at"]),
			Success:   getBool(item["success"]),
		}
	}
	return result
}

func getString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func getInt(value interface{}) int {
	if num, ok := value.(int); ok {
		return num
	}
	if num, ok := value.(float64); ok {
		return int(num)
	}
	return 0
}

func getBool(value interface{}) bool {
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}

func getTime(value interface{}) time.Time {
	if t, ok := value.(time.Time); ok {
		return t
	}
	return time.Time{}
}
