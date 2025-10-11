package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"risknexus/backend/internal/dto"
	"risknexus/backend/internal/repo"
)

type AIChatService struct {
	aiRepo *repo.AIRepo
	client *http.Client
}

func NewAIChatService(aiRepo *repo.AIRepo) *AIChatService {
	return &AIChatService{
		aiRepo: aiRepo,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SendChatMessage отправляет сообщения в AI провайдер и возвращает ответ
func (s *AIChatService) SendChatMessage(ctx context.Context, providerID string, messages []dto.ChatMessage, model string) (*dto.SendChatMessageResponse, error) {
	// Получаем провайдера из БД
	provider, err := s.aiRepo.Get(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения провайдера: %w", err)
	}
	if provider == nil {
		return nil, fmt.Errorf("провайдер не найден")
	}
	if !provider.IsActive {
		return nil, fmt.Errorf("провайдер неактивен")
	}

	// Если модель не указана, используем дефолтную
	if model == "" {
		model = "llama3.2"
	}

	// Формируем запрос для OpenWeb UI API
	requestBody := map[string]interface{}{
		"model":    model,
		"messages": messages,
		"stream":   false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга запроса: %w", err)
	}

	// Создаем HTTP запрос к провайдеру
	apiURL := provider.BaseURL + "/api/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if provider.APIKey != nil && *provider.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+*provider.APIKey)
	}

	// Отправляем запрос
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса к провайдеру %s: %w", provider.Name, err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("провайдер %s вернул ошибку %d: %s", provider.Name, resp.StatusCode, string(body))
	}

	// Парсим ответ от провайдера
	var openWebUIResponse struct {
		Choices []struct {
			Message dto.ChatMessage `json:"message"`
		} `json:"choices"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &openWebUIResponse); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	if len(openWebUIResponse.Choices) == 0 {
		return nil, fmt.Errorf("провайдер %s не вернул ответ", provider.Name)
	}

	return &dto.SendChatMessageResponse{
		Message: openWebUIResponse.Choices[0].Message,
		Model:   openWebUIResponse.Model,
	}, nil
}
