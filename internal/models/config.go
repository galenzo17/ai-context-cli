package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConfigManager handles model configuration persistence
type ConfigManager struct {
	configDir  string
	configFile string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".ai-context-cli")
	configFile := filepath.Join(configDir, "models.json")
	
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	return &ConfigManager{
		configDir:  configDir,
		configFile: configFile,
	}, nil
}

// ModelConfig represents the persisted model configuration
type ModelConfig struct {
	Models      map[string]AIModel   `json:"models"`
	Preferences ModelPreferences     `json:"preferences"`
	Version     string               `json:"version"`
}

// LoadConfig loads model configuration from disk
func (cm *ConfigManager) LoadConfig() (*ModelConfig, error) {
	data, err := os.ReadFile(cm.configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return &ModelConfig{
				Models:      make(map[string]AIModel),
				Preferences: ModelPreferences{
					LastUsed:      make(map[string]time.Time),
					Favorites:     []string{},
					CustomConfigs: make(map[string]AIModel),
					UsageStats:    make(map[string]int),
				},
				Version: "1.0",
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config ModelConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves model configuration to disk
func (cm *ConfigManager) SaveConfig(config *ModelConfig) error {
	config.Version = "1.0"
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// SaveModelConfig saves a specific model configuration
func (cm *ConfigManager) SaveModelConfig(model AIModel) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}
	
	if config.Models == nil {
		config.Models = make(map[string]AIModel)
	}
	
	config.Models[model.ID] = model
	
	return cm.SaveConfig(config)
}

// LoadModelConfig loads a specific model configuration
func (cm *ConfigManager) LoadModelConfig(modelID string) (AIModel, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return AIModel{}, err
	}
	
	model, exists := config.Models[modelID]
	if !exists {
		return AIModel{}, fmt.Errorf("model %s not found in config", modelID)
	}
	
	return model, nil
}

// SavePreferences saves model preferences
func (cm *ConfigManager) SavePreferences(prefs ModelPreferences) error {
	config, err := cm.LoadConfig()
	if err != nil {
		return err
	}
	
	config.Preferences = prefs
	
	return cm.SaveConfig(config)
}

// LoadPreferences loads model preferences
func (cm *ConfigManager) LoadPreferences() (ModelPreferences, error) {
	config, err := cm.LoadConfig()
	if err != nil {
		return ModelPreferences{}, err
	}
	
	return config.Preferences, nil
}

// GetConfigPath returns the path to the configuration file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configFile
}

// GetConfigDir returns the path to the configuration directory
func (cm *ConfigManager) GetConfigDir() string {
	return cm.configDir
}