package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"ai-context-cli/pkg/types"
)

type Config struct {
	DefaultModel      string                    `json:"default_model"`
	Models            []types.AIModel           `json:"models"`
	ContextTemplates  []types.ContextTemplate   `json:"context_templates"`
	ConfigDir         string                    `json:"-"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".ai-context-cli")
	configFile := filepath.Join(configDir, "config.json")

	config := &Config{
		ConfigDir: configDir,
		Models: []types.AIModel{
			{
				Name:     "gpt-3.5-turbo",
				Provider: "openai",
				APIEndpoint: "https://api.openai.com/v1/chat/completions",
			},
		},
		ContextTemplates: []types.ContextTemplate{
			{
				ID:          "default",
				Name:        "Default Context",
				Description: "Basic context template",
				Template:    "You are a helpful AI assistant. {{.context}}",
				Variables:   []string{"context"},
			},
		},
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
		return config, config.Save()
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	config.ConfigDir = configDir
	return config, nil
}

func (c *Config) Save() error {
	configFile := filepath.Join(c.ConfigDir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}