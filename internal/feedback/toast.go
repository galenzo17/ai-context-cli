package feedback

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ToastType represents the type of toast notification
type ToastType int

const (
	ToastSuccess ToastType = iota
	ToastError
	ToastWarning
	ToastInfo
)

// ToastExpireMsg is sent when a toast should expire
type ToastExpireMsg struct {
	ID int
}

// ToastModel represents a toast notification
type ToastModel struct {
	id       int
	message  string
	toastType ToastType
	visible  bool
	duration time.Duration
	style    lipgloss.Style
}

// NewToast creates a new toast notification
func NewToast(id int, message string, toastType ToastType) ToastModel {
	var style lipgloss.Style
	
	switch toastType {
	case ToastSuccess:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#10B981")).
			Padding(0, 1).
			Bold(true)
	case ToastError:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#EF4444")).
			Padding(0, 1).
			Bold(true)
	case ToastWarning:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#F59E0B")).
			Padding(0, 1).
			Bold(true)
	case ToastInfo:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3B82F6")).
			Padding(0, 1).
			Bold(true)
	}
	
	return ToastModel{
		id:       id,
		message:  message,
		toastType: toastType,
		visible:  true,
		duration: 3 * time.Second,
		style:    style,
	}
}

// Show displays the toast and sets auto-expire timer
func (t ToastModel) Show() (ToastModel, tea.Cmd) {
	t.visible = true
	return t, t.expireAfter()
}

// Hide hides the toast
func (t ToastModel) Hide() ToastModel {
	t.visible = false
	return t
}

// Update handles toast expire messages
func (t ToastModel) Update(msg tea.Msg) (ToastModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ToastExpireMsg:
		if msg.ID == t.id {
			t.visible = false
		}
	}
	return t, nil
}

// View renders the toast if visible
func (t ToastModel) View() string {
	if !t.visible {
		return ""
	}
	
	// Add icon based on type
	var icon string
	switch t.toastType {
	case ToastSuccess:
		icon = "✅ "
	case ToastError:
		icon = "❌ "
	case ToastWarning:
		icon = "⚠️ "
	case ToastInfo:
		icon = "ℹ️ "
	}
	
	return t.style.Render(icon + t.message)
}

// IsVisible returns whether the toast is currently visible
func (t ToastModel) IsVisible() bool {
	return t.visible
}

// ID returns the toast ID
func (t ToastModel) ID() int {
	return t.id
}

// expireAfter returns a command that will hide this toast after the duration
func (t ToastModel) expireAfter() tea.Cmd {
	return tea.Tick(t.duration, func(time.Time) tea.Msg {
		return ToastExpireMsg{ID: t.id}
	})
}

// ToastManager manages multiple toast notifications
type ToastManager struct {
	toasts  []ToastModel
	nextID  int
	maxSize int
}

// NewToastManager creates a new toast manager
func NewToastManager() ToastManager {
	return ToastManager{
		toasts:  make([]ToastModel, 0),
		nextID:  1,
		maxSize: 5,
	}
}

// AddToast adds a new toast notification
func (tm ToastManager) AddToast(message string, toastType ToastType) (ToastManager, tea.Cmd) {
	toast := NewToast(tm.nextID, message, toastType)
	tm.nextID++
	
	// Remove oldest toast if at max capacity
	if len(tm.toasts) >= tm.maxSize {
		tm.toasts = tm.toasts[1:]
	}
	
	tm.toasts = append(tm.toasts, toast)
	
	// Show the toast
	toast, cmd := toast.Show()
	tm.toasts[len(tm.toasts)-1] = toast
	
	return tm, cmd
}

// Update updates all toasts
func (tm ToastManager) Update(msg tea.Msg) (ToastManager, tea.Cmd) {
	var cmds []tea.Cmd
	
	for i, toast := range tm.toasts {
		updatedToast, cmd := toast.Update(msg)
		tm.toasts[i] = updatedToast
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	// Remove invisible toasts
	var visibleToasts []ToastModel
	for _, toast := range tm.toasts {
		if toast.IsVisible() {
			visibleToasts = append(visibleToasts, toast)
		}
	}
	tm.toasts = visibleToasts
	
	return tm, tea.Batch(cmds...)
}

// View renders all visible toasts
func (tm ToastManager) View() string {
	if len(tm.toasts) == 0 {
		return ""
	}
	
	var toastViews []string
	for _, toast := range tm.toasts {
		if view := toast.View(); view != "" {
			toastViews = append(toastViews, view)
		}
	}
	
	if len(toastViews) == 0 {
		return ""
	}
	
	// Stack toasts vertically
	return lipgloss.JoinVertical(lipgloss.Left, toastViews...)
}