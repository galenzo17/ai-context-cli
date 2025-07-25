package preview

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ai-context-cli/internal/context"
)

// ContextPreviewModel represents the context preview interface
type ContextPreviewModel struct {
	contextResult *context.ContextResult
	scanResult    *context.ScanResult
	
	// Display options
	showFullContent bool
	currentSection  int
	editMode        bool
	templateMode    bool
	currentTemplate int
	
	// UI state
	width        int
	height       int
	cursor       int
	viewport     ViewportInfo
	errorMessage string
	
	// Edit state
	editingContent string
	originalContent string
	
	// Available templates
	templates []ContextTemplate
}

// ViewportInfo tracks what's currently visible
type ViewportInfo struct {
	offset int
	size   int
}

// ContextTemplate represents a predefined context template
type ContextTemplate struct {
	Name        string
	Description string
	Template    string
	Icon        string
}

// PreviewMsg represents messages for the preview system
type PreviewMsg struct {
	Type string
	Data interface{}
}

// TokenEstimate represents token count estimation
type TokenEstimate struct {
	Characters int
	Words      int
	Tokens     int
	Cost       float64
}

// NewContextPreviewModel creates a new context preview model
func NewContextPreviewModel(contextResult *context.ContextResult, scanResult *context.ScanResult) *ContextPreviewModel {
	templates := getDefaultTemplates()
	
	return &ContextPreviewModel{
		contextResult:   contextResult,
		scanResult:      scanResult,
		width:          80,
		height:         20,
		templates:      templates,
		currentSection: 0,
		viewport: ViewportInfo{
			offset: 0,
			size:   15,
		},
	}
}

// getDefaultTemplates returns predefined context templates
func getDefaultTemplates() []ContextTemplate {
	return []ContextTemplate{
		{
			Name:        "Development Focus",
			Description: "Optimized for code development and debugging",
			Template:    "development",
			Icon:        "💻",
		},
		{
			Name:        "Documentation",
			Description: "Focused on generating documentation",
			Template:    "documentation", 
			Icon:        "📚",
		},
		{
			Name:        "Code Review",
			Description: "Structured for code review and analysis",
			Template:    "review",
			Icon:        "🔍",
		},
		{
			Name:        "Bug Analysis",
			Description: "Targeted for debugging and issue resolution",
			Template:    "debug",
			Icon:        "🐛",
		},
		{
			Name:        "Full Context",
			Description: "Complete project context with all details",
			Template:    "full",
			Icon:        "📋",
		},
	}
}

// Update handles preview messages and key events
func (m *ContextPreviewModel) Update(msg tea.Msg) (*ContextPreviewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()
	case PreviewMsg:
		return m.handlePreviewMsg(msg)
	}
	
	return m, nil
}

// handleKeyPress processes keyboard input
func (m *ContextPreviewModel) handleKeyPress(msg tea.KeyMsg) (*ContextPreviewModel, tea.Cmd) {
	if m.editMode {
		return m.handleEditMode(msg)
	}
	
	if m.templateMode {
		return m.handleTemplateMode(msg)
	}
	
	switch msg.String() {
	case "esc":
		// Exit preview mode
		return m, m.exitPreview()
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.updateViewport()
		}
	case "down", "j":
		if m.cursor < len(m.contextResult.Sections)-1 {
			m.cursor++
			m.updateViewport()
		}
	case "left", "h":
		if m.currentSection > 0 {
			m.currentSection--
		}
	case "right", "l":
		if m.currentSection < len(m.contextResult.Sections)-1 {
			m.currentSection++
		}
	case "enter", " ":
		m.showFullContent = !m.showFullContent
	case "e":
		// Enter edit mode
		m.editMode = true
		if m.currentSection < len(m.contextResult.Sections) {
			m.originalContent = m.contextResult.Sections[m.currentSection].Content
			m.editingContent = m.originalContent
		}
	case "t":
		// Enter template mode
		m.templateMode = true
		m.currentTemplate = 0
	case "r":
		// Refresh context
		return m, m.refreshContext()
	case "s":
		// Save current context
		return m, m.saveContext()
	case "home":
		m.cursor = 0
		m.currentSection = 0
		m.updateViewport()
	case "end":
		m.cursor = len(m.contextResult.Sections) - 1
		m.currentSection = len(m.contextResult.Sections) - 1
		m.updateViewport()
	}
	
	return m, nil
}

// handleEditMode processes input in edit mode
func (m *ContextPreviewModel) handleEditMode(msg tea.KeyMsg) (*ContextPreviewModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel edit
		m.editMode = false
		m.editingContent = ""
		m.originalContent = ""
	case "ctrl+s":
		// Save edit
		if m.currentSection < len(m.contextResult.Sections) {
			m.contextResult.Sections[m.currentSection].Content = m.editingContent
		}
		m.editMode = false
		m.editingContent = ""
		m.originalContent = ""
	case "ctrl+z":
		// Undo edit
		m.editingContent = m.originalContent
	default:
		// Handle text input (simplified for demo)
		if len(msg.String()) == 1 {
			m.editingContent += msg.String()
		} else if msg.String() == "backspace" {
			if len(m.editingContent) > 0 {
				m.editingContent = m.editingContent[:len(m.editingContent)-1]
			}
		}
	}
	
	return m, nil
}

// handleTemplateMode processes input in template selection mode
func (m *ContextPreviewModel) handleTemplateMode(msg tea.KeyMsg) (*ContextPreviewModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel template selection
		m.templateMode = false
	case "up", "k":
		if m.currentTemplate > 0 {
			m.currentTemplate--
		}
	case "down", "j":
		if m.currentTemplate < len(m.templates)-1 {
			m.currentTemplate++
		}
	case "enter", " ":
		// Apply selected template
		return m, m.applyTemplate(m.templates[m.currentTemplate])
	}
	
	return m, nil
}

// updateViewport adjusts the viewport to keep cursor visible
func (m *ContextPreviewModel) updateViewport() {
	m.viewport.size = m.height - 8 // Reserve space for header and footer
	
	// Adjust offset to keep cursor visible
	if m.cursor < m.viewport.offset {
		m.viewport.offset = m.cursor
	} else if m.cursor >= m.viewport.offset+m.viewport.size {
		m.viewport.offset = m.cursor - m.viewport.size + 1
	}
	
	// Ensure offset doesn't go negative
	if m.viewport.offset < 0 {
		m.viewport.offset = 0
	}
}

// View renders the context preview interface
func (m *ContextPreviewModel) View() string {
	var result strings.Builder
	
	// Header
	result.WriteString(m.renderHeader())
	result.WriteString("\n\n")
	
	// Error message
	if m.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)
		result.WriteString(errorStyle.Render("⚠️ " + m.errorMessage))
		result.WriteString("\n\n")
	}
	
	// Content based on mode
	if m.editMode {
		result.WriteString(m.renderEditMode())
	} else if m.templateMode {
		result.WriteString(m.renderTemplateMode())
	} else {
		result.WriteString(m.renderContextPreview())
	}
	
	// Footer
	result.WriteString("\n\n")
	result.WriteString(m.renderFooter())
	
	return result.String()
}

// renderHeader renders the header with context summary
func (m *ContextPreviewModel) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width)
	
	// Calculate token estimate
	estimate := m.calculateTokenEstimate()
	
	header := fmt.Sprintf("📋 Context Preview - %s | %d sections | ~%s tokens | ~$%.4f",
		m.contextResult.ProjectName,
		len(m.contextResult.Sections),
		formatNumber(estimate.Tokens),
		estimate.Cost)
	
	return headerStyle.Render(header)
}

// renderContextPreview renders the main context preview
func (m *ContextPreviewModel) renderContextPreview() string {
	var result strings.Builder
	
	if len(m.contextResult.Sections) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
		result.WriteString(emptyStyle.Render("No context sections available"))
		return result.String()
	}
	
	// Section navigation
	sectionNavStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		Bold(true)
	
	navText := fmt.Sprintf("Section %d/%d: %s", 
		m.currentSection+1, 
		len(m.contextResult.Sections),
		m.contextResult.Sections[m.currentSection].Title)
	result.WriteString(sectionNavStyle.Render(navText))
	result.WriteString("\n\n")
	
	// Section content
	section := m.contextResult.Sections[m.currentSection]
	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#374151")).
		Width(m.width-4).
		Padding(1, 2)
	
	content := section.Content
	if !m.showFullContent && len(content) > 500 {
		content = content[:500] + "...\n\nPress ENTER to show full content"
	}
	
	result.WriteString(contentStyle.Render(content))
	
	// Section metadata
	if len(section.Files) > 0 {
		result.WriteString("\n\n")
		metadataStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
		
		fileList := strings.Join(section.Files, ", ")
		if len(fileList) > 60 {
			fileList = fileList[:60] + "..."
		}
		result.WriteString(metadataStyle.Render(fmt.Sprintf("Files: %s", fileList)))
	}
	
	return result.String()
}

// renderEditMode renders the edit interface
func (m *ContextPreviewModel) renderEditMode() string {
	var result strings.Builder
	
	editHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F59E0B"))
	
	result.WriteString(editHeaderStyle.Render("✏️ Edit Mode - Section: " + m.contextResult.Sections[m.currentSection].Title))
	result.WriteString("\n\n")
	
	// Edit area
	editStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F59E0B")).
		Width(m.width-4).
		Height(m.height-10).
		Padding(1)
	
	result.WriteString(editStyle.Render(m.editingContent))
	
	// Edit instructions
	result.WriteString("\n\n")
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	result.WriteString(instructionStyle.Render("Ctrl+S: Save • Ctrl+Z: Undo • ESC: Cancel"))
	
	return result.String()
}

// renderTemplateMode renders the template selection interface
func (m *ContextPreviewModel) renderTemplateMode() string {
	var result strings.Builder
	
	templateHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#10B981"))
	
	result.WriteString(templateHeaderStyle.Render("🎨 Template Selection"))
	result.WriteString("\n\n")
	
	// Template list
	for i, template := range m.templates {
		isSelected := i == m.currentTemplate
		
		var templateStyle lipgloss.Style
		if isSelected {
			templateStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#10B981")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 1)
		} else {
			templateStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#374151")).
				Padding(0, 1)
		}
		
		templateText := fmt.Sprintf("%s %s - %s", template.Icon, template.Name, template.Description)
		result.WriteString(templateStyle.Render(templateText))
		result.WriteString("\n")
	}
	
	return result.String()
}

// renderFooter renders the footer with controls and statistics
func (m *ContextPreviewModel) renderFooter() string {
	var result strings.Builder
	
	// Statistics
	estimate := m.calculateTokenEstimate()
	statsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B82F6")).
		BorderTop(true).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1, 0)
	
	stats := fmt.Sprintf("📊 %s chars | %s words | ~%s tokens | Size: %s | Files: %d",
		formatNumber(estimate.Characters),
		formatNumber(estimate.Words), 
		formatNumber(estimate.Tokens),
		context.FormatSize(m.contextResult.TotalSize),
		m.contextResult.TotalFiles)
	
	result.WriteString(statsStyle.Render(stats))
	result.WriteString("\n")
	
	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	var instructions string
	if m.editMode {
		instructions = "Edit mode active"
	} else if m.templateMode {
		instructions = "↑↓: select template • Enter: apply • ESC: cancel"
	} else {
		instructions = "←→: navigate sections • Enter: toggle full view • E: edit • T: templates • S: save • R: refresh • ESC: exit"
	}
	
	result.WriteString(instructionStyle.Render(instructions))
	
	return result.String()
}

// calculateTokenEstimate estimates token count and cost
func (m *ContextPreviewModel) calculateTokenEstimate() TokenEstimate {
	var totalChars, totalWords int
	
	for _, section := range m.contextResult.Sections {
		totalChars += len(section.Content)
		totalWords += len(strings.Fields(section.Content))
	}
	
	// Rough token estimation (1 token ≈ 4 characters for GPT models)
	estimatedTokens := totalChars / 4
	
	// Rough cost estimation (assuming GPT-4 pricing)
	costPer1KTokens := 0.03 // $0.03 per 1K tokens (input)
	estimatedCost := float64(estimatedTokens) / 1000.0 * costPer1KTokens
	
	return TokenEstimate{
		Characters: totalChars,
		Words:      totalWords,
		Tokens:     estimatedTokens,
		Cost:       estimatedCost,
	}
}

// refreshContext refreshes the context data
func (m *ContextPreviewModel) refreshContext() tea.Cmd {
	return func() tea.Msg {
		return PreviewMsg{
			Type: "refresh_requested",
			Data: nil,
		}
	}
}

// saveContext saves the current context
func (m *ContextPreviewModel) saveContext() tea.Cmd {
	return func() tea.Msg {
		return PreviewMsg{
			Type: "save_requested",
			Data: m.contextResult,
		}
	}
}

// applyTemplate applies a selected template
func (m *ContextPreviewModel) applyTemplate(template ContextTemplate) tea.Cmd {
	return func() tea.Msg {
		return PreviewMsg{
			Type: "template_applied",
			Data: template,
		}
	}
}

// exitPreview exits the preview mode
func (m *ContextPreviewModel) exitPreview() tea.Cmd {
	return func() tea.Msg {
		return PreviewMsg{
			Type: "exit_preview",
			Data: nil,
		}
	}
}

// GetContextResult returns the current context result
func (m *ContextPreviewModel) GetContextResult() *context.ContextResult {
	return m.contextResult
}

// SetSize updates the preview dimensions
func (m *ContextPreviewModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.updateViewport()
}

// handlePreviewMsg processes preview-specific messages
func (m *ContextPreviewModel) handlePreviewMsg(msg PreviewMsg) (*ContextPreviewModel, tea.Cmd) {
	switch msg.Type {
	case "context_updated":
		if contextResult, ok := msg.Data.(*context.ContextResult); ok {
			m.contextResult = contextResult
		}
	case "template_applied":
		m.templateMode = false
		// Template application logic would go here
	}
	
	return m, nil
}

// formatNumber formats a number with K/M suffixes
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
}