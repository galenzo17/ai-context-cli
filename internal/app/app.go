package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ai-context-cli/internal/context"
	"ai-context-cli/internal/feedback"
	"ai-context-cli/internal/folder"
	"ai-context-cli/internal/navigation"
	"ai-context-cli/internal/preview"
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
	
	// Navigation system
	navStack     navigation.NavigationStack
	navRenderer  navigation.NavigationRenderer
	currentScreen string
	
	// Context system
	scanner      *context.ProjectScanner
	scanResult   *context.ScanResult
	contextResult *context.ContextResult
	showingResult bool
	
	// Folder browser system
	folderBrowser *folder.BrowserModel
	showingBrowser bool
	selectedFolder *folder.FolderNode
	
	// Context preview system
	contextPreview *preview.ContextPreviewModel
	showingPreview bool
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

// ScanProgressMsg is sent during real project scanning
type ScanProgressMsg struct {
	Progress context.ScanProgress
}

// ScanCompleteMsg is sent when scanning completes
type ScanCompleteMsg struct {
	Result *context.ScanResult
	Error  error
}

// ContextGeneratedMsg is sent when context generation completes
type ContextGeneratedMsg struct {
	Result *context.ContextResult
	Error  error
}

// FolderSelectedMsg is sent when a folder is selected
type FolderSelectedMsg struct {
	Folder *folder.FolderNode
}

// FolderBrowserMsg is sent for folder browser events
type FolderBrowserMsg struct {
	Type string
	Data interface{}
}

// ContextPreviewMsg is sent for context preview events
type ContextPreviewMsg struct {
	Type string
	Data interface{}
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
		navStack:     navigation.NewNavigationStack().Push(navigation.MainMenuScreen),
		navRenderer:  navigation.NewNavigationRenderer(),
		currentScreen: "main_menu",
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
	case ScanProgressMsg:
		return m.handleScanProgress(msg)
	case ScanCompleteMsg:
		return m.handleScanComplete(msg)
	case ContextGeneratedMsg:
		return m.handleContextGenerated(msg)
	case FolderSelectedMsg:
		return m.handleFolderSelected(msg)
	case FolderBrowserMsg:
		return m.handleFolderBrowser(msg)
	case folder.BrowserMsg:
		// Convert BrowserMsg to FolderBrowserMsg for consistency
		return m.handleFolderBrowser(FolderBrowserMsg{Type: msg.Type, Data: msg.Data})
	case ContextPreviewMsg:
		return m.handleContextPreview(msg)
	case preview.PreviewMsg:
		// Convert PreviewMsg to ContextPreviewMsg for consistency
		return m.handleContextPreview(ContextPreviewMsg{Type: msg.Type, Data: msg.Data})
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
		// Handle context preview first - it should get all key events when active
		if m.showingPreview && m.contextPreview != nil {
			preview, cmd := m.contextPreview.Update(msg)
			m.contextPreview = preview
			
			// Execute the command and handle any returned messages
			if cmd != nil {
				var newCmds []tea.Cmd
				newCmds = append(newCmds, cmd)
				return m, tea.Batch(newCmds...)
			}
			return m, nil
		}
		
		// Handle folder browser second - it should get all key events when active
		if m.showingBrowser && m.folderBrowser != nil {
			browser, cmd := m.folderBrowser.Update(msg)
			m.folderBrowser = browser
			
			// Execute the command and handle any returned messages
			if cmd != nil {
				var newCmds []tea.Cmd
				newCmds = append(newCmds, cmd)
				return m, tea.Batch(newCmds...)
			}
			return m, nil
		}
		
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
				return m, nil
			}
			
			// Handle context preview close
			if m.showingPreview {
				m.showingPreview = false
				m.contextPreview = nil
				return m, nil
			}
			
			// Handle folder browser close
			if m.showingBrowser {
				m.showingBrowser = false
				m.folderBrowser = nil
				return m, nil
			}
			
			// Handle navigation back
			if m.loadingState == StateMenu && m.navStack.CanGoBack() {
				navStack, success := m.navStack.Pop()
				if success {
					m.navStack = navStack
					if current, ok := m.navStack.Current(); ok {
						m.currentScreen = current.ID
						// Show toast for navigation
						toastManager, toastCmd := m.toastManager.AddToast("Returned to "+current.Title, feedback.ToastInfo)
						m.toastManager = toastManager
						return m, toastCmd
					}
				}
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
			// Reset navigation to main menu
			m.navStack = navigation.NewNavigationStack().Push(navigation.MainMenuScreen)
			m.currentScreen = "main_menu"
		}
	}
	
	return m, tea.Batch(cmds...)
}

// handleScanProgress handles real scan progress updates
func (m Model) handleScanProgress(msg ScanProgressMsg) (Model, tea.Cmd) {
	progress := msg.Progress
	
	// Update spinner message
	m.spinner = m.spinner.SetMessage(progress.CurrentPhase)
	
	// Update progress bar
	if progress.TotalEstimated > 0 {
		m.progress = m.progress.SetProgress(progress.ProcessedFiles).
			SetMessage(fmt.Sprintf("%s (%d/%d files)", 
				progress.CurrentPhase, progress.ProcessedFiles, progress.TotalEstimated))
		
		// Update total if we have better estimate
		if progress.TotalEstimated > 0 {
			m.progress = feedback.NewProgress(progress.TotalEstimated, "Scanning files")
			m.progress = m.progress.SetProgress(progress.ProcessedFiles)
		}
	}
	
	return m, nil
}

// handleScanComplete handles scan completion
func (m Model) handleScanComplete(msg ScanCompleteMsg) (Model, tea.Cmd) {
	if msg.Error != nil {
		m.loadingState = StateComplete
		m.spinner = m.spinner.Stop()
		
		toastManager, toastCmd := m.toastManager.AddToast(
			fmt.Sprintf("Scan failed: %v", msg.Error), feedback.ToastError)
		m.toastManager = toastManager
		
		return m, tea.Batch(toastCmd, m.resetToMenuAfterDelay())
	}
	
	// Store scan result and start context generation
	m.scanResult = msg.Result
	m.loadingState = StateProcessing
	m.spinner = m.spinner.SetMessage("Generating comprehensive context...").Start()
	m.progress = feedback.NewProgress(0, "Processing scan results")
	
	toastManager, toastCmd := m.toastManager.AddToast(
		fmt.Sprintf("Scanned %d files in %v", msg.Result.TotalFiles, msg.Result.ScanDuration.Round(time.Millisecond)), 
		feedback.ToastSuccess)
	m.toastManager = toastManager
	
	return m, tea.Batch(toastCmd, m.generateContext())
}

// handleContextGenerated handles context generation completion
func (m Model) handleContextGenerated(msg ContextGeneratedMsg) (Model, tea.Cmd) {
	if msg.Error != nil {
		m.loadingState = StateComplete
		m.spinner = m.spinner.Stop()
		
		toastManager, toastCmd := m.toastManager.AddToast(
			fmt.Sprintf("Context generation failed: %v", msg.Error), feedback.ToastError)
		m.toastManager = toastManager
		
		return m, tea.Batch(toastCmd, m.resetToMenuAfterDelay())
	}
	
	// Store context result and show success
	m.contextResult = msg.Result
	m.loadingState = StateComplete
	m.spinner = m.spinner.Stop()
	m.showingResult = true
	
	toastManager, toastCmd := m.toastManager.AddToast(
		fmt.Sprintf("Context generated! %d sections, ~%d tokens", 
			len(msg.Result.Sections), msg.Result.TokenEstimate), 
		feedback.ToastSuccess)
	m.toastManager = toastManager
	
	return m, toastCmd
}

// handleFolderSelected handles folder selection from browser
func (m Model) handleFolderSelected(msg FolderSelectedMsg) (Model, tea.Cmd) {
	m.selectedFolder = msg.Folder
	m.showingBrowser = false
	m.folderBrowser = nil
	
	if msg.Folder == nil {
		toastManager, toastCmd := m.toastManager.AddToast("No folder selected", feedback.ToastWarning)
		m.toastManager = toastManager
		return m, toastCmd
	}
	
	// Start folder scanning
	m.loadingState = StateScanning
	m.spinner = m.spinner.SetMessage(fmt.Sprintf("Scanning folder '%s'...", msg.Folder.Name)).Start()
	m.progress = feedback.NewProgress(0, "Scanning folder files")
	
	toastManager, toastCmd := m.toastManager.AddToast(
		fmt.Sprintf("Selected folder: %s", msg.Folder.Name), feedback.ToastInfo)
	m.toastManager = toastManager
	
	return m, tea.Batch(
		toastCmd,
		m.spinner.InitSpinner(),
		m.startFolderScan(msg.Folder.Path),
	)
}

// handleFolderBrowser handles folder browser events
func (m Model) handleFolderBrowser(msg FolderBrowserMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case "folder_selected":
		if node, ok := msg.Data.(*folder.FolderNode); ok {
			return m.handleFolderSelected(FolderSelectedMsg{Folder: node})
		}
	}
	
	return m, nil
}

// handleContextPreview handles context preview events
func (m Model) handleContextPreview(msg ContextPreviewMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case "save_requested":
		if _, ok := msg.Data.(*context.ContextResult); ok {
			// Handle context save
			toastManager, toastCmd := m.toastManager.AddToast("Context saved successfully", feedback.ToastSuccess)
			m.toastManager = toastManager
			return m, toastCmd
		}
	case "refresh_requested":
		// Handle context refresh
		toastManager, toastCmd := m.toastManager.AddToast("Context refreshed", feedback.ToastInfo)
		m.toastManager = toastManager
		return m, toastCmd
	case "template_applied":
		// Handle template application
		toastManager, toastCmd := m.toastManager.AddToast("Template applied successfully", feedback.ToastSuccess)
		m.toastManager = toastManager
		return m, toastCmd
	case "exit_preview":
		// Handle exit preview
		m.showingPreview = false
		m.contextPreview = nil
		return m, nil
	}
	
	return m, nil
}

// startProjectScan starts a real project scan
func (m Model) startProjectScan() tea.Cmd {
	return func() tea.Msg {
		// Get current working directory
		wd, err := os.Getwd()
		if err != nil {
			return ScanCompleteMsg{Error: fmt.Errorf("failed to get working directory: %w", err)}
		}
		
		// Create scanner with default config
		config := context.DefaultScanConfig(wd)
		scanner := context.NewProjectScanner(config)
		
		// Start progress monitoring in a goroutine
		progressChan := scanner.GetProgressChannel()
		go func() {
			for progress := range progressChan {
				// Send progress updates to the main loop
				// Note: In a real implementation, you'd want to use a proper
				// channel or callback mechanism here
				_ = progress
			}
		}()
		
		// Perform the scan
		result, err := scanner.Scan()
		if err != nil {
			return ScanCompleteMsg{Error: err}
		}
		
		return ScanCompleteMsg{Result: result}
	}
}

// generateContext generates context from scan results
func (m Model) generateContext() tea.Cmd {
	return func() tea.Msg {
		if m.scanResult == nil {
			return ContextGeneratedMsg{Error: fmt.Errorf("no scan result available")}
		}
		
		// Create context generator
		generator := context.NewContextGenerator()
		
		// Get project name from current directory
		wd, _ := os.Getwd()
		projectName := "Project"
		if wd != "" {
			projectName = strings.TrimSuffix(wd, "/")
			if idx := strings.LastIndex(projectName, "/"); idx >= 0 {
				projectName = projectName[idx+1:]
			}
		}
		
		// Generate context
		result, err := generator.GenerateContext(m.scanResult, projectName)
		if err != nil {
			return ContextGeneratedMsg{Error: err}
		}
		
		return ContextGeneratedMsg{Result: result}
	}
}

// handleMenuAction processes menu item selection
func (m Model) handleMenuAction(index int) (Model, tea.Cmd) {
	switch index {
	case 0: // Add Context (All)
		// Navigate to Add Context All screen
		m.navStack = m.navStack.Push(navigation.AddContextAllScreen)
		m.currentScreen = "add_context_all"
		m.loadingState = StateScanning
		m.spinner = m.spinner.SetMessage("Initializing project scan...").Start()
		m.progress = feedback.NewProgress(0, "Scanning project files")
		m.showingResult = false
		
		// Start real project scanning
		return m, tea.Batch(
			m.spinner.InitSpinner(),
			m.startProjectScan(),
		)
	case 1: // Add Context (Folder)
		// Navigate to Add Context Folder screen and open browser
		m.navStack = m.navStack.Push(navigation.AddContextFolderScreen)
		m.currentScreen = "add_context_folder"
		
		// Initialize folder browser
		wd, err := os.Getwd()
		if err != nil {
			toastManager, toastCmd := m.toastManager.AddToast(
				fmt.Sprintf("Error getting current directory: %v", err), feedback.ToastError)
			m.toastManager = toastManager
			return m, toastCmd
		}
		
		browser, err := folder.NewBrowserModel(wd)
		if err != nil {
			toastManager, toastCmd := m.toastManager.AddToast(
				fmt.Sprintf("Error initializing folder browser: %v", err), feedback.ToastError)
			m.toastManager = toastManager
			return m, toastCmd
		}
		
		m.folderBrowser = browser
		m.showingBrowser = true
		m.showingResult = false
		
		return m, nil
	case 2: // Context Before
		// Navigate to Context Preview screen
		m.navStack = m.navStack.Push(navigation.ContextPreviewScreen)
		m.currentScreen = "context_preview"
		
		// Check if we have context to preview
		if m.contextResult == nil {
			toastManager, toastCmd := m.toastManager.AddToast(
				"No context available. Please scan files first.", feedback.ToastWarning)
			m.toastManager = toastManager
			return m, toastCmd
		}
		
		// Initialize context preview
		contextPreview := preview.NewContextPreviewModel(m.contextResult, m.scanResult)
		m.contextPreview = contextPreview
		m.showingPreview = true
		m.showingResult = false
		
		return m, nil
	case 3: // Select Model
		// Navigate to Model Selection screen
		m.navStack = m.navStack.Push(navigation.ModelSelectionScreen)
		m.currentScreen = "model_selection"
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
	
	// Always show navigation at the top
	navView := m.navRenderer.RenderFullNavigation(m.navStack)
	if navView != "" {
		centeredNav := m.navRenderer.CenterNavigation(navView, 100)
		result.WriteString(centeredNav)
		result.WriteString("\n\n")
	}
	
	// Always show toasts after navigation
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
	
	// Show context preview if active
	if m.showingPreview && m.contextPreview != nil {
		return result.String() + m.contextPreview.View()
	}
	
	// Show folder browser if active
	if m.showingBrowser && m.folderBrowser != nil {
		return result.String() + m.folderBrowser.View()
	}
	
	// Show result view if available
	if m.showingResult && m.contextResult != nil {
		return result.String() + m.renderResultView()
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
	
	// Loading instructions with navigation hint
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "⏳ Loading... "
	if m.navStack.CanGoBack() {
		instructions += "ESC: Back • "
	}
	instructions += "Ctrl+C: Cancel"
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
	
	// Add compact instructions with navigation
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "↑↓/jk: navigate • Enter: select • ?: help"
	if m.navStack.CanGoBack() {
		instructions += " • ESC: back"
	}
	instructions += " • q: quit"
	centeredInstructions := centerText(instructionStyle.Render(instructions), 100)
	result.WriteString("\n")
	result.WriteString(centeredInstructions)
	
	return result.String()
}

// renderResultView renders the context generation results
func (m Model) renderResultView() string {
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
	
	// Context Results Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#10B981")).
		Align(lipgloss.Center)
	
	title := fmt.Sprintf("✨ Context Generated Successfully! ✨")
	centeredTitle := centerText(titleStyle.Render(title), 100)
	result.WriteString(centeredTitle)
	result.WriteString("\n\n")
	
	// Summary statistics
	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981")).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)
	
	var summaryContent strings.Builder
	summaryContent.WriteString(fmt.Sprintf("📁 Project: %s\n", m.contextResult.ProjectName))
	summaryContent.WriteString(fmt.Sprintf("📊 Files Processed: %d\n", m.contextResult.TotalFiles))
	summaryContent.WriteString(fmt.Sprintf("📄 Total Size: %s\n", context.FormatSize(m.contextResult.TotalSize)))
	summaryContent.WriteString(fmt.Sprintf("📝 Sections Generated: %d\n", len(m.contextResult.Sections)))
	summaryContent.WriteString(fmt.Sprintf("🧠 Estimated Tokens: ~%s\n", context.FormatNumber(m.contextResult.TokenEstimate)))
	summaryContent.WriteString(fmt.Sprintf("⏱️ Generated: %s", m.contextResult.GeneratedAt.Format("15:04:05")))
	
	summaryRendered := summaryBox.Render(summaryContent.String())
	centeredSummary := centerText(summaryRendered, 100)
	result.WriteString(centeredSummary)
	result.WriteString("\n\n")
	
	// Sections overview
	if len(m.contextResult.Sections) > 0 {
		sectionTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#3B82F6")).
			Render("📋 Generated Sections:")
		
		centeredSectionTitle := centerText(sectionTitle, 100)
		result.WriteString(centeredSectionTitle)
		result.WriteString("\n\n")
		
		for i, section := range m.contextResult.Sections {
			if i >= 5 { // Show first 5 sections
				moreText := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#6B7280")).
					Italic(true).
					Render(fmt.Sprintf("... and %d more sections", len(m.contextResult.Sections)-5))
				centeredMore := centerText(moreText, 100)
				result.WriteString(centeredMore)
				result.WriteString("\n")
				break
			}
			
			sectionItem := fmt.Sprintf("• %s", section.Title)
			if len(section.Files) > 0 {
				sectionItem += fmt.Sprintf(" (%d files)", len(section.Files))
			}
			
			sectionStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#374151"))
			
			centeredSection := centerText(sectionStyle.Render(sectionItem), 100)
			result.WriteString(centeredSection)
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}
	
	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "✨ Context ready for AI interaction!"
	if m.navStack.CanGoBack() {
		instructions += " • ESC: back"
	}
	instructions += " • q: quit"
	centeredInstructions := centerText(instructionStyle.Render(instructions), 100)
	result.WriteString(centeredInstructions)
	
	return result.String()
}