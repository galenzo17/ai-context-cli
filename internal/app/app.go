package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	Title       string
	Description string
	Icon        string
	DetailHelp  string
}

type Model struct {
	menuItems    []MenuItem
	cursor       int
	selected     map[int]struct{}
	showingHelp  bool
	helpForItem  int
}

func NewModel() Model {
	return Model{
		menuItems: []MenuItem{
			{
				Title:       "📂 Add Context to All Files",
				Description: "Scan entire project and add all files to AI context",
				Icon:        "📂",
				DetailHelp:  "Recursively scans your project directory and adds all code files to the AI context. Useful for giving the AI complete understanding of your project structure and codebase.",
			},
			{
				Title:       "📁 Add Context to Specific Folder",
				Description: "Choose a folder to add to AI context",
				Icon:        "📁",
				DetailHelp:  "Browse and select a specific folder to add to the AI context. This allows you to focus the AI's attention on a particular part of your project.",
			},
			{
				Title:       "📋 Preview Context Before Sending",
				Description: "Review and edit context before AI interaction",
				Icon:        "📋",
				DetailHelp:  "Shows you exactly what context will be sent to the AI model. You can review, edit, or modify the context before starting your conversation.",
			},
			{
				Title:       "🤖 Select AI Model",
				Description: "Choose and configure AI model settings",
				Icon:        "🤖",
				DetailHelp:  "Select from available AI models (GPT-4, Claude, etc.), configure API keys, and adjust model-specific settings like temperature and max tokens.",
			},
			{
				Title:       "🚪 Exit",
				Description: "Quit the application",
				Icon:        "🚪",
				DetailHelp:  "Exit the Context Engine application safely.",
			},
		},
		selected:    make(map[int]struct{}),
		showingHelp: false,
		helpForItem: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.showingHelp {
				// Close help modal
				m.showingHelp = false
				m.helpForItem = -1
				return m, nil
			}
			return m, tea.Quit
		case "esc":
			if m.showingHelp {
				m.showingHelp = false
				m.helpForItem = -1
			}
			return m, nil
		case "up", "k":
			if !m.showingHelp && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.showingHelp && m.cursor < len(m.menuItems)-1 {
				m.cursor++
			}
		case "?":
			// Show help for current item
			m.showingHelp = true
			m.helpForItem = m.cursor
		case "enter", " ":
			if m.showingHelp {
				// Close help modal
				m.showingHelp = false
				m.helpForItem = -1
				return m, nil
			}
			if m.cursor == len(m.menuItems)-1 {
				return m, tea.Quit
			}
			// Here you would implement the actual functionality
			// For now, just toggle selection for demo
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

// Helper function to center text within a given width
func centerText(text string, width int) string {
	lines := strings.Split(text, "\n")
	var centeredLines []string
	
	for _, line := range lines {
		// Handle ANSI styled text by using lipgloss.Width
		lineLen := lipgloss.Width(line)
		if lineLen >= width {
			centeredLines = append(centeredLines, line)
			continue
		}
		padding := (width - lineLen) / 2
		centeredLine := strings.Repeat(" ", padding) + line
		centeredLines = append(centeredLines, centeredLine)
	}
	
	return strings.Join(centeredLines, "\n")
}

// Helper function to create a boxed button
func (m Model) createButton(item MenuItem, index int, isSelected bool) string {
	// Define colors
	primaryColor := lipgloss.Color("#7D56F4")   // Purple
	selectedColor := lipgloss.Color("#3B82F6")  // Blue  
	normalColor := lipgloss.Color("#6B7280")    // Gray
	bgSelectedColor := lipgloss.Color("#1E1B4B") // Dark purple background
	
	// Button dimensions - wider for more info
	buttonWidth := 50
	buttonHeight := 2
	
	// Create the button content with title and description
	content := item.Title + "\n" + item.Description
	centeredText := centerText(content, buttonWidth-2) // -2 for borders
	
	// Define styles based on selection state
	var buttonStyle lipgloss.Style
	if isSelected {
		buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(selectedColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(bgSelectedColor).
			Width(buttonWidth).
			Height(buttonHeight).
			Align(lipgloss.Center).
			Bold(true)
	} else {
		buttonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(normalColor).
			Foreground(primaryColor).
			Width(buttonWidth).
			Height(buttonHeight).
			Align(lipgloss.Center).
			Bold(false)
	}
	
	return buttonStyle.Render(centeredText)
}

// Helper function to create help modal
func (m Model) createHelpModal(item MenuItem) string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3B82F6")).
		Background(lipgloss.Color("#1E1B4B")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 2).
		Width(60).
		Bold(true)
	
	content := "Help: " + item.Title + "\n\n" + item.DetailHelp + "\n\nPress ESC or Enter to close"
	return modalStyle.Render(content)
}

func (m Model) View() string {
	var result strings.Builder
	
	// If showing help modal, render it over everything
	if m.showingHelp && m.helpForItem >= 0 && m.helpForItem < len(m.menuItems) {
		// Still show the base interface but dimmed
		baseView := m.renderBaseView()
		
		// Create overlay with help modal
		helpModal := m.createHelpModal(m.menuItems[m.helpForItem])
		centeredModal := centerText(helpModal, 100)
		
		// Simple overlay - just show the modal over the base view
		result.WriteString(baseView)
		result.WriteString("\n\n")
		result.WriteString(centeredModal)
		result.WriteString("\n\n")
		
		return result.String()
	}
	
	return m.renderBaseView()
}

func (m Model) renderBaseView() string {
	var result strings.Builder
	
	// Compact banner
	bannerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)
	
	compactBanner := []string{
		"╔═══════════════════════════╗",
		"║      Context Engine       ║", 
		"╚═══════════════════════════╝",
	}
	
	for _, line := range compactBanner {
		centeredLine := centerText(bannerStyle.Render(line), 100)
		result.WriteString(centeredLine)
		result.WriteString("\n")
	}
	result.WriteString("\n") // Single line spacing after banner
	
	// Create buttons layout
	for i, item := range m.menuItems {
		isSelected := i == m.cursor
		button := m.createButton(item, i, isSelected)
		
		// Center each button
		centeredButton := centerText(button, 100)
		result.WriteString(centeredButton)
		result.WriteString("\n") // Single line spacing between buttons
	}
	
	// Add compact instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "↑↓/jk: navigate • Enter: select • ?: help • q: quit"
	centeredInstructions := centerText(instructionStyle.Render(instructions), 100)
	result.WriteString("\n")
	result.WriteString(centeredInstructions)
	
	return result.String()
}