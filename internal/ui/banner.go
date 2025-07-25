package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const version = "v0.1.0"

var logoArt = []string{
	"╔══════════════════════════════════════════════════════════════════╗",
	"║  ___ ___  _  _ _____ _____  _______   ___ _  _  ___ ___ _  _ ___   ║",
	"║ / __/ _ \\| \\| |_   _| __\\ \\/ /_   _| | __| \\| |/ __|_ _| \\| | __|  ║",
	"║| (_| (_) | .` | | | | _| >  <  | |   | _|| .` | (_ || || .` | _|   ║",
	"║ \\___\\___/|_|\\_| |_| |___/_/\\_\\ |_|   |___|_|\\_|\\___|___|_|\\_|___|  ║",
	"╚══════════════════════════════════════════════════════════════════╝",
}

type BannerConfig struct {
	Width       int
	ShowVersion bool
	ColorScheme string
}

func GetTerminalWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil || width < 80 {
		return 80 // Minimum width fallback
	}
	return width
}

func centerText(text string, width int) string {
	textLen := lipgloss.Width(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

func RenderBanner(config BannerConfig) string {
	var result strings.Builder
	
	// Get terminal width if not specified
	if config.Width == 0 {
		config.Width = GetTerminalWidth()
	}

	// Create gradient colors - violet to blue to cyan
	gradientColors := []lipgloss.Color{
		lipgloss.Color("#8B5CF6"), // Violet
		lipgloss.Color("#7C3AED"), // Purple  
		lipgloss.Color("#6366F1"), // Indigo
		lipgloss.Color("#3B82F6"), // Blue
		lipgloss.Color("#0EA5E9"), // Sky blue
		lipgloss.Color("#06B6D4"), // Cyan
		lipgloss.Color("#14B8A6"), // Teal
	}

	// Create styles for different parts
	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6366F1")).
		Bold(true)
	
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#10B981")).
		Italic(true)

	// Check if terminal is wide enough for full logo
	logoWidth := 70 // Width of the boxed design with single-line ASCII
	if config.Width < logoWidth+4 {
		// Use compact version for narrow terminals
		compactLogo := []string{
			"╔═══════════════════════════╗",
			"║      Context Engine       ║",
			"╚═══════════════════════════╝",
		}
		
		for i, line := range compactLogo {
			// Use gradient colors for compact version too
			colorIndex := i % len(gradientColors)
			compactStyle := lipgloss.NewStyle().
				Foreground(gradientColors[colorIndex]).
				Bold(true)
			centeredLine := centerText(compactStyle.Render(line), config.Width)
			result.WriteString(centeredLine + "\n")
		}
	} else {
		// Render full ASCII art logo with gradient colors
		for i, line := range logoArt {
			if line == "" {
				result.WriteString("\n")
				continue
			}
			
			// Apply different colors based on line position
			var lineStyle lipgloss.Style
			if i == 0 || i == len(logoArt)-1 {
				// Border lines use border style
				lineStyle = borderStyle
			} else {
				// Content lines use gradient colors
				colorIndex := (i - 1) % len(gradientColors)
				lineStyle = lipgloss.NewStyle().
					Foreground(gradientColors[colorIndex]).
					Bold(true)
			}
			
			centeredLine := centerText(lineStyle.Render(line), config.Width)
			result.WriteString(centeredLine + "\n")
		}
	}

	// Add version in bottom right if requested
	if config.ShowVersion {
		result.WriteString("\n")
		versionText := fmt.Sprintf("%s", version)
		styledVersion := versionStyle.Render(versionText)
		// Use lipgloss.Width to account for ANSI codes
		actualWidth := lipgloss.Width(styledVersion)
		padding := config.Width - actualWidth
		if padding < 0 {
			padding = 0
		}
		versionLine := strings.Repeat(" ", padding) + styledVersion
		result.WriteString(versionLine + "\n")
	}

	// Add some spacing
	result.WriteString("\n")

	return result.String()
}

func RenderBannerDefault() string {
	return RenderBanner(BannerConfig{
		Width:       GetTerminalWidth(),
		ShowVersion: true,
		ColorScheme: "default",
	})
}