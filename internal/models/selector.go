package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectorModel handles the model selection interface
type SelectorModel struct {
	registry     *ModelRegistry
	models       []AIModel
	cursor       int
	selected     int
	showDetails  bool
	showConfig   bool
	showTesting  bool
	
	// Filter and view options
	filterProvider   string
	filterCapability string
	viewMode        ViewMode
	
	// Dimensions
	width  int
	height int
	
	// Styles
	styles SelectorStyles
}

// ViewMode represents different ways to display models
type ViewMode int

const (
	ViewModeList ViewMode = iota
	ViewModeGrid
	ViewModeComparison
)

// SelectorStyles contains styling for the model selector
type SelectorStyles struct {
	TitleStyle       lipgloss.Style
	HeaderStyle      lipgloss.Style
	ModelStyle       lipgloss.Style
	SelectedStyle    lipgloss.Style
	ProviderStyle    lipgloss.Style
	StatusStyle      lipgloss.Style
	DescriptionStyle lipgloss.Style
	KeyStyle         lipgloss.Style
	HelpStyle        lipgloss.Style
	BorderStyle      lipgloss.Style
}

// NewSelectorModel creates a new model selector
func NewSelectorModel() *SelectorModel {
	registry := NewModelRegistry()
	
	return &SelectorModel{
		registry:    registry,
		models:      registry.GetAllModels(),
		cursor:      0,
		selected:    -1,
		showDetails: false,
		showConfig:  false,
		showTesting: false,
		viewMode:    ViewModeList,
		styles:      newSelectorStyles(),
	}
}

// newSelectorStyles creates the default styles for the selector
func newSelectorStyles() SelectorStyles {
	return SelectorStyles{
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7C3AED")).
			Padding(0, 1),
		
		HeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1),
		
		ModelStyle: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			MarginBottom(1),
		
		SelectedStyle: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Background(lipgloss.Color("#1A1A2E")).
			MarginBottom(1),
		
		ProviderStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true),
		
		StatusStyle: lipgloss.NewStyle().
			Bold(true),
		
		DescriptionStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Width(60),
		
		KeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true),
		
		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1),
		
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(1, 2),
	}
}

// Init initializes the model selector
func (m *SelectorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the model selector
func (m *SelectorModel) Update(msg tea.Msg) (*SelectorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
			
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			
		case "down", "j":
			if m.cursor < len(m.models)-1 {
				m.cursor++
			}
			
		case "enter", " ":
			if m.cursor < len(m.models) {
				m.selected = m.cursor
				return m, m.selectModel()
			}
			
		case "d":
			m.showDetails = !m.showDetails
			
		case "c":
			m.showConfig = !m.showConfig
			
		case "t":
			if m.cursor < len(m.models) {
				m.showTesting = true
				return m, m.testConnection()
			}
			
		case "v":
			m.cycleViewMode()
			
		case "f":
			// Filter functionality (to be implemented)
			
		case "r":
			// Refresh models
			m.models = m.registry.GetAllModels()
		}
	}
	
	return m, nil
}

// selectModel returns a command to select the current model
func (m *SelectorModel) selectModel() tea.Cmd {
	return func() tea.Msg {
		return ModelSelectedMsg{
			Model: m.models[m.cursor],
		}
	}
}

// testConnection returns a command to test model connection
func (m *SelectorModel) testConnection() tea.Cmd {
	return func() tea.Msg {
		model := m.models[m.cursor]
		
		// Create connection tester and test the model
		tester := NewConnectionTester()
		result := tester.TestConnection(model)
		
		// Update the registry with the test result
		status := StatusError
		if result.Success {
			status = StatusConnected
		} else {
			status = StatusDisconnected
		}
		
		m.registry.UpdateModelStatus(model.ID, status, &result.Timestamp)
		
		// Refresh models list to show updated status
		m.models = m.registry.GetAllModels()
		
		return ModelTestResultMsg{
			ModelID: model.ID,
			Result:  result,
		}
	}
}

// cycleViewMode cycles through different view modes
func (m *SelectorModel) cycleViewMode() {
	switch m.viewMode {
	case ViewModeList:
		m.viewMode = ViewModeGrid
	case ViewModeGrid:
		m.viewMode = ViewModeComparison
	case ViewModeComparison:
		m.viewMode = ViewModeList
	}
}

// View renders the model selector
func (m *SelectorModel) View() string {
	var sections []string
	
	// Title
	title := m.styles.TitleStyle.Render("🤖 AI Model Selection")
	sections = append(sections, title)
	
	// Header with stats
	header := m.renderHeader()
	sections = append(sections, header)
	
	// Model list based on view mode
	switch m.viewMode {
	case ViewModeList:
		modelList := m.renderListView()
		sections = append(sections, modelList)
	case ViewModeGrid:
		modelGrid := m.renderGridView()
		sections = append(sections, modelGrid)
	case ViewModeComparison:
		comparison := m.renderComparisonView()
		sections = append(sections, comparison)
	}
	
	// Details panel if showing details
	if m.showDetails && m.cursor < len(m.models) {
		details := m.renderDetailsPanel()
		sections = append(sections, details)
	}
	
	// Configuration panel if showing config
	if m.showConfig && m.cursor < len(m.models) {
		config := m.renderConfigPanel()
		sections = append(sections, config)
	}
	
	// Help text
	help := m.renderHelp()
	sections = append(sections, help)
	
	return strings.Join(sections, "\n\n")
}

// renderHeader renders the header with model count and filters
func (m *SelectorModel) renderHeader() string {
	var parts []string
	
	totalModels := len(m.models)
	parts = append(parts, fmt.Sprintf("Total Models: %d", totalModels))
	
	// Provider filter
	if m.filterProvider != "" {
		parts = append(parts, fmt.Sprintf("Provider: %s", m.filterProvider))
	}
	
	// View mode
	var viewModeStr string
	switch m.viewMode {
	case ViewModeList:
		viewModeStr = "List"
	case ViewModeGrid:
		viewModeStr = "Grid"
	case ViewModeComparison:
		viewModeStr = "Comparison"
	}
	parts = append(parts, fmt.Sprintf("View: %s", viewModeStr))
	
	return m.styles.HeaderStyle.Render(strings.Join(parts, " | "))
}

// renderListView renders models in a list format
func (m *SelectorModel) renderListView() string {
	var items []string
	
	for i, model := range m.models {
		var style lipgloss.Style
		if i == m.cursor {
			style = m.styles.SelectedStyle
		} else {
			style = m.styles.ModelStyle
		}
		
		// Model info
		name := fmt.Sprintf("%s (%s)", model.Name, model.Provider)
		status := m.renderStatusBadge(model.Status)
		cost := fmt.Sprintf("$%.4f/1K", model.CostPer1K)
		tokens := fmt.Sprintf("%dK tokens", model.MaxTokens/1000)
		
		// Build the item content
		content := fmt.Sprintf("%s %s\n%s\nCost: %s | Max: %s",
			name, status, model.Description, cost, tokens)
		
		// Add selection indicator
		if i == m.cursor {
			content = "► " + content
		} else {
			content = "  " + content
		}
		
		items = append(items, style.Render(content))
	}
	
	return strings.Join(items, "\n")
}

// renderGridView renders models in a grid format
func (m *SelectorModel) renderGridView() string {
	// For now, fallback to list view (grid would require more complex layout)
	return m.renderListView()
}

// renderComparisonView renders models in a comparison table
func (m *SelectorModel) renderComparisonView() string {
	if len(m.models) == 0 {
		return "No models available for comparison"
	}
	
	var table []string
	
	// Header
	header := fmt.Sprintf("%-20s %-12s %-10s %-10s %-10s",
		"Model", "Provider", "Cost/1K", "Max Tokens", "Status")
	table = append(table, m.styles.HeaderStyle.Render(header))
	
	// Separator
	table = append(table, strings.Repeat("─", 70))
	
	// Rows
	for i, model := range m.models {
		var style lipgloss.Style
		if i == m.cursor {
			style = m.styles.SelectedStyle.Copy().Border(lipgloss.NormalBorder(), false)
		} else {
			style = lipgloss.NewStyle()
		}
		
		status := string(model.Status)
		if model.Status == StatusUnknown {
			status = "untested"
		}
		
		row := fmt.Sprintf("%-20s %-12s $%-9.4f %-10s %-10s",
			model.Name,
			model.Provider,
			model.CostPer1K,
			fmt.Sprintf("%dK", model.MaxTokens/1000),
			status)
		
		if i == m.cursor {
			row = "► " + row
		} else {
			row = "  " + row
		}
		
		table = append(table, style.Render(row))
	}
	
	return strings.Join(table, "\n")
}

// renderDetailsPanel renders detailed information about the current model
func (m *SelectorModel) renderDetailsPanel() string {
	model := m.models[m.cursor]
	
	var details []string
	details = append(details, m.styles.HeaderStyle.Render("📋 Model Details"))
	
	// Basic info
	details = append(details, fmt.Sprintf("Name: %s", model.Name))
	details = append(details, fmt.Sprintf("Provider: %s", model.Provider))
	details = append(details, fmt.Sprintf("ID: %s", model.ID))
	details = append(details, fmt.Sprintf("Description: %s", model.Description))
	
	// Capabilities
	details = append(details, "Capabilities:")
	for _, cap := range model.Capabilities {
		details = append(details, fmt.Sprintf("  • %s", cap))
	}
	
	// Costs and limits
	details = append(details, fmt.Sprintf("Cost per 1K tokens: $%.4f", model.CostPer1K))
	details = append(details, fmt.Sprintf("Maximum tokens: %d", model.MaxTokens))
	
	// Status
	status := m.renderStatusBadge(model.Status)
	details = append(details, fmt.Sprintf("Status: %s", status))
	
	if model.LastTested != nil {
		details = append(details, fmt.Sprintf("Last tested: %s", 
			model.LastTested.Format("2006-01-02 15:04:05")))
	}
	
	return m.styles.BorderStyle.Render(strings.Join(details, "\n"))
}

// renderConfigPanel renders the API configuration panel
func (m *SelectorModel) renderConfigPanel() string {
	model := m.models[m.cursor]
	
	var config []string
	config = append(config, m.styles.HeaderStyle.Render("⚙️  API Configuration"))
	
	config = append(config, fmt.Sprintf("Base URL: %s", model.APIConfig.BaseURL))
	
	// Mask API key for security
	apiKey := model.APIConfig.APIKey
	if apiKey != "" {
		if len(apiKey) > 10 {
			apiKey = apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		} else {
			apiKey = "***"
		}
	} else {
		apiKey = "Not configured"
	}
	config = append(config, fmt.Sprintf("API Key: %s", apiKey))
	
	config = append(config, fmt.Sprintf("Timeout: %s", model.APIConfig.Timeout))
	config = append(config, fmt.Sprintf("Retry count: %d", model.APIConfig.RetryCount))
	
	return m.styles.BorderStyle.Render(strings.Join(config, "\n"))
}

// renderStatusBadge renders a status badge for a model
func (m *SelectorModel) renderStatusBadge(status ConnectionStatus) string {
	var color lipgloss.Color
	var text string
	
	switch status {
	case StatusConnected:
		color = lipgloss.Color("#10B981")
		text = "✓ Connected"
	case StatusDisconnected:
		color = lipgloss.Color("#F59E0B")
		text = "⚠ Disconnected"
	case StatusError:
		color = lipgloss.Color("#EF4444")
		text = "✗ Error"
	default:
		color = lipgloss.Color("#6B7280")
		text = "? Unknown"
	}
	
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(text)
}

// renderHelp renders the help text
func (m *SelectorModel) renderHelp() string {
	help := []string{
		"Navigation: " + m.styles.KeyStyle.Render("↑/k") + " up, " + 
			m.styles.KeyStyle.Render("↓/j") + " down, " + 
			m.styles.KeyStyle.Render("enter") + " select",
		"Actions: " + m.styles.KeyStyle.Render("d") + " details, " + 
			m.styles.KeyStyle.Render("c") + " config, " + 
			m.styles.KeyStyle.Render("t") + " test, " + 
			m.styles.KeyStyle.Render("v") + " view mode",
		"Other: " + m.styles.KeyStyle.Render("r") + " refresh, " + 
			m.styles.KeyStyle.Render("q/esc") + " quit",
	}
	
	return m.styles.HelpStyle.Render(strings.Join(help, "\n"))
}

// Message types for model selection
type ModelSelectedMsg struct {
	Model AIModel
}

type ModelTestResultMsg struct {
	ModelID string
	Result  TestResult
}