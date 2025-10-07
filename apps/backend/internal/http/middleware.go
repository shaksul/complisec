package http

import (
	"context"
	"log"
	"strings"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *domain.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("DEBUG: AuthMiddleware called for %s %s", c.Method(), c.Path())

		// Обработка ошибок
		defer func() {
			if r := recover(); r != nil {
				log.Printf("ERROR: AuthMiddleware panic: %v", r)
			}
		}()

		var tokenString string

		// Проверяем Authorization header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
			}
		} else {
			// Fallback: проверяем query parameter для file preview/download flows
			tokenString = c.Query("access_token")
			if tokenString == "" {
				log.Printf("DEBUG: No authorization header and no access_token query")
				return c.Status(401).JSON(fiber.Map{"error": "Authorization header required"})
			}
		}

		// Валидация токена
		_, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Получение информации о пользователе из токена
		user, roles, err := authService.GetUserFromAccessToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Сохранение информации о пользователе в контексте
		c.Locals("user_id", user.ID)
		c.Locals("tenant_id", user.TenantID)
		c.Locals("roles", roles)

		log.Printf("DEBUG: AuthMiddleware authenticated user_id=%s, tenant_id=%s, roles=%v", user.ID, user.TenantID, roles)

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
		// Безопасное получение user_id из контекста
		userIDRaw := c.Locals("user_id")
		if userIDRaw == nil {
			log.Printf("ERROR: RequirePermission user_id not found in context")
			return c.Status(401).JSON(fiber.Map{"error": "User not authenticated"})
		}

		userID, ok := userIDRaw.(string)
		if !ok {
			log.Printf("ERROR: RequirePermission invalid user_id type in context: %T", userIDRaw)
			return c.Status(401).JSON(fiber.Map{"error": "Invalid user context"})
		}

		// Безопасно получаем роли, проверяя на nil
		var roles []string
		if rolesRaw := c.Locals("roles"); rolesRaw != nil {
			if rolesSlice, ok := rolesRaw.([]string); ok {
				roles = rolesSlice
			}
		}

		log.Printf("DEBUG: RequirePermission user_id=%s roles=%v permission=%s", userID, roles, permission)

		// Проверяем, есть ли у пользователя роль Admin (имеет все права)
		hasPermission := false
		for _, role := range roles {
			if role == "Admin" {
				hasPermission = true
				log.Printf("DEBUG: RequirePermission user has Admin role, granting access")
				break
			}
		}

		// Если не админ, проверяем конкретное право
		if !hasPermission && globalPermissionChecker != nil {
			var err error
			hasPermission, err = globalPermissionChecker.HasPermission(c.Context(), userID, permission)
			if err != nil {
				log.Printf("ERROR: RequirePermission permission check failed: %v", err)
				return c.Status(500).JSON(fiber.Map{"error": "Permission check failed"})
			}
		}

		if !hasPermission {
			log.Printf("WARN: RequirePermission access denied user_id=%s roles=%v permission=%s", userID, roles, permission)
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}

		log.Printf("DEBUG: RequirePermission access granted user_id=%s permission=%s", userID, permission)
		return c.Next()
	}
}
