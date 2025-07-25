package feedback

import (
	"testing"
)

func TestSpinnerCreation(t *testing.T) {
	spinner := NewSpinner("Loading test...")
	
	if spinner.message != "Loading test..." {
		t.Errorf("Expected message 'Loading test...', got '%s'", spinner.message)
	}
	
	if spinner.active {
		t.Error("Expected spinner to be inactive initially")
	}
	
	if len(spinner.frames) == 0 {
		t.Error("Expected spinner to have frames")
	}
}

func TestSpinnerStartStop(t *testing.T) {
	spinner := NewSpinner("Test")
	
	// Test start
	spinner = spinner.Start()
	if !spinner.active {
		t.Error("Expected spinner to be active after Start()")
	}
	
	// Test stop
	spinner = spinner.Stop()
	if spinner.active {
		t.Error("Expected spinner to be inactive after Stop()")
	}
}

func TestSpinnerView(t *testing.T) {
	spinner := NewSpinner("Loading...")
	
	// Inactive spinner should return empty view
	view := spinner.View()
	if view != "" {
		t.Error("Expected empty view for inactive spinner")
	}
	
	// Active spinner should return content
	spinner = spinner.Start()
	view = spinner.View()
	if view == "" {
		t.Error("Expected non-empty view for active spinner")
	}
}

func TestProgressCreation(t *testing.T) {
	progress := NewProgress(100, "Processing...")
	
	if progress.total != 100 {
		t.Errorf("Expected total 100, got %d", progress.total)
	}
	
	if progress.current != 0 {
		t.Errorf("Expected current 0, got %d", progress.current)
	}
	
	if progress.message != "Processing..." {
		t.Errorf("Expected message 'Processing...', got '%s'", progress.message)
	}
}

func TestProgressSetProgress(t *testing.T) {
	progress := NewProgress(10, "Test")
	
	// Test normal progress
	progress = progress.SetProgress(5)
	if progress.current != 5 {
		t.Errorf("Expected current 5, got %d", progress.current)
	}
	
	// Test overflow protection
	progress = progress.SetProgress(15)
	if progress.current != 10 {
		t.Errorf("Expected current 10 (clamped), got %d", progress.current)
	}
	
	// Test negative protection
	progress = progress.SetProgress(-5)
	if progress.current != 0 {
		t.Errorf("Expected current 0 (clamped), got %d", progress.current)
	}
}

func TestProgressPercentage(t *testing.T) {
	progress := NewProgress(100, "Test")
	
	// Test 0%
	if progress.Percentage() != 0 {
		t.Errorf("Expected 0%%, got %d%%", progress.Percentage())
	}
	
	// Test 50%
	progress = progress.SetProgress(50)
	if progress.Percentage() != 50 {
		t.Errorf("Expected 50%%, got %d%%", progress.Percentage())
	}
	
	// Test 100%
	progress = progress.SetProgress(100)
	if progress.Percentage() != 100 {
		t.Errorf("Expected 100%%, got %d%%", progress.Percentage())
	}
	
	// Test completion
	if !progress.IsComplete() {
		t.Error("Expected progress to be complete at 100%")
	}
}

func TestProgressView(t *testing.T) {
	progress := NewProgress(10, "Testing")
	
	// Test empty progress
	view := progress.View()
	if view == "" {
		t.Error("Expected non-empty view for progress bar")
	}
	
	// Test with progress
	progress = progress.SetProgress(5)
	view = progress.View()
	if view == "" {
		t.Error("Expected non-empty view for progress with value")
	}
}

func TestToastCreation(t *testing.T) {
	toast := NewToast(1, "Test message", ToastSuccess)
	
	if toast.id != 1 {
		t.Errorf("Expected ID 1, got %d", toast.id)
	}
	
	if toast.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", toast.message)
	}
	
	if toast.toastType != ToastSuccess {
		t.Errorf("Expected ToastSuccess type, got %v", toast.toastType)
	}
	
	if !toast.visible {
		t.Error("Expected toast to be visible initially")
	}
}

func TestToastTypes(t *testing.T) {
	testCases := []struct {
		toastType ToastType
		name      string
	}{
		{ToastSuccess, "success"},
		{ToastError, "error"},
		{ToastWarning, "warning"},
		{ToastInfo, "info"},
	}
	
	for _, tc := range testCases {
		toast := NewToast(1, "Test", tc.toastType)
		view := toast.View()
		
		if view == "" {
			t.Errorf("Expected non-empty view for %s toast", tc.name)
		}
	}
}

func TestToastVisibility(t *testing.T) {
	toast := NewToast(1, "Test", ToastInfo)
	
	// Should be visible initially
	if !toast.IsVisible() {
		t.Error("Expected toast to be visible initially")
	}
	
	// Hide toast
	toast = toast.Hide()
	if toast.IsVisible() {
		t.Error("Expected toast to be hidden after Hide()")
	}
	
	// View should be empty when hidden
	view := toast.View()
	if view != "" {
		t.Error("Expected empty view for hidden toast")
	}
}

func TestToastManager(t *testing.T) {
	manager := NewToastManager()
	
	if len(manager.toasts) != 0 {
		t.Error("Expected empty toast manager initially")
	}
	
	// Add toast
	manager, _ = manager.AddToast("Test message", ToastSuccess)
	if len(manager.toasts) != 1 {
		t.Errorf("Expected 1 toast, got %d", len(manager.toasts))
	}
	
	// Check view
	view := manager.View()
	if view == "" {
		t.Error("Expected non-empty view with toasts")
	}
}

func TestToastManagerMaxSize(t *testing.T) {
	manager := NewToastManager()
	
	// Add more toasts than max size
	for i := 0; i < 7; i++ {
		manager, _ = manager.AddToast("Test", ToastInfo)
	}
	
	if len(manager.toasts) > manager.maxSize {
		t.Errorf("Expected max %d toasts, got %d", manager.maxSize, len(manager.toasts))
	}
}

func TestToastExpiration(t *testing.T) {
	toast := NewToast(1, "Test", ToastSuccess)
	
	// Simulate expiration message
	expireMsg := ToastExpireMsg{ID: 1}
	updatedToast, _ := toast.Update(expireMsg)
	
	if updatedToast.IsVisible() {
		t.Error("Expected toast to be hidden after expiration")
	}
	
	// Wrong ID should not affect toast
	expireMsg = ToastExpireMsg{ID: 2}
	updatedToast, _ = toast.Update(expireMsg)
	
	if !updatedToast.IsVisible() {
		t.Error("Expected toast to remain visible with wrong expiration ID")
	}
}