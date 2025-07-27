package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ConnectionTester handles testing connections to AI model APIs
type ConnectionTester struct {
	client *http.Client
}

// NewConnectionTester creates a new connection tester
func NewConnectionTester() *ConnectionTester {
	return &ConnectionTester{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TestConnection tests the connection to a model's API
func (ct *ConnectionTester) TestConnection(model AIModel) TestResult {
	start := time.Now()
	
	result := TestResult{
		Success:   false,
		Timestamp: start,
	}
	
	// Check if API key is configured
	if model.APIConfig.APIKey == "" {
		result.Error = "API key not configured"
		result.Latency = time.Since(start)
		return result
	}
	
	// Test based on provider
	switch model.Provider {
	case "OpenAI":
		return ct.testOpenAI(model, start)
	case "Anthropic":
		return ct.testAnthropic(model, start)
	case "Google":
		return ct.testGoogle(model, start)
	case "Ollama":
		return ct.testOllama(model, start)
	default:
		result.Error = fmt.Sprintf("Testing not implemented for provider: %s", model.Provider)
		result.Latency = time.Since(start)
		return result
	}
}

// testOpenAI tests OpenAI API connection
func (ct *ConnectionTester) testOpenAI(model AIModel, start time.Time) TestResult {
	result := TestResult{
		Success:   false,
		Timestamp: start,
	}
	
	// Create a simple test request to list models
	url := model.APIConfig.BaseURL + "/models"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	req.Header.Set("Authorization", "Bearer "+model.APIConfig.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ct.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("Request failed: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	
	result.Latency = time.Since(start)
	
	if resp.StatusCode == 200 {
		result.Success = true
	} else {
		result.Error = fmt.Sprintf("API returned status %d", resp.StatusCode)
	}
	
	return result
}

// testAnthropic tests Anthropic API connection  
func (ct *ConnectionTester) testAnthropic(model AIModel, start time.Time) TestResult {
	result := TestResult{
		Success:   false,
		Timestamp: start,
	}
	
	// Create a minimal test request
	testPayload := map[string]interface{}{
		"model":      model.ID,
		"max_tokens": 1,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Hi",
			},
		},
	}
	
	jsonData, err := json.Marshal(testPayload)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to marshal request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	url := model.APIConfig.BaseURL + "/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	req.Header.Set("x-api-key", model.APIConfig.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ct.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("Request failed: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	
	result.Latency = time.Since(start)
	
	if resp.StatusCode == 200 {
		result.Success = true
	} else {
		result.Error = fmt.Sprintf("API returned status %d", resp.StatusCode)
	}
	
	return result
}

// testGoogle tests Google API connection
func (ct *ConnectionTester) testGoogle(model AIModel, start time.Time) TestResult {
	result := TestResult{
		Success:   false,
		Timestamp: start,
	}
	
	// For Google, test with a simple generate request
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", 
		model.APIConfig.BaseURL, model.ID, model.APIConfig.APIKey)
	
	testPayload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "Hi"},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": 1,
		},
	}
	
	jsonData, err := json.Marshal(testPayload)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to marshal request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ct.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("Request failed: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	
	result.Latency = time.Since(start)
	
	if resp.StatusCode == 200 {
		result.Success = true
	} else {
		result.Error = fmt.Sprintf("API returned status %d", resp.StatusCode)
	}
	
	return result
}

// testOllama tests Ollama local API connection
func (ct *ConnectionTester) testOllama(model AIModel, start time.Time) TestResult {
	result := TestResult{
		Success:   false,
		Timestamp: start,
	}
	
	// For Ollama, test by listing models
	url := model.APIConfig.BaseURL + "/tags"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Latency = time.Since(start)
		return result
	}
	
	// Set a shorter timeout for local requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	
	resp, err := ct.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("Request failed (is Ollama running?): %v", err)
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	
	result.Latency = time.Since(start)
	
	if resp.StatusCode == 200 {
		result.Success = true
	} else {
		result.Error = fmt.Sprintf("Ollama returned status %d", resp.StatusCode)
	}
	
	return result
}

// TestAllModels tests connections to all models in the registry
func (ct *ConnectionTester) TestAllModels(registry *ModelRegistry) map[string]TestResult {
	results := make(map[string]TestResult)
	models := registry.GetAllModels()
	
	for _, model := range models {
		result := ct.TestConnection(model)
		results[model.ID] = result
		
		// Update model status in registry
		status := StatusError
		if result.Success {
			status = StatusConnected
		} else {
			status = StatusDisconnected
		}
		
		registry.UpdateModelStatus(model.ID, status, &result.Timestamp)
	}
	
	return results
}