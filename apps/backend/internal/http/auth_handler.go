package http

import (
	"context"
	"log"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *domain.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService *domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
	user, roles, err := h.authService.Login(context.Background(), req.Email, req.Password, req.TenantID)
	if err != nil {
		log.Printf("WARN: AuthHandler.login failed: %v", err)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user, roles)
	if err != nil {
		log.Printf("ERROR: AuthHandler.login token generation failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
	}

	// Get user permissions
	permissions, err := h.authService.GetUserPermissions(context.Background(), user.ID)
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

	_, err := h.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	user, roles, err := h.authService.GetUserFromToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user, roles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(response)
}

func (h *AuthHandler) me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.authService.GetUser(context.Background(), userID)
	if err != nil || user == nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	roles, err := h.authService.GetUserRoles(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user roles"})
	}

	permissions, err := h.authService.GetUserPermissions(context.Background(), userID)
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
