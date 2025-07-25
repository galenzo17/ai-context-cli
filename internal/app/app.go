package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ai-context-cli/internal/feedback"
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
	
	// Feedback system
	spinner      feedback.SpinnerModel
	progress     feedback.ProgressModel
	toastManager feedback.ToastManager
	loadingState LoadingState
}

// LoadingState represents different loading states
type LoadingState int

const (
	StateMenu LoadingState = iota
	StateScanning
	StateProcessing
	StateComplete
)

// SimulateOperationMsg is sent to simulate different operations
type SimulateOperationMsg struct {
	Operation string
	Step      int
	Total     int
}

// ProgressUpdateMsg is sent to update progress
type ProgressUpdateMsg struct {
	Current int
	Total   int
	Message string
}

// OperationCompleteMsg is sent when an operation completes
type OperationCompleteMsg struct {
	Success bool
	Message string
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
		selected:     make(map[int]struct{}),
		showingHelp:  false,
		helpForItem:  -1,
		spinner:      feedback.NewSpinner("Loading..."),
		progress:     feedback.NewProgress(0, ""),
		toastManager: feedback.NewToastManager(),
		loadingState: StateMenu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	// Update spinner
	spinner, spinnerCmd := m.spinner.Update(msg)
	m.spinner = spinner
	if spinnerCmd != nil {
		cmds = append(cmds, spinnerCmd)
	}
	
	// Update toast manager
	toastManager, toastCmd := m.toastManager.Update(msg)
	m.toastManager = toastManager
	if toastCmd != nil {
		cmds = append(cmds, toastCmd)
	}
	
	switch msg := msg.(type) {
	case SimulateOperationMsg:
		return m.handleSimulateOperation(msg)
	case ProgressUpdateMsg:
		m.progress = m.progress.SetProgress(msg.Current).SetMessage(msg.Message)
		if msg.Current < msg.Total {
			// Continue simulation
			return m, m.simulateProgressStep(msg.Current+1, msg.Total, msg.Message)
		} else {
			// Operation complete
			return m, m.completeOperation(true, "Operation completed successfully!")
		}
	case OperationCompleteMsg:
		m.loadingState = StateComplete
		m.spinner = m.spinner.Stop()
		
		// Show toast notification
		var toastType feedback.ToastType
		if msg.Success {
			toastType = feedback.ToastSuccess
		} else {
			toastType = feedback.ToastError
		}
		
		toastManager, toastCmd := m.toastManager.AddToast(msg.Message, toastType)
		m.toastManager = toastManager
		
		// Reset to menu after showing result
		return m, tea.Batch(toastCmd, m.resetToMenuAfterDelay())
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
			
			// Only handle menu actions if not in loading state
			if m.loadingState == StateMenu {
				return m.handleMenuAction(m.cursor)
			}
		case "r":
			// Reset to menu (used internally after operations)
			m.loadingState = StateMenu
			m.spinner = m.spinner.Stop()
			m.progress = feedback.NewProgress(0, "")
		}
	}
	
	return m, tea.Batch(cmds...)
}

// handleMenuAction processes menu item selection
func (m Model) handleMenuAction(index int) (Model, tea.Cmd) {
	switch index {
	case 0: // Add Context (All)
		m.loadingState = StateScanning
		m.spinner = m.spinner.SetMessage("Scanning project files...").Start()
		m.progress = feedback.NewProgress(100, "Adding context to all files")
		return m, tea.Batch(
			m.spinner.InitSpinner(),
			m.simulateFileScanning(),
		)
	case 1: // Add Context (Folder)
		m.loadingState = StateScanning
		m.spinner = m.spinner.SetMessage("Selecting folder...").Start()
		m.progress = feedback.NewProgress(50, "Adding context to folder")
		return m, tea.Batch(
			m.spinner.InitSpinner(),
			m.simulateFolderSelection(),
		)
	case 2: // Context Before
		m.loadingState = StateProcessing
		m.spinner = m.spinner.SetMessage("Preparing context preview...").Start()
		return m, tea.Batch(
			m.spinner.InitSpinner(),
			m.simulateContextPreview(),
		)
	case 3: // Select Model
		m.loadingState = StateProcessing
		m.spinner = m.spinner.SetMessage("Loading available models...").Start()
		return m, tea.Batch(
			m.spinner.InitSpinner(),
			m.simulateModelLoading(),
		)
	default:
		return m, nil
	}
}

// handleSimulateOperation processes simulation messages
func (m Model) handleSimulateOperation(msg SimulateOperationMsg) (Model, tea.Cmd) {
	switch msg.Operation {
	case "fileScanning":
		m.progress = m.progress.SetProgress(msg.Step).SetMessage(fmt.Sprintf("Scanned %d files", msg.Step))
		if msg.Step < msg.Total {
			return m, m.simulateProgressStep(msg.Step+1, msg.Total, "fileScanning")
		} else {
			return m, m.completeOperation(true, "Successfully added context from all files!")
		}
	case "folderSelection":
		m.progress = m.progress.SetProgress(msg.Step).SetMessage(fmt.Sprintf("Processing folder %d/%d", msg.Step, msg.Total))
		if msg.Step < msg.Total {
			return m, m.simulateProgressStep(msg.Step+1, msg.Total, "folderSelection")
		} else {
			return m, m.completeOperation(true, "Folder context added successfully!")
		}
	default:
		return m, m.completeOperation(true, "Operation completed!")
	}
}

// Simulation commands
func (m Model) simulateFileScanning() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return SimulateOperationMsg{Operation: "fileScanning", Step: 1, Total: 100}
	})
}

func (m Model) simulateFolderSelection() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return SimulateOperationMsg{Operation: "folderSelection", Step: 1, Total: 50}
	})
}

func (m Model) simulateContextPreview() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return OperationCompleteMsg{Success: true, Message: "Context preview ready!"}
	})
}

func (m Model) simulateModelLoading() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return OperationCompleteMsg{Success: true, Message: "Models loaded successfully!"}
	})
}

func (m Model) simulateProgressStep(step, total int, operation string) tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return SimulateOperationMsg{Operation: operation, Step: step, Total: total}
	})
}

func (m Model) completeOperation(success bool, message string) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return OperationCompleteMsg{Success: success, Message: message}
	})
}

func (m Model) resetToMenuAfterDelay() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}} // Reset signal
	})
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
	
	// Always show toasts at the top
	if toastView := m.toastManager.View(); toastView != "" {
		centeredToast := centerText(toastView, 100)
		result.WriteString(centeredToast)
		result.WriteString("\n\n")
	}
	
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
	
	// Show loading state interface
	if m.loadingState != StateMenu {
		return result.String() + m.renderLoadingView()
	}
	
	return result.String() + m.renderBaseView()
}

// renderLoadingView renders the loading interface
func (m Model) renderLoadingView() string {
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
	result.WriteString("\n")
	
	// Show spinner if active
	if spinnerView := m.spinner.View(); spinnerView != "" {
		centeredSpinner := centerText(spinnerView, 100)
		result.WriteString(centeredSpinner)
		result.WriteString("\n\n")
	}
	
	// Show progress bar if operation has progress
	if m.loadingState == StateScanning && m.progress.Percentage() > 0 {
		progressView := m.progress.View()
		centeredProgress := centerText(progressView, 100)
		result.WriteString(centeredProgress)
		result.WriteString("\n\n")
	}
	
	// Loading instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "⏳ Loading... Press Ctrl+C to cancel"
	centeredInstructions := centerText(instructionStyle.Render(instructions), 100)
	result.WriteString(centeredInstructions)
	
	return result.String()
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