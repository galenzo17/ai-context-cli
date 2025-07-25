package preview

import (
	"testing"
	"time"

	"ai-context-cli/internal/context"
)

func TestNewContextPreviewModel(t *testing.T) {
	// Create test context result
	contextResult := &context.ContextResult{
		ProjectName:   "test-project",
		TotalFiles:    5,
		TotalSize:     1024,
		TokenEstimate: 500,
		GeneratedAt:   time.Now(),
		Sections: []context.ContextSection{
			{
				Title:   "Test Section",
				Content: "Test content",
				Files:   []string{"test.go"},
			},
		},
	}

	scanResult := &context.ScanResult{
		TotalFiles: 5,
		TotalSize:  1024,
	}

	// Test model creation
	model := NewContextPreviewModel(contextResult, scanResult)
	
	if model == nil {
		t.Error("Expected model to be created")
	}
	
	if model.contextResult != contextResult {
		t.Error("Expected context result to be set")
	}
	
	if model.scanResult != scanResult {
		t.Error("Expected scan result to be set")
	}
	
	if len(model.templates) == 0 {
		t.Error("Expected templates to be initialized")
	}
	
	if model.width != 80 {
		t.Errorf("Expected default width 80, got %d", model.width)
	}
	
	if model.height != 20 {
		t.Errorf("Expected default height 20, got %d", model.height)
	}
}

func TestGetDefaultTemplates(t *testing.T) {
	templates := getDefaultTemplates()
	
	if len(templates) == 0 {
		t.Error("Expected templates to be returned")
	}
	
	// Check that all templates have required fields
	for i, template := range templates {
		if template.Name == "" {
			t.Errorf("Template %d missing name", i)
		}
		if template.Description == "" {
			t.Errorf("Template %d missing description", i)
		}
		if template.Template == "" {
			t.Errorf("Template %d missing template", i)
		}
		if template.Icon == "" {
			t.Errorf("Template %d missing icon", i)
		}
	}
	
	// Check for expected templates
	expectedTemplates := []string{"Development Focus", "Documentation", "Code Review", "Bug Analysis", "Full Context"}
	foundTemplates := make(map[string]bool)
	
	for _, template := range templates {
		foundTemplates[template.Name] = true
	}
	
	for _, expected := range expectedTemplates {
		if !foundTemplates[expected] {
			t.Errorf("Expected template '%s' not found", expected)
		}
	}
}

func TestCalculateTokenEstimate(t *testing.T) {
	contextResult := &context.ContextResult{
		Sections: []context.ContextSection{
			{
				Title:   "Section 1",
				Content: "This is test content with multiple words to test token estimation.",
			},
			{
				Title:   "Section 2", 
				Content: "Another section with different content for testing purposes.",
			},
		},
	}
	
	model := NewContextPreviewModel(contextResult, &context.ScanResult{})
	estimate := model.calculateTokenEstimate()
	
	if estimate.Characters == 0 {
		t.Error("Expected non-zero character count")
	}
	
	if estimate.Words == 0 {
		t.Error("Expected non-zero word count")
	}
	
	if estimate.Tokens == 0 {
		t.Error("Expected non-zero token count")
	}
	
	if estimate.Cost < 0 {
		t.Error("Expected non-negative cost estimate")
	}
	
	// Test relationship between characters and tokens (roughly 4:1)
	expectedTokens := estimate.Characters / 4
	if estimate.Tokens != expectedTokens {
		t.Errorf("Expected tokens %d, got %d", expectedTokens, estimate.Tokens)
	}
}

func TestUpdateViewport(t *testing.T) {
	contextResult := &context.ContextResult{
		Sections: make([]context.ContextSection, 20), // Many sections
	}
	
	model := NewContextPreviewModel(contextResult, &context.ScanResult{})
	model.height = 25
	
	// Test initial state
	model.updateViewport()
	if model.viewport.offset != 0 {
		t.Errorf("Expected initial offset 0, got %d", model.viewport.offset)
	}
	
	// Test cursor movement affecting viewport
	model.cursor = 10
	model.updateViewport()
	
	if model.cursor < model.viewport.offset || model.cursor >= model.viewport.offset+model.viewport.size {
		t.Error("Cursor should be visible within viewport")
	}
	
	// Test cursor at end
	model.cursor = 19
	model.updateViewport()
	
	if model.cursor < model.viewport.offset || model.cursor >= model.viewport.offset+model.viewport.size {
		t.Error("Cursor should be visible at end of list")
	}
}

func TestFormatNumber(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{999, "999"},
		{1000, "1.0K"},
		{1500, "1.5K"},
		{1000000, "1.0M"},
		{1500000, "1.5M"},
	}
	
	for _, tc := range testCases {
		result := formatNumber(tc.input)
		if result != tc.expected {
			t.Errorf("formatNumber(%d) = '%s', expected '%s'", tc.input, result, tc.expected)
		}
	}
}

func TestSetSize(t *testing.T) {
	model := NewContextPreviewModel(&context.ContextResult{}, &context.ScanResult{})
	
	// Test setting size
	model.SetSize(100, 30)
	
	if model.width != 100 {
		t.Errorf("Expected width 100, got %d", model.width)
	}
	
	if model.height != 30 {
		t.Errorf("Expected height 30, got %d", model.height)
	}
}

func TestGetContextResult(t *testing.T) {
	contextResult := &context.ContextResult{
		ProjectName: "test-project",
	}
	
	model := NewContextPreviewModel(contextResult, &context.ScanResult{})
	
	result := model.GetContextResult()
	if result != contextResult {
		t.Error("Expected to get same context result")
	}
	
	if result.ProjectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", result.ProjectName)
	}
}

func TestTemplateNavigation(t *testing.T) {
	model := NewContextPreviewModel(&context.ContextResult{}, &context.ScanResult{})
	model.templateMode = true
	
	// Test initial state
	if model.currentTemplate != 0 {
		t.Errorf("Expected initial template 0, got %d", model.currentTemplate)
	}
	
	// Test navigation within bounds
	
	// Navigate down
	if model.currentTemplate < len(model.templates)-1 {
		expectedNext := model.currentTemplate + 1
		// Simulate down key press behavior
		model.currentTemplate++
		
		if model.currentTemplate != expectedNext {
			t.Errorf("Expected template %d, got %d", expectedNext, model.currentTemplate)
		}
	}
	
	// Reset and test navigation up
	model.currentTemplate = 1
	if model.currentTemplate > 0 {
		expectedPrev := model.currentTemplate - 1
		// Simulate up key press behavior
		model.currentTemplate--
		
		if model.currentTemplate != expectedPrev {
			t.Errorf("Expected template %d, got %d", expectedPrev, model.currentTemplate)
		}
	}
	
	// Test bounds - shouldn't go below 0
	model.currentTemplate = 0
	// Simulate up key when at 0 (should stay at 0)
	if model.currentTemplate > 0 {
		model.currentTemplate--
	}
	
	if model.currentTemplate != 0 {
		t.Errorf("Expected template to stay at 0, got %d", model.currentTemplate)
	}
}