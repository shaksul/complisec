package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"risknexus/backend/internal/repo"
)

type AIService struct {
	repo *repo.AIRepo
}

func NewAIService(r *repo.AIRepo) *AIService {
	return &AIService{repo: r}
}

func (s *AIService) List(ctx context.Context, tenantID string) ([]repo.AIProvider, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *AIService) Create(ctx context.Context, p repo.AIProvider) error {
	return s.repo.Create(ctx, p)
}

func (s *AIService) Query(ctx context.Context, provider repo.AIProvider, role, input string, contextData any) (string, error) {
	payload := map[string]any{"role": role, "input": input, "context": contextData}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", provider.BaseURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if provider.APIKey != nil {
		req.Header.Set("Authorization", "Bearer "+*provider.APIKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if out, ok := data["output"].(string); ok {
		return out, nil
	}
	return "(no output)", nil
}
