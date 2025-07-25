package navigation

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Screen represents a navigation screen
type Screen struct {
	ID          string
	Title       string
	ParentID    string
	Path        []string // Full path from root
	ShowBack    bool     // Whether to show back navigation
	Breadcrumbs []Breadcrumb
}

// Breadcrumb represents a single breadcrumb item
type Breadcrumb struct {
	Title  string
	Active bool
}

// NavigationStack manages the navigation history
type NavigationStack struct {
	history    []Screen
	current    int
	maxHistory int
}

// NewNavigationStack creates a new navigation stack
func NewNavigationStack() NavigationStack {
	return NavigationStack{
		history:    make([]Screen, 0),
		current:    -1,
		maxHistory: 20,
	}
}

// Push adds a new screen to the navigation history
func (ns NavigationStack) Push(screen Screen) NavigationStack {
	// Remove any history after current position (for branch navigation)
	if ns.current >= 0 && ns.current < len(ns.history)-1 {
		ns.history = ns.history[:ns.current+1]
	}
	
	// Add new screen
	ns.history = append(ns.history, screen)
	ns.current = len(ns.history) - 1
	
	// Maintain max history size
	if len(ns.history) > ns.maxHistory {
		ns.history = ns.history[1:]
		ns.current--
	}
	
	return ns
}

// Pop removes the current screen and goes back
func (ns NavigationStack) Pop() (NavigationStack, bool) {
	if ns.current <= 0 {
		return ns, false // Can't go back further
	}
	
	ns.current--
	return ns, true
}

// Current returns the current screen
func (ns NavigationStack) Current() (Screen, bool) {
	if ns.current < 0 || ns.current >= len(ns.history) {
		return Screen{}, false
	}
	return ns.history[ns.current], true
}

// Previous returns the previous screen without modifying the stack
func (ns NavigationStack) Previous() (Screen, bool) {
	if ns.current <= 0 {
		return Screen{}, false
	}
	return ns.history[ns.current-1], true
}

// CanGoBack returns true if there's a previous screen
func (ns NavigationStack) CanGoBack() bool {
	return ns.current > 0
}

// GetPath returns the current navigation path
func (ns NavigationStack) GetPath() []string {
	if current, ok := ns.Current(); ok {
		return current.Path
	}
	return []string{}
}

// Clear clears the navigation history
func (ns NavigationStack) Clear() NavigationStack {
	ns.history = ns.history[:0]
	ns.current = -1
	return ns
}

// NavigationRenderer handles the visual rendering of navigation elements
type NavigationRenderer struct {
	breadcrumbStyle lipgloss.Style
	separatorStyle  lipgloss.Style
	activeStyle     lipgloss.Style
	backStyle       lipgloss.Style
}

// NewNavigationRenderer creates a new navigation renderer
func NewNavigationRenderer() NavigationRenderer {
	return NavigationRenderer{
		breadcrumbStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Bold(false),
		separatorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")).
			Bold(false),
		activeStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true),
		backStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true),
	}
}

// RenderBreadcrumbs renders the breadcrumb navigation
func (nr NavigationRenderer) RenderBreadcrumbs(screen Screen) string {
	if len(screen.Breadcrumbs) == 0 {
		return ""
	}
	
	var parts []string
	for i, crumb := range screen.Breadcrumbs {
		var style lipgloss.Style
		if crumb.Active {
			style = nr.activeStyle
		} else {
			style = nr.breadcrumbStyle
		}
		
		parts = append(parts, style.Render(crumb.Title))
		
		// Add separator except for last item
		if i < len(screen.Breadcrumbs)-1 {
			parts = append(parts, nr.separatorStyle.Render(" › "))
		}
	}
	
	return strings.Join(parts, "")
}

// RenderBackButton renders the back navigation indicator
func (nr NavigationRenderer) RenderBackButton(canGoBack bool) string {
	if !canGoBack {
		return ""
	}
	
	return nr.backStyle.Render("← ESC: Back")
}

// RenderFullNavigation renders the complete navigation bar
func (nr NavigationRenderer) RenderFullNavigation(stack NavigationStack) string {
	current, hasCurrent := stack.Current()
	if !hasCurrent {
		return ""
	}
	
	var parts []string
	
	// Back button (left side)
	if backButton := nr.RenderBackButton(stack.CanGoBack()); backButton != "" {
		parts = append(parts, backButton)
	}
	
	// Breadcrumbs (center/right side)
	if breadcrumbs := nr.RenderBreadcrumbs(current); breadcrumbs != "" {
		if len(parts) > 0 {
			// Add spacing between back button and breadcrumbs
			parts = append(parts, strings.Repeat(" ", 4))
		}
		parts = append(parts, breadcrumbs)
	}
	
	return strings.Join(parts, "")
}

// CenterNavigation centers the navigation bar within a given width
func (nr NavigationRenderer) CenterNavigation(navigation string, width int) string {
	if navigation == "" {
		return ""
	}
	
	// Calculate actual width without ANSI codes
	actualWidth := lipgloss.Width(navigation)
	if actualWidth >= width {
		return navigation
	}
	
	padding := (width - actualWidth) / 2
	return strings.Repeat(" ", padding) + navigation
}

// Common screen definitions
var (
	MainMenuScreen = Screen{
		ID:       "main_menu",
		Title:    "Main Menu",
		Path:     []string{"Context Engine"},
		ShowBack: false,
		Breadcrumbs: []Breadcrumb{
			{Title: "Context Engine", Active: true},
		},
	}
	
	AddContextAllScreen = Screen{
		ID:       "add_context_all",
		Title:    "Add Context - All Files",
		ParentID: "main_menu",
		Path:     []string{"Context Engine", "Add Context", "All Files"},
		ShowBack: true,
		Breadcrumbs: []Breadcrumb{
			{Title: "Context Engine", Active: false},
			{Title: "Add Context", Active: false},
			{Title: "All Files", Active: true},
		},
	}
	
	AddContextFolderScreen = Screen{
		ID:       "add_context_folder",
		Title:    "Add Context - Folder",
		ParentID: "main_menu",
		Path:     []string{"Context Engine", "Add Context", "Folder"},
		ShowBack: true,
		Breadcrumbs: []Breadcrumb{
			{Title: "Context Engine", Active: false},
			{Title: "Add Context", Active: false},
			{Title: "Folder", Active: true},
		},
	}
	
	ContextPreviewScreen = Screen{
		ID:       "context_preview",
		Title:    "Context Preview",
		ParentID: "main_menu",
		Path:     []string{"Context Engine", "Context Preview"},
		ShowBack: true,
		Breadcrumbs: []Breadcrumb{
			{Title: "Context Engine", Active: false},
			{Title: "Context Preview", Active: true},
		},
	}
	
	ModelSelectionScreen = Screen{
		ID:       "model_selection",
		Title:    "Model Selection",
		ParentID: "main_menu",
		Path:     []string{"Context Engine", "Model Selection"},
		ShowBack: true,
		Breadcrumbs: []Breadcrumb{
			{Title: "Context Engine", Active: false},
			{Title: "Model Selection", Active: true},
		},
	}
)