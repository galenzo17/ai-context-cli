package feedback

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressModel represents a progress bar
type ProgressModel struct {
	current   int
	total     int
	width     int
	message   string
	style     lipgloss.Style
	barStyle  lipgloss.Style
	fillStyle lipgloss.Style
}

// NewProgress creates a new progress bar
func NewProgress(total int, message string) ProgressModel {
	return ProgressModel{
		current: 0,
		total:   total,
		width:   40,
		message: message,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
		barStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")).
			Bold(false),
		fillStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true),
	}
}

// SetProgress updates the current progress
func (p ProgressModel) SetProgress(current int) ProgressModel {
	if current > p.total {
		current = p.total
	}
	if current < 0 {
		current = 0
	}
	p.current = current
	return p
}

// Increment increases progress by 1
func (p ProgressModel) Increment() ProgressModel {
	return p.SetProgress(p.current + 1)
}

// SetMessage updates the progress message
func (p ProgressModel) SetMessage(message string) ProgressModel {
	p.message = message
	return p
}

// IsComplete returns true if progress is at 100%
func (p ProgressModel) IsComplete() bool {
	return p.current >= p.total
}

// Percentage returns the current percentage (0-100)
func (p ProgressModel) Percentage() int {
	if p.total == 0 {
		return 100
	}
	return (p.current * 100) / p.total
}

// View renders the progress bar
func (p ProgressModel) View() string {
	if p.total == 0 {
		return ""
	}

	// Calculate progress
	percentage := p.Percentage()
	filled := (p.current * p.width) / p.total
	
	// Build progress bar
	var bar strings.Builder
	bar.WriteString("[")
	
	// Filled portion
	if filled > 0 {
		bar.WriteString(p.fillStyle.Render(strings.Repeat("█", filled)))
	}
	
	// Empty portion
	empty := p.width - filled
	if empty > 0 {
		bar.WriteString(p.barStyle.Render(strings.Repeat("░", empty)))
	}
	
	bar.WriteString("]")
	
	// Progress text
	progressText := fmt.Sprintf(" %d%% (%d/%d)", percentage, p.current, p.total)
	
	// Message
	message := ""
	if p.message != "" {
		message = p.style.Render(p.message) + "\n"
	}
	
	return message + bar.String() + p.style.Render(progressText)
}

// ViewCompact renders a compact version of the progress bar
func (p ProgressModel) ViewCompact() string {
	if p.total == 0 {
		return ""
	}
	
	percentage := p.Percentage()
	return fmt.Sprintf("%s %d%% (%d/%d)", 
		p.style.Render(p.message), 
		percentage, 
		p.current, 
		p.total)
}