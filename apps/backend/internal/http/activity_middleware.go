package http

import (
	"context"
	"log"
	"strings"
	"time"

	"risknexus/backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type ActivityMiddleware struct {
	userService *domain.UserService
}

func NewActivityMiddleware(userService *domain.UserService) *ActivityMiddleware {
	return &ActivityMiddleware{
		userService: userService,
	}
}

// LogUserActivity middleware logs user activities
func (m *ActivityMiddleware) LogUserActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from context (set by auth middleware)
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			// If no user ID, continue without logging
			return c.Next()
		}

		// Get request information
		method := c.Method()
		path := c.Path()
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")

		// Generate action and description based on the request
		action, description := m.generateActionDescription(method, path, c.AllParams())

		// Prepare metadata
		metadata := map[string]interface{}{
			"method": method,
			"path":   path,
			"params": c.AllParams(),
			"query":  c.Queries(),
		}

	// Log activity in background (don't block the request)
	// ВАЖНО: Создаем новый контекст, т.к. c.Context() становится невалидным после завершения обработки запроса
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := m.userService.LogUserActivity(ctx, userID, action, description, ipAddress, userAgent, metadata)
		if err != nil {
			log.Printf("Failed to log user activity: %v", err)
		}
	}()

	return c.Next()
	}
}

// generateActionDescription creates action and description based on HTTP method and path
func (m *ActivityMiddleware) generateActionDescription(method, path string, params map[string]string) (string, string) {
	// Remove query parameters and normalize path
	cleanPath := strings.Split(path, "?")[0]

	// Define action mappings
	switch {
	case strings.HasPrefix(cleanPath, "/api/users"):
		return m.getUserAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/documents"):
		return m.getDocumentAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/risks"):
		return m.getRiskAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/incidents"):
		return m.getIncidentAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/assets"):
		return m.getAssetAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/training"):
		return m.getTrainingAction(method, cleanPath, params)
	case strings.HasPrefix(cleanPath, "/api/auth"):
		return m.getAuthAction(method, cleanPath, params)
	default:
		return m.getGenericAction(method, cleanPath, params)
	}
}

func (m *ActivityMiddleware) getUserAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if strings.Contains(path, "/activity") {
			return "user_activity_view", "Просмотр активности пользователя"
		}
		if strings.Contains(path, "/detail") {
			return "user_detail_view", "Просмотр детальной информации о пользователе"
		}
		if params["id"] != "" {
			return "user_view", "Просмотр пользователя"
		}
		return "users_list_view", "Просмотр списка пользователей"
	case "POST":
		return "user_create", "Создание нового пользователя"
	case "PUT", "PATCH":
		return "user_update", "Обновление информации о пользователе"
	case "DELETE":
		return "user_delete", "Удаление пользователя"
	default:
		return "user_action", "Действие с пользователями"
	}
}

func (m *ActivityMiddleware) getDocumentAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if params["id"] != "" {
			return "document_view", "Просмотр документа"
		}
		return "documents_list_view", "Просмотр списка документов"
	case "POST":
		return "document_create", "Создание нового документа"
	case "PUT", "PATCH":
		return "document_update", "Обновление документа"
	case "DELETE":
		return "document_delete", "Удаление документа"
	default:
		return "document_action", "Действие с документами"
	}
}

func (m *ActivityMiddleware) getRiskAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if params["id"] != "" {
			return "risk_view", "Просмотр риска"
		}
		return "risks_list_view", "Просмотр списка рисков"
	case "POST":
		return "risk_create", "Создание нового риска"
	case "PUT", "PATCH":
		return "risk_update", "Обновление риска"
	case "DELETE":
		return "risk_delete", "Удаление риска"
	default:
		return "risk_action", "Действие с рисками"
	}
}

func (m *ActivityMiddleware) getIncidentAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if params["id"] != "" {
			return "incident_view", "Просмотр инцидента"
		}
		return "incidents_list_view", "Просмотр списка инцидентов"
	case "POST":
		return "incident_create", "Создание нового инцидента"
	case "PUT", "PATCH":
		return "incident_update", "Обновление инцидента"
	case "DELETE":
		return "incident_delete", "Удаление инцидента"
	default:
		return "incident_action", "Действие с инцидентами"
	}
}

func (m *ActivityMiddleware) getAssetAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if params["id"] != "" {
			return "asset_view", "Просмотр актива"
		}
		return "assets_list_view", "Просмотр списка активов"
	case "POST":
		return "asset_create", "Создание нового актива"
	case "PUT", "PATCH":
		return "asset_update", "Обновление актива"
	case "DELETE":
		return "asset_delete", "Удаление актива"
	default:
		return "asset_action", "Действие с активами"
	}
}

func (m *ActivityMiddleware) getTrainingAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		if params["id"] != "" {
			return "training_view", "Просмотр обучения"
		}
		return "training_list_view", "Просмотр списка обучения"
	case "POST":
		return "training_create", "Создание нового обучения"
	case "PUT", "PATCH":
		return "training_update", "Обновление обучения"
	case "DELETE":
		return "training_delete", "Удаление обучения"
	default:
		return "training_action", "Действие с обучением"
	}
}

func (m *ActivityMiddleware) getAuthAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "POST":
		if strings.Contains(path, "login") {
			return "login", "Вход в систему"
		}
		if strings.Contains(path, "logout") {
			return "logout", "Выход из системы"
		}
		if strings.Contains(path, "refresh") {
			return "token_refresh", "Обновление токена"
		}
		return "auth_action", "Действие аутентификации"
	case "GET":
		if strings.Contains(path, "profile") {
			return "profile_view", "Просмотр профиля"
		}
		return "auth_info_view", "Просмотр информации аутентификации"
	default:
		return "auth_action", "Действие аутентификации"
	}
}

func (m *ActivityMiddleware) getGenericAction(method, path string, params map[string]string) (string, string) {
	switch method {
	case "GET":
		return "data_view", "Просмотр данных"
	case "POST":
		return "data_create", "Создание данных"
	case "PUT", "PATCH":
		return "data_update", "Обновление данных"
	case "DELETE":
		return "data_delete", "Удаление данных"
	default:
		return "api_action", "API действие"
	}
}
