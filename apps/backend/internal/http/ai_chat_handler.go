package http

import (
	"log"

	"risknexus/backend/internal/domain"
	"risknexus/backend/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type AIChatHandler struct {
	service *domain.AIChatService
}

func NewAIChatHandler(service *domain.AIChatService) *AIChatHandler {
	return &AIChatHandler{service: service}
}

func (h *AIChatHandler) Register(r fiber.Router) {
	r.Post("/ai/chat", h.sendChatMessage)
}

// sendChatMessage обрабатывает POST /ai/chat
func (h *AIChatHandler) sendChatMessage(c *fiber.Ctx) error {
	var req dto.SendChatMessageRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("AI Chat: ошибка парсинга запроса: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Валидация
	if req.ProviderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID провайдера обязателен",
		})
	}
	if len(req.Messages) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Сообщения не могут быть пустыми",
		})
	}

	// Отправляем запрос к AI провайдеру
	response, err := h.service.SendChatMessage(c.Context(), req.ProviderID, req.Messages, req.Model)
	if err != nil {
		log.Printf("AI Chat: ошибка отправки сообщения: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}
