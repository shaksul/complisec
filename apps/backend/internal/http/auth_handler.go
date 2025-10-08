package http

import (
	"context"
	"log"
	"time"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *domain.AuthService
	userService *domain.UserService
	validator   *validator.Validate
}

func NewAuthHandler(authService *domain.AuthService, userService *domain.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(r fiber.Router) {
	r.Post("/auth/login", h.login)
	r.Post("/auth/refresh", h.refresh)
	// Protected by middleware; registered under protected group in main
}

func (h *AuthHandler) RegisterProtected(r fiber.Router) {
	// Requires AuthMiddleware
	r.Get("/auth/me", h.me)
}

func (h *AuthHandler) login(c *fiber.Ctx) error {
	log.Println("DEBUG: AuthHandler.login called")
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("ERROR: AuthHandler.login invalid body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("ERROR: AuthHandler.login validation failed: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	log.Printf("DEBUG: AuthHandler.login attempt email=%s tenantID=%s", req.Email, req.TenantID)

	// Get request information for logging
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	user, roles, err := h.authService.Login(c.Context(), req.Email, req.Password, req.TenantID)
	if err != nil {
		log.Printf("WARN: AuthHandler.login failed: %v", err)

		// Log failed login attempt
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			h.userService.LogLoginAttempt(ctx, "", req.TenantID, req.Email, ipAddress, userAgent, false, err.Error())
		}()

		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	// Log successful login attempt
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		h.userService.LogLoginAttempt(ctx, user.ID, req.TenantID, req.Email, ipAddress, userAgent, true, "")
	}()

	accessToken, refreshToken, err := h.authService.GenerateTokens(user, roles)
	if err != nil {
		log.Printf("ERROR: AuthHandler.login token generation failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
	}

	// Get user permissions
	permissions, err := h.authService.GetUserPermissions(c.Context(), user.ID)
	if err != nil {
		log.Printf("ERROR: AuthHandler.login failed to get permissions: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user permissions"})
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Roles:       roles,
			Permissions: permissions,
		},
	}

	return c.JSON(response)
}

func (h *AuthHandler) refresh(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	// Используем доменный метод RefreshToken вместо ручного разбора
	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("ERROR: AuthHandler.refresh failed: %v", err)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Получаем информацию о пользователе из нового access токена для возврата в ответе
	token, err := h.authService.ValidateToken(accessToken)
	if err != nil {
		log.Printf("ERROR: AuthHandler.refresh failed to validate new access token: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to validate new token"})
	}

	userID, _, roles, err := h.authService.ExtractUserFromToken(token)
	if err != nil {
		log.Printf("ERROR: AuthHandler.refresh failed to extract user from token: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to extract user info"})
	}

	// Получаем полную информацию о пользователе
	user, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		log.Printf("ERROR: AuthHandler.refresh failed to get user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user info"})
	}

	// Получаем разрешения пользователя
	permissions, err := h.authService.GetUserPermissions(c.Context(), userID)
	if err != nil {
		log.Printf("ERROR: AuthHandler.refresh failed to get permissions: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user permissions"})
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Roles:       roles,
			Permissions: permissions,
		},
	}

	return c.JSON(response)
}

func (h *AuthHandler) me(c *fiber.Ctx) error {
	// Безопасное получение user_id из контекста
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		log.Printf("ERROR: AuthHandler.me user_id not found in context")
		return c.Status(401).JSON(fiber.Map{"error": "User not authenticated"})
	}

	userID, ok := userIDRaw.(string)
	if !ok {
		log.Printf("ERROR: AuthHandler.me invalid user_id type in context: %T", userIDRaw)
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
	}

	user, err := h.authService.GetUser(c.Context(), userID)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	roles, err := h.authService.GetUserRoles(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user roles"})
	}

	permissions, err := h.authService.GetUserPermissions(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user permissions"})
	}

	return c.JSON(dto.LoginResponse{
		AccessToken:  "",
		RefreshToken: "",
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			IsActive:    user.IsActive,
			Roles:       roles,
			Permissions: permissions,
			CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}
