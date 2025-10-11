package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"risknexus/backend/internal/repo"
	"strings"
)

type AIService struct {
	repo       *repo.AIRepo
	ragService *RAGService
}

func NewAIService(r *repo.AIRepo) *AIService {
	return &AIService{repo: r}
}

func (s *AIService) SetRAGService(ragService *RAGService) {
	s.ragService = ragService
}

func (s *AIService) List(ctx context.Context, tenantID string) ([]repo.AIProvider, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *AIService) Get(ctx context.Context, id string) (*repo.AIProvider, error) {
	return s.repo.Get(ctx, id)
}

func (s *AIService) Create(ctx context.Context, p repo.AIProvider) error {
	return s.repo.Create(ctx, p)
}

func (s *AIService) Update(ctx context.Context, p repo.AIProvider) error {
	return s.repo.Update(ctx, p)
}

func (s *AIService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *AIService) Query(ctx context.Context, provider repo.AIProvider, role, input string, contextData any) (string, error) {
	// Определяем формат на основе URL
	var payload map[string]any

	if strings.Contains(provider.BaseURL, "ollama") {
		// Ollama формат
		payload = map[string]any{
			"model":  provider.DefaultModel,
			"prompt": fmt.Sprintf("Ты - корпоративный ассистент по информационной безопасности. Роль: %s\n\n%s", role, input),
			"stream": false,
		}
	} else {
		// OpenAI-совместимый формат
		payload = map[string]any{
			"model": provider.DefaultModel,
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": fmt.Sprintf("Ты - корпоративный ассистент по информационной безопасности. Роль: %s", role),
				},
				{
					"role":    "user",
					"content": input,
				},
			},
			"max_tokens":  2000,
			"temperature": 0.7,
		}
	}
	body, _ := json.Marshal(payload)

	log.Printf("DEBUG: AI Query to %s - payload: %s", provider.BaseURL, string(body)[:200])

	req, _ := http.NewRequestWithContext(ctx, "POST", provider.BaseURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if provider.APIKey != nil {
		req.Header.Set("Authorization", "Bearer "+*provider.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR: AI provider HTTP error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("DEBUG: AI provider response status: %d", resp.StatusCode)

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("ERROR: AI provider JSON decode error: %v", err)
		return "", err
	}

	log.Printf("DEBUG: AI provider response data keys: %v", getKeys(data))

	// Ollama формат ответа
	if strings.Contains(provider.BaseURL, "ollama") {
		if response, ok := data["response"].(string); ok {
			return response, nil
		}
	} else {
		// OpenAI-совместимый формат ответа
		if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := choice["message"].(map[string]interface{}); ok {
					if content, ok := message["content"].(string); ok {
						return content, nil
					}
				}
			}
		}
	}

	// Fallback для других форматов
	if out, ok := data["output"].(string); ok {
		return out, nil
	}
	if answer, ok := data["answer"].(string); ok {
		return answer, nil
	}
	if response, ok := data["response"].(string); ok {
		return response, nil
	}
	if message, ok := data["message"].(string); ok {
		return message, nil
	}

	log.Printf("WARN: AI provider returned unexpected format: %+v", data)
	return "(no output)", nil
}

func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// QueryWithRAG - запрос с контекстом из RAG
func (s *AIService) QueryWithRAG(ctx context.Context, tenantID, userID, providerID, role, query string, useRAG bool) (*AIQueryWithRAGResult, error) {
	var ragContext string
	var sources []RAGSource

	// 1. Если RAG включен, получаем контекст
	if useRAG && s.ragService != nil {
		ragResult, err := s.ragService.Query(ctx, tenantID, userID, query, true, 5)
		if err == nil && ragResult != nil {
			sources = ragResult.Sources
			// Формируем контекст из найденных фрагментов
			for i, src := range sources {
				ragContext += "\n\n[Источник " + string(rune(i+1)) + ": " + src.Title + "]\n" + src.Content
			}
		}
	}

	// 2. Получаем AI провайдера
	provider, err := s.Get(ctx, providerID)
	if err != nil || provider == nil {
		return &AIQueryWithRAGResult{
			Answer:  "Ошибка: AI провайдер не найден",
			Sources: sources,
		}, err
	}

	// 3. Формируем промпт с контекстом
	var finalInput string
	if ragContext != "" {
		finalInput = `Ты - корпоративный ассистент по информационной безопасности.

КОНТЕКСТ ИЗ ДОКУМЕНТОВ КОМПАНИИ:
` + ragContext + `

ВОПРОС ПОЛЬЗОВАТЕЛЯ:
` + query + `

ИНСТРУКЦИИ:
1. Ответь на вопрос подробно и структурированно
2. Используй ТОЛЬКО информацию из предоставленного контекста
3. Если в контексте нет нужной информации - честно об этом скажи
4. Не придумывай информацию, которой нет в контексте
5. Используй маркированные списки для структурирования ответа
6. Обязательно укажи из какого источника взята информация

Твой ответ:`
	} else {
		finalInput = query
	}

	// 4. Отправляем запрос в AI провайдер
	log.Printf("DEBUG: Sending to AI provider %s: role=%s, input length=%d", provider.Name, role, len(finalInput))
	answer, err := s.Query(ctx, *provider, role, finalInput, nil)
	if err != nil {
		log.Printf("ERROR: AI provider error: %v", err)
		return &AIQueryWithRAGResult{
			Answer:  "Ошибка при обращении к AI: " + err.Error(),
			Sources: sources,
		}, err
	}

	log.Printf("DEBUG: AI provider response length=%d", len(answer))

	// Если ответ пустой, но есть источники - сформируем базовый ответ
	if answer == "" && len(sources) > 0 {
		answer = "К сожалению, AI провайдер не вернул ответ. Найдены следующие источники:"
		for i, src := range sources {
			answer += fmt.Sprintf("\n\n%d. %s (релевантность: %.2f)", i+1, src.Title, src.Score)
		}
	}

	return &AIQueryWithRAGResult{
		Answer:  answer,
		Sources: sources,
	}, nil
}

type AIQueryWithRAGResult struct {
	Answer  string      `json:"answer"`
	Sources []RAGSource `json:"sources"`
}
