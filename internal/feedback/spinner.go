package feedback

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerMsg is sent when the spinner should update
type SpinnerMsg struct{}

// SpinnerModel represents a loading spinner
type SpinnerModel struct {
	frames   []string
	current  int
	active   bool
	message  string
	style    lipgloss.Style
}

// NewSpinner creates a new spinner instance
func NewSpinner(message string) SpinnerModel {
	return SpinnerModel{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		current: 0,
		active:  false,
		message: message,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true),
	}
}

// Start activates the spinner
func (s SpinnerModel) Start() SpinnerModel {
	s.active = true
	return s
}

// Stop deactivates the spinner
func (s SpinnerModel) Stop() SpinnerModel {
	s.active = false
	return s
}

// SetMessage updates the spinner message
func (s SpinnerModel) SetMessage(message string) SpinnerModel {
	s.message = message
	return s
}

// Update handles spinner tick messages
func (s SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	switch msg.(type) {
	case SpinnerMsg:
		if s.active {
			s.current = (s.current + 1) % len(s.frames)
			return s, s.tick()
		}
	}
	return s, nil
}

// View renders the spinner
func (s SpinnerModel) View() string {
	if !s.active {
		return ""
	}
	
	spinner := s.style.Render(s.frames[s.current])
	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Render(s.message)
	
	return spinner + " " + message
}

// tick returns a command that sends a SpinnerMsg after 100ms
func (s SpinnerModel) tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerMsg{}
	})
}

// InitSpinner returns a command to start the spinner
func (s SpinnerModel) InitSpinner() tea.Cmd {
	if s.active {
		return s.tick()
	}
	return nil
}