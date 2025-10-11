package dto

type CreateAIProviderDTO struct {
	Name           string   `json:"name" validate:"required"`
	BaseURL        string   `json:"base_url" validate:"required,url"`
	APIKey         string   `json:"api_key"`
	Roles          []string `json:"roles" validate:"required,dive,required"`
	Models         []string `json:"models"`
	DefaultModel   string   `json:"default_model"`
	PromptTemplate string   `json:"prompt_template"`
}

type UpdateAIProviderDTO struct {
	Name           string   `json:"name" validate:"required"`
	BaseURL        string   `json:"base_url" validate:"required,url"`
	APIKey         string   `json:"api_key"`
	Roles          []string `json:"roles" validate:"required,dive,required"`
	Models         []string `json:"models"`
	DefaultModel   string   `json:"default_model"`
	PromptTemplate string   `json:"prompt_template"`
	IsActive       bool     `json:"is_active"`
}

type QueryAIRequest struct {
	ProviderID string      `json:"provider_id" validate:"required,uuid"`
	Role       string      `json:"role" validate:"required"`
	Input      string      `json:"input" validate:"required"`
	Context    interface{} `json:"context"`
}

type QueryAIResponse struct {
	Output string `json:"output"`
}
