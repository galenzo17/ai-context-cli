package navigation

import (
	"fmt"
	"testing"
)

func TestNavigationStackCreation(t *testing.T) {
	stack := NewNavigationStack()
	
	if len(stack.history) != 0 {
		t.Error("Expected empty history for new navigation stack")
	}
	
	if stack.current != -1 {
		t.Errorf("Expected current to be -1, got %d", stack.current)
	}
	
	if stack.CanGoBack() {
		t.Error("Expected CanGoBack to be false for empty stack")
	}
}

func TestNavigationStackPush(t *testing.T) {
	stack := NewNavigationStack()
	
	// Push first screen
	screen1 := Screen{ID: "test1", Title: "Test 1"}
	stack = stack.Push(screen1)
	
	if len(stack.history) != 1 {
		t.Errorf("Expected history length 1, got %d", len(stack.history))
	}
	
	if stack.current != 0 {
		t.Errorf("Expected current to be 0, got %d", stack.current)
	}
	
	// Push second screen
	screen2 := Screen{ID: "test2", Title: "Test 2"}
	stack = stack.Push(screen2)
	
	if len(stack.history) != 2 {
		t.Errorf("Expected history length 2, got %d", len(stack.history))
	}
	
	if stack.current != 1 {
		t.Errorf("Expected current to be 1, got %d", stack.current)
	}
	
	if !stack.CanGoBack() {
		t.Error("Expected CanGoBack to be true with multiple screens")
	}
}

func TestNavigationStackPop(t *testing.T) {
	stack := NewNavigationStack()
	
	// Add screens
	screen1 := Screen{ID: "test1", Title: "Test 1"}
	screen2 := Screen{ID: "test2", Title: "Test 2"}
	stack = stack.Push(screen1).Push(screen2)
	
	// Pop back
	stack, success := stack.Pop()
	if !success {
		t.Error("Expected pop to succeed")
	}
	
	if stack.current != 0 {
		t.Errorf("Expected current to be 0 after pop, got %d", stack.current)
	}
	
	// Try to pop again
	stack, success = stack.Pop()
	if success {
		t.Error("Expected pop to fail when at first screen")
	}
}

func TestNavigationStackCurrent(t *testing.T) {
	stack := NewNavigationStack()
	
	// Empty stack
	_, ok := stack.Current()
	if ok {
		t.Error("Expected Current to return false for empty stack")
	}
	
	// With screen
	screen := Screen{ID: "test", Title: "Test Screen"}
	stack = stack.Push(screen)
	
	current, ok := stack.Current()
	if !ok {
		t.Error("Expected Current to return true with screen")
	}
	
	if current.ID != "test" {
		t.Errorf("Expected current screen ID 'test', got '%s'", current.ID)
	}
}

func TestNavigationStackPrevious(t *testing.T) {
	stack := NewNavigationStack()
	
	// Empty stack
	_, ok := stack.Previous()
	if ok {
		t.Error("Expected Previous to return false for empty stack")
	}
	
	// Single screen
	screen1 := Screen{ID: "test1", Title: "Test 1"}
	stack = stack.Push(screen1)
	
	_, ok = stack.Previous()
	if ok {
		t.Error("Expected Previous to return false for single screen")
	}
	
	// Multiple screens
	screen2 := Screen{ID: "test2", Title: "Test 2"}
	stack = stack.Push(screen2)
	
	prev, ok := stack.Previous()
	if !ok {
		t.Error("Expected Previous to return true with multiple screens")
	}
	
	if prev.ID != "test1" {
		t.Errorf("Expected previous screen ID 'test1', got '%s'", prev.ID)
	}
}

func TestNavigationStackMaxHistory(t *testing.T) {
	stack := NewNavigationStack()
	
	// Add more screens than max history
	for i := 0; i < 25; i++ {
		screen := Screen{ID: fmt.Sprintf("test%d", i), Title: fmt.Sprintf("Test %d", i)}
		stack = stack.Push(screen)
	}
	
	if len(stack.history) > stack.maxHistory {
		t.Errorf("Expected history length <= %d, got %d", stack.maxHistory, len(stack.history))
	}
}

func TestNavigationStackClear(t *testing.T) {
	stack := NewNavigationStack()
	
	// Add screens
	screen1 := Screen{ID: "test1", Title: "Test 1"}
	screen2 := Screen{ID: "test2", Title: "Test 2"}
	stack = stack.Push(screen1).Push(screen2)
	
	// Clear
	stack = stack.Clear()
	
	if len(stack.history) != 0 {
		t.Errorf("Expected empty history after clear, got %d", len(stack.history))
	}
	
	if stack.current != -1 {
		t.Errorf("Expected current to be -1 after clear, got %d", stack.current)
	}
}

func TestNavigationStackGetPath(t *testing.T) {
	stack := NewNavigationStack()
	
	// Empty stack
	path := stack.GetPath()
	if len(path) != 0 {
		t.Error("Expected empty path for empty stack")
	}
	
	// With screen
	screen := Screen{
		ID:   "test",
		Path: []string{"Root", "Sub", "Test"},
	}
	stack = stack.Push(screen)
	
	path = stack.GetPath()
	if len(path) != 3 {
		t.Errorf("Expected path length 3, got %d", len(path))
	}
	
	if path[0] != "Root" || path[1] != "Sub" || path[2] != "Test" {
		t.Errorf("Expected path [Root, Sub, Test], got %v", path)
	}
}

func TestBreadcrumbRendering(t *testing.T) {
	renderer := NewNavigationRenderer()
	
	// Empty breadcrumbs
	screen := Screen{Breadcrumbs: []Breadcrumb{}}
	result := renderer.RenderBreadcrumbs(screen)
	if result != "" {
		t.Error("Expected empty result for no breadcrumbs")
	}
	
	// Single breadcrumb
	screen = Screen{
		Breadcrumbs: []Breadcrumb{
			{Title: "Home", Active: true},
		},
	}
	result = renderer.RenderBreadcrumbs(screen)
	if result == "" {
		t.Error("Expected non-empty result for breadcrumbs")
	}
	
	// Multiple breadcrumbs
	screen = Screen{
		Breadcrumbs: []Breadcrumb{
			{Title: "Home", Active: false},
			{Title: "Section", Active: false},
			{Title: "Current", Active: true},
		},
	}
	result = renderer.RenderBreadcrumbs(screen)
	if result == "" {
		t.Error("Expected non-empty result for multiple breadcrumbs")
	}
}

func TestBackButtonRendering(t *testing.T) {
	renderer := NewNavigationRenderer()
	
	// Can't go back
	result := renderer.RenderBackButton(false)
	if result != "" {
		t.Error("Expected empty result when can't go back")
	}
	
	// Can go back
	result = renderer.RenderBackButton(true)
	if result == "" {
		t.Error("Expected non-empty result when can go back")
	}
}

func TestFullNavigationRendering(t *testing.T) {
	renderer := NewNavigationRenderer()
	stack := NewNavigationStack()
	
	// Empty stack
	result := renderer.RenderFullNavigation(stack)
	if result != "" {
		t.Error("Expected empty result for empty navigation stack")
	}
	
	// With screen
	screen := Screen{
		ID:    "test",
		Title: "Test",
		Breadcrumbs: []Breadcrumb{
			{Title: "Home", Active: false},
			{Title: "Test", Active: true},
		},
	}
	stack = stack.Push(screen)
	
	result = renderer.RenderFullNavigation(stack)
	if result == "" {
		t.Error("Expected non-empty result for navigation with screen")
	}
}

func TestNavigationCentering(t *testing.T) {
	renderer := NewNavigationRenderer()
	
	// Empty navigation
	result := renderer.CenterNavigation("", 100)
	if result != "" {
		t.Error("Expected empty result for empty navigation")
	}
	
	// Short navigation
	result = renderer.CenterNavigation("Test", 100)
	if result == "" {
		t.Error("Expected non-empty result for short navigation")
	}
	
	// Long navigation (wider than target width)
	longNav := "This is a very long navigation that exceeds the target width"
	result = renderer.CenterNavigation(longNav, 10)
	if result != longNav {
		t.Error("Expected original navigation when too long to center")
	}
}

func TestPredefinedScreens(t *testing.T) {
	screens := []Screen{
		MainMenuScreen,
		AddContextAllScreen,
		AddContextFolderScreen,
		ContextPreviewScreen,
		ModelSelectionScreen,
	}
	
	for _, screen := range screens {
		if screen.ID == "" {
			t.Errorf("Expected non-empty ID for screen %s", screen.Title)
		}
		
		if screen.Title == "" {
			t.Errorf("Expected non-empty Title for screen %s", screen.ID)
		}
		
		if len(screen.Path) == 0 {
			t.Errorf("Expected non-empty Path for screen %s", screen.ID)
		}
		
		if len(screen.Breadcrumbs) == 0 {
			t.Errorf("Expected non-empty Breadcrumbs for screen %s", screen.ID)
		}
	}
	
	// Test main menu specifics
	if MainMenuScreen.ShowBack {
		t.Error("Expected MainMenuScreen.ShowBack to be false")
	}
	
	// Test other screens show back
	if !AddContextAllScreen.ShowBack {
		t.Error("Expected AddContextAllScreen.ShowBack to be true")
	}
}