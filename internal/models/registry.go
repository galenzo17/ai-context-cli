package models

import (
	"time"
)

// ModelRegistry manages available AI models
type ModelRegistry struct {
	models      map[string]AIModel
	preferences ModelPreferences
	configMgr   *ConfigManager
}

// NewModelRegistry creates a new model registry with default models
func NewModelRegistry() *ModelRegistry {
	configMgr, _ := NewConfigManager() // Ignore error for now
	
	registry := &ModelRegistry{
		models: make(map[string]AIModel),
		preferences: ModelPreferences{
			LastUsed:      make(map[string]time.Time),
			Favorites:     []string{},
			CustomConfigs: make(map[string]AIModel),
			UsageStats:    make(map[string]int),
		},
		configMgr: configMgr,
	}
	
	// Load default models
	registry.loadDefaultModels()
	
	// Load saved configuration if available
	if configMgr != nil {
		registry.loadSavedConfig()
	}
	
	return registry
}

// loadDefaultModels populates the registry with common AI models
func (r *ModelRegistry) loadDefaultModels() {
	defaultModels := []AIModel{
		{
			ID:          "gpt-4o",
			Name:        "GPT-4o",
			Provider:    "OpenAI",
			Description: "Most capable GPT-4 model, great for complex reasoning and code",
			MaxTokens:   128000,
			CostPer1K:   0.005,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityCodeGeneration),
				string(CapabilityCodeReview),
				string(CapabilityQA),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "https://api.openai.com/v1",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			Status: StatusUnknown,
		},
		{
			ID:          "gpt-4o-mini",
			Name:        "GPT-4o Mini",
			Provider:    "OpenAI",
			Description: "Faster and cheaper GPT-4 model, excellent for most tasks",
			MaxTokens:   128000,
			CostPer1K:   0.0015,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityCodeGeneration),
				string(CapabilityQA),
				string(CapabilitySummarization),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "https://api.openai.com/v1",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			Status: StatusUnknown,
		},
		{
			ID:          "claude-3-5-sonnet",
			Name:        "Claude 3.5 Sonnet",
			Provider:    "Anthropic",
			Description: "Advanced reasoning and code generation, excellent for complex tasks",
			MaxTokens:   200000,
			CostPer1K:   0.003,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityCodeGeneration),
				string(CapabilityCodeReview),
				string(CapabilityQA),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "https://api.anthropic.com/v1",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			Status: StatusUnknown,
		},
		{
			ID:          "claude-3-haiku",
			Name:        "Claude 3 Haiku",
			Provider:    "Anthropic",
			Description: "Fast and efficient model, great for quick tasks and analysis",
			MaxTokens:   200000,
			CostPer1K:   0.00025,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityQA),
				string(CapabilitySummarization),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "https://api.anthropic.com/v1",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			Status: StatusUnknown,
		},
		{
			ID:          "gemini-pro",
			Name:        "Gemini Pro",
			Provider:    "Google",
			Description: "Google's advanced model with multimodal capabilities",
			MaxTokens:   30720,
			CostPer1K:   0.0005,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityCodeGeneration),
				string(CapabilityQA),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "https://generativelanguage.googleapis.com/v1beta",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			Status: StatusUnknown,
		},
		{
			ID:          "ollama-llama3",
			Name:        "Llama 3 (Local)",
			Provider:    "Ollama",
			Description: "Local Llama 3 model running via Ollama (free, private)",
			MaxTokens:   8192,
			CostPer1K:   0.0,
			Capabilities: []string{
				string(CapabilityTextGeneration),
				string(CapabilityCodeGeneration),
				string(CapabilityQA),
			},
			APIConfig: APIConfiguration{
				BaseURL:    "http://localhost:11434/api",
				Timeout:    60 * time.Second,
				RetryCount: 1,
			},
			Status: StatusUnknown,
		},
	}
	
	for _, model := range defaultModels {
		r.models[model.ID] = model
	}
}

// GetAllModels returns all registered models
func (r *ModelRegistry) GetAllModels() []AIModel {
	models := make([]AIModel, 0, len(r.models))
	for _, model := range r.models {
		models = append(models, model)
	}
	return models
}

// GetModel returns a specific model by ID
func (r *ModelRegistry) GetModel(id string) (AIModel, bool) {
	model, exists := r.models[id]
	return model, exists
}

// UpdateModelStatus updates the connection status of a model
func (r *ModelRegistry) UpdateModelStatus(id string, status ConnectionStatus, lastTested *time.Time) {
	if model, exists := r.models[id]; exists {
		model.Status = status
		model.LastTested = lastTested
		r.models[id] = model
	}
}

// UpdateAPIConfig updates the API configuration for a model
func (r *ModelRegistry) UpdateAPIConfig(id string, config APIConfiguration) {
	if model, exists := r.models[id]; exists {
		model.APIConfig = config
		r.models[id] = model
	}
}

// GetPreferences returns the current model preferences
func (r *ModelRegistry) GetPreferences() ModelPreferences {
	return r.preferences
}

// UpdatePreferences updates model preferences
func (r *ModelRegistry) UpdatePreferences(prefs ModelPreferences) {
	r.preferences = prefs
}

// GetModelsByProvider returns models filtered by provider
func (r *ModelRegistry) GetModelsByProvider(provider string) []AIModel {
	var models []AIModel
	for _, model := range r.models {
		if model.Provider == provider {
			models = append(models, model)
		}
	}
	return models
}

// GetProviders returns all unique providers
func (r *ModelRegistry) GetProviders() []string {
	providers := make(map[string]bool)
	for _, model := range r.models {
		providers[model.Provider] = true
	}
	
	result := make([]string, 0, len(providers))
	for provider := range providers {
		result = append(result, provider)
	}
	return result
}

// loadSavedConfig loads saved configuration from disk
func (r *ModelRegistry) loadSavedConfig() {
	if r.configMgr == nil {
		return
	}
	
	config, err := r.configMgr.LoadConfig()
	if err != nil {
		return // Ignore errors, use defaults
	}
	
	// Merge saved models with defaults
	for id, savedModel := range config.Models {
		r.models[id] = savedModel
	}
	
	// Load preferences
	r.preferences = config.Preferences
}

// SaveConfig saves the current configuration to disk
func (r *ModelRegistry) SaveConfig() error {
	if r.configMgr == nil {
		return nil
	}
	
	config := &ModelConfig{
		Models:      r.models,
		Preferences: r.preferences,
	}
	
	return r.configMgr.SaveConfig(config)
}

// SaveModelConfig saves a specific model configuration
func (r *ModelRegistry) SaveModelConfig(model AIModel) error {
	if r.configMgr == nil {
		return nil
	}
	
	r.models[model.ID] = model
	return r.configMgr.SaveModelConfig(model)
}