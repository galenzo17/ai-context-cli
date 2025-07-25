package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	model := NewModel()
	
	if len(model.choices) == 0 {
		t.Error("Expected choices to be populated")
	}
	
	if model.cursor != 0 {
		t.Error("Expected cursor to start at 0")
	}
	
	if model.selected == nil {
		t.Error("Expected selected map to be initialized")
	}
}

func TestModelUpdate(t *testing.T) {
	model := NewModel()
	
	// Test cursor movement
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)
	
	if m.cursor != 1 {
		t.Errorf("Expected cursor to be 1, got %d", m.cursor)
	}
	
	// Test quit
	msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(msg)
	
	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}
}

func TestModelView(t *testing.T) {
	model := NewModel()
	view := model.View()
	
	if view == "" {
		t.Error("Expected view to return non-empty string")
	}
}