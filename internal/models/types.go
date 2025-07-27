package models

import "time"

// AIModel represents an AI model configuration
type AIModel struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Provider     string             `json:"provider"`
	Description  string             `json:"description"`
	MaxTokens    int                `json:"max_tokens"`
	CostPer1K    float64           `json:"cost_per_1k"`
	Capabilities []string          `json:"capabilities"`
	APIConfig    APIConfiguration  `json:"api_config"`
	Status       ConnectionStatus  `json:"status"`
	LastTested   *time.Time        `json:"last_tested,omitempty"`
}

// APIConfiguration holds API settings for a model
type APIConfiguration struct {
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key"`
	Headers    map[string]string `json:"headers,omitempty"`
	Timeout    time.Duration     `json:"timeout"`
	RetryCount int               `json:"retry_count"`
}

// ConnectionStatus represents the current connection state
type ConnectionStatus string

const (
	StatusUnknown     ConnectionStatus = "unknown"
	StatusConnected   ConnectionStatus = "connected"
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusError       ConnectionStatus = "error"
)

// ModelCapability represents different model capabilities
type ModelCapability string

const (
	CapabilityTextGeneration ModelCapability = "text_generation"
	CapabilityCodeGeneration ModelCapability = "code_generation"
	CapabilityCodeReview     ModelCapability = "code_review"
	CapabilityQA            ModelCapability = "question_answering"
	CapabilitySummarization ModelCapability = "summarization"
	CapabilityTranslation   ModelCapability = "translation"
)

// TestResult represents the result of a connection test
type TestResult struct {
	Success   bool          `json:"success"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// ModelPreferences holds user preferences for model selection
type ModelPreferences struct {
	DefaultModelID string                     `json:"default_model_id"`
	LastUsed       map[string]time.Time       `json:"last_used"`
	Favorites      []string                   `json:"favorites"`
	CustomConfigs  map[string]AIModel         `json:"custom_configs"`
	UsageStats     map[string]int             `json:"usage_stats"`
}