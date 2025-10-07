package http

import (
	"context"
	"fmt"
	"strings"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *domain.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("DEBUG: AuthMiddleware called for %s %s\n", c.Method(), c.Path())
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Fallback: allow token via query for file preview/download flows
			tokenFromQuery := c.Query("access_token")
			if tokenFromQuery == "" {
				fmt.Printf("DEBUG: No authorization header and no access_token query\n")
				return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
			}

		_, err := authService.ValidateToken(tokenFromQuery)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		user, roles, err := authService.GetUserFromAccessToken(tokenFromQuery)
			if err != nil {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
			}

			// Store user info in context
			c.Locals("user_id", user.ID)
			c.Locals("tenant_id", user.TenantID)
			c.Locals("roles", roles)

			return c.Next()
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
		}

		_, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		user, roles, err := authService.GetUserFromAccessToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Store user info in context
		c.Locals("user_id", user.ID)
		c.Locals("tenant_id", user.TenantID)
		c.Locals("roles", roles)
		fmt.Printf("DEBUG: AuthMiddleware set user_id=%s, tenant_id=%s, roles=%v\n", user.ID, user.TenantID, roles)

		return c.Next()
	}
}

// PermissionChecker интерфейс для проверки прав
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID, permission string) (bool, error)
}

var globalPermissionChecker PermissionChecker

// SetPermissionChecker устанавливает глобальный проверщик прав
func SetPermissionChecker(checker PermissionChecker) {
	globalPermissionChecker = checker
}

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		roles := c.Locals("roles").([]string)

		// Проверяем, есть ли у пользователя роль Admin (имеет все права)
		hasPermission := false
		for _, role := range roles {
			if role == "Admin" {
				hasPermission = true
				break
			}
		}

		// Если не админ, проверяем конкретное право
		if !hasPermission && globalPermissionChecker != nil {
			var err error
			hasPermission, err = globalPermissionChecker.HasPermission(c.Context(), userID, permission)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Permission check failed"})
			}
		}

		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		return c.Next()
	}
}
