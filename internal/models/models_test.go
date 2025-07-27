package models

import (
	"testing"
	"time"
)

func TestNewModelRegistry(t *testing.T) {
	registry := NewModelRegistry()
	
	if registry == nil {
		t.Error("Expected registry to be initialized")
	}
	
	models := registry.GetAllModels()
	if len(models) == 0 {
		t.Error("Expected models to be populated")
	}
	
	// Check that we have the expected models
	expectedModels := []string{
		"gpt-4o",
		"gpt-4o-mini", 
		"claude-3-5-sonnet",
		"claude-3-haiku",
		"gemini-pro",
		"ollama-llama3",
	}
	
	for _, expectedID := range expectedModels {
		model, exists := registry.GetModel(expectedID)
		if !exists {
			t.Errorf("Expected model %s to exist", expectedID)
		}
		
		if model.ID != expectedID {
			t.Errorf("Expected model ID %s, got %s", expectedID, model.ID)
		}
		
		if model.Name == "" {
			t.Error("Expected model to have a name")
		}
		
		if model.Provider == "" {
			t.Error("Expected model to have a provider")
		}
	}
}

func TestModelRegistryOperations(t *testing.T) {
	registry := NewModelRegistry()
	
	// Test getting a specific model
	model, exists := registry.GetModel("gpt-4o")
	if !exists {
		t.Error("Expected GPT-4o model to exist")
	}
	
	if model.Provider != "OpenAI" {
		t.Errorf("Expected provider OpenAI, got %s", model.Provider)
	}
	
	// Test updating model status
	now := time.Now()
	registry.UpdateModelStatus("gpt-4o", StatusConnected, &now)
	
	updated, _ := registry.GetModel("gpt-4o")
	if updated.Status != StatusConnected {
		t.Error("Expected model status to be updated to connected")
	}
	
	if updated.LastTested == nil {
		t.Error("Expected last tested time to be set")
	}
	
	// Test filtering by provider
	openaiModels := registry.GetModelsByProvider("OpenAI")
	if len(openaiModels) < 2 {
		t.Error("Expected at least 2 OpenAI models")
	}
	
	for _, model := range openaiModels {
		if model.Provider != "OpenAI" {
			t.Error("Expected all models to be from OpenAI")
		}
	}
	
	// Test getting providers
	providers := registry.GetProviders()
	expectedProviders := []string{"OpenAI", "Anthropic", "Google", "Ollama"}
	
	for _, expected := range expectedProviders {
		found := false
		for _, provider := range providers {
			if provider == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected provider %s to be in list", expected)
		}
	}
}

func TestNewSelectorModel(t *testing.T) {
	selector := NewSelectorModel()
	
	if selector == nil {
		t.Error("Expected selector to be initialized")
	}
	
	if selector.registry == nil {
		t.Error("Expected selector to have a registry")
	}
	
	if len(selector.models) == 0 {
		t.Error("Expected selector to have models")
	}
	
	if selector.cursor != 0 {
		t.Error("Expected cursor to start at 0")
	}
	
	if selector.viewMode != ViewModeList {
		t.Error("Expected default view mode to be list")
	}
}

func TestViewModeToggle(t *testing.T) {
	selector := NewSelectorModel()
	
	// Test view mode cycling
	if selector.viewMode != ViewModeList {
		t.Error("Expected initial view mode to be list")
	}
	
	selector.cycleViewMode()
	if selector.viewMode != ViewModeGrid {
		t.Error("Expected view mode to cycle to grid")
	}
	
	selector.cycleViewMode()
	if selector.viewMode != ViewModeComparison {
		t.Error("Expected view mode to cycle to comparison")
	}
	
	selector.cycleViewMode()
	if selector.viewMode != ViewModeList {
		t.Error("Expected view mode to cycle back to list")
	}
}