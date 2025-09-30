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

	log.Printf("DEBUG: AuthHandler.login attempt email=%s", req.Email)
	user, roles, err := h.authService.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		log.Printf("WARN: AuthHandler.login failed: %v", err)
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user.ID, user.TenantID, roles)
	if err != nil {
		log.Printf("ERROR: AuthHandler.login token generation failed: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roles,
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

	token, err := h.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	userID, tenantID, _, err := h.authService.GetUserFromToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Get user roles from database for refresh token
	roles, err := h.authService.GetUserRoles(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get user roles"})
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(userID, tenantID, roles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(response)
}
