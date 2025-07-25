package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const version = "v0.1.0"

var logoArt = []string{
	"    █████╗ ██╗     ██████╗ ██████╗ ███╗   ██╗████████╗███████╗██╗  ██╗████████╗",
	"   ██╔══██╗██║    ██╔════╝██╔═══██╗████╗  ██║╚══██╔══╝██╔════╝╚██╗██╔╝╚══██╔══╝",
	"   ███████║██║    ██║     ██║   ██║██╔██╗ ██║   ██║   █████╗   ╚███╔╝    ██║   ",
	"   ██╔══██║██║    ██║     ██║   ██║██║╚██╗██║   ██║   ██╔══╝   ██╔██╗    ██║   ",
	"   ██║  ██║██║    ╚██████╗╚██████╔╝██║ ╚████║   ██║   ███████╗██╔╝ ██╗   ██║   ",
	"   ╚═╝  ╚═╝╚═╝     ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝   ╚═╝   ",
	"",
	"                           ██████╗██╗     ██╗",
	"                          ██╔════╝██║     ██║",
	"                          ██║     ██║     ██║",
	"                          ██║     ██║     ██║",
	"                          ╚██████╗███████╗██║",
	"                           ╚═════╝╚══════╝╚═╝",
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

	// Create gradient colors
	primaryColor := lipgloss.Color("#7D56F4")   // Purple
	accentColor := lipgloss.Color("#10B981")    // Green

	// Create styles
	logoStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)

	versionStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Italic(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)

	// Check if terminal is wide enough for full logo
	logoWidth := 78 // Width of the ASCII art
	if config.Width < logoWidth+4 {
		// Use compact version for narrow terminals
		compactLogo := []string{
			"╔═══════════════════════════════════╗",
			"║           AI CONTEXT CLI          ║",
			"║     Advanced Context Engine       ║",
			"╚═══════════════════════════════════╝",
		}
		
		for _, line := range compactLogo {
			centeredLine := centerText(logoStyle.Render(line), config.Width)
			result.WriteString(centeredLine + "\n")
		}
	} else {
		// Render full ASCII art logo
		for _, line := range logoArt {
			if line == "" {
				result.WriteString("\n")
				continue
			}
			centeredLine := centerText(logoStyle.Render(line), config.Width)
			result.WriteString(centeredLine + "\n")
		}
	}

	// Add subtitle
	result.WriteString("\n")
	subtitle := "Advanced Context Engineering for AI Models"
	centeredSubtitle := centerText(subtitleStyle.Render(subtitle), config.Width)
	result.WriteString(centeredSubtitle + "\n")

	// Add version if requested
	if config.ShowVersion {
		result.WriteString("\n")
		versionText := fmt.Sprintf("Version %s", version)
		centeredVersion := centerText(versionStyle.Render(versionText), config.Width)
		result.WriteString(centeredVersion + "\n")
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