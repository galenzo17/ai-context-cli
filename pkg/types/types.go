package types

type AIModel struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	APIEndpoint string `json:"api_endpoint"`
	APIKey      string `json:"api_key,omitempty"`
}

type ContextTemplate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Template    string `json:"template"`
	Variables   []string `json:"variables"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatSession struct {
	ID       string        `json:"id"`
	Model    AIModel       `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Context  string        `json:"context"`
}