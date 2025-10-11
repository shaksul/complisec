package dto

// ChatMessage представляет одно сообщение в чате
type ChatMessage struct {
	Role    string `json:"role"` // "user" | "assistant" | "system"
	Content string `json:"content"`
}

// SendChatMessageRequest - запрос на отправку сообщения в AI чат
type SendChatMessageRequest struct {
	ProviderID string        `json:"provider_id" validate:"required"`
	Messages   []ChatMessage `json:"messages" validate:"required,min=1"`
	Model      string        `json:"model"`
	Stream     bool          `json:"stream"`
}

// SendChatMessageResponse - ответ от AI чата
type SendChatMessageResponse struct {
	Message ChatMessage `json:"message"`
	Model   string      `json:"model"`
}

// GetModelsResponse - список доступных моделей
type GetModelsResponse struct {
	Models []string `json:"models"`
}
