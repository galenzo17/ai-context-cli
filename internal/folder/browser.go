package folder

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BrowserModel represents the folder browser UI
type BrowserModel struct {
	tree         *FolderTree
	visibleNodes []*FolderNode
	cursor       int
	viewport     ViewportInfo
	width        int
	height       int
	showStats    bool
	confirmMode  bool
	errorMessage string
}

// ViewportInfo tracks what's currently visible
type ViewportInfo struct {
	offset int
	size   int
}

// BrowserMsg represents messages for the browser
type BrowserMsg struct {
	Type string
	Data interface{}
}

// NewBrowserModel creates a new folder browser
func NewBrowserModel(rootPath string) (*BrowserModel, error) {
	tree, err := NewFolderTree(rootPath)
	if err != nil {
		return nil, err
	}
	
	browser := &BrowserModel{
		tree:      tree,
		cursor:    0,
		width:     80,
		height:    20,
		showStats: true,
	}
	
	browser.refreshView()
	return browser, nil
}

// refreshView updates the visible nodes list
func (m *BrowserModel) refreshView() {
	m.visibleNodes = m.tree.GetVisibleNodes()
	
	// Ensure cursor is within bounds
	if m.cursor >= len(m.visibleNodes) {
		m.cursor = len(m.visibleNodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	
	// Update viewport
	m.updateViewport()
}

// updateViewport adjusts the viewport to keep cursor visible
func (m *BrowserModel) updateViewport() {
	m.viewport.size = m.height - 4 // Reserve space for header and footer
	
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

// Update handles browser messages and key events
func (m *BrowserModel) Update(msg tea.Msg) (*BrowserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.refreshView()
	case BrowserMsg:
		return m.handleBrowserMsg(msg)
	}
	
	return m, nil
}

// handleKeyPress processes keyboard input
func (m *BrowserModel) handleKeyPress(msg tea.KeyMsg) (*BrowserModel, tea.Cmd) {
	if m.confirmMode {
		return m.handleConfirmMode(msg)
	}
	
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.updateViewport()
		}
	case "down", "j":
		if m.cursor < len(m.visibleNodes)-1 {
			m.cursor++
			m.updateViewport()
		}
	case "left", "h":
		return m.handleLeft()
	case "right", "l", "enter":
		return m.handleRight()
	case "space":
		return m.handleSelection()
	case "s":
		m.showStats = !m.showStats
	case "r":
		return m.handleRefresh()
	case "c":
		if m.getCurrentNode() != nil && m.getCurrentNode().IsDir {
			m.confirmMode = true
		}
	case "home":
		m.cursor = 0
		m.updateViewport()
	case "end":
		m.cursor = len(m.visibleNodes) - 1
		m.updateViewport()
	case "pageup":
		m.cursor -= m.viewport.size
		if m.cursor < 0 {
			m.cursor = 0
		}
		m.updateViewport()
	case "pagedown":
		m.cursor += m.viewport.size
		if m.cursor >= len(m.visibleNodes) {
			m.cursor = len(m.visibleNodes) - 1
		}
		m.updateViewport()
	}
	
	return m, nil
}

// handleConfirmMode processes input in confirmation mode
func (m *BrowserModel) handleConfirmMode(msg tea.KeyMsg) (*BrowserModel, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		// Confirm selection
		m.confirmMode = false
		return m, m.selectFolder()
	case "n", "esc":
		// Cancel
		m.confirmMode = false
	}
	
	return m, nil
}

// handleLeft processes left arrow (collapse/go up)
func (m *BrowserModel) handleLeft() (*BrowserModel, tea.Cmd) {
	currentNode := m.getCurrentNode()
	if currentNode == nil {
		return m, nil
	}
	
	if currentNode.IsDir && currentNode.IsExpanded {
		// Collapse current directory
		m.tree.CollapseNode(currentNode)
		m.refreshView()
	} else if currentNode.Parent != nil {
		// Go to parent directory
		parentIndex := m.findNodeIndex(currentNode.Parent)
		if parentIndex >= 0 {
			m.cursor = parentIndex
			m.updateViewport()
		}
	}
	
	return m, nil
}

// handleRight processes right arrow (expand/enter)
func (m *BrowserModel) handleRight() (*BrowserModel, tea.Cmd) {
	currentNode := m.getCurrentNode()
	if currentNode == nil {
		return m, nil
	}
	
	if currentNode.IsDir {
		if !currentNode.IsExpanded {
			// Expand directory
			err := m.tree.ExpandNode(currentNode)
			if err != nil {
				m.errorMessage = fmt.Sprintf("Error expanding folder: %v", err)
			} else {
				m.refreshView()
			}
		} else if len(currentNode.Children) > 0 {
			// Go to first child
			m.cursor++
			m.updateViewport()
		}
	}
	
	return m, nil
}

// handleSelection toggles selection of current item
func (m *BrowserModel) handleSelection() (*BrowserModel, tea.Cmd) {
	currentNode := m.getCurrentNode()
	if currentNode != nil {
		m.tree.SelectNode(currentNode)
	}
	
	return m, nil
}

// handleRefresh refreshes the current view
func (m *BrowserModel) handleRefresh() (*BrowserModel, tea.Cmd) {
	err := m.tree.refreshTree()
	if err != nil {
		m.errorMessage = fmt.Sprintf("Error refreshing: %v", err)
	} else {
		m.refreshView()
		m.errorMessage = ""
	}
	
	return m, nil
}

// selectFolder returns a command to select the current folder
func (m *BrowserModel) selectFolder() tea.Cmd {
	return func() tea.Msg {
		return BrowserMsg{
			Type: "folder_selected",
			Data: m.getCurrentNode(),
		}
	}
}

// getCurrentNode returns the currently highlighted node
func (m *BrowserModel) getCurrentNode() *FolderNode {
	if m.cursor >= 0 && m.cursor < len(m.visibleNodes) {
		return m.visibleNodes[m.cursor]
	}
	return nil
}

// findNodeIndex finds the index of a node in the visible list
func (m *BrowserModel) findNodeIndex(node *FolderNode) int {
	for i, n := range m.visibleNodes {
		if n == node {
			return i
		}
	}
	return -1
}

// View renders the folder browser
func (m *BrowserModel) View() string {
	var result strings.Builder
	
	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.width)
	
	currentPath := m.tree.GetPath()
	if len(currentPath) > m.width-20 {
		currentPath = "..." + currentPath[len(currentPath)-(m.width-23):]
	}
	
	header := fmt.Sprintf("📁 Browse Folders: %s", currentPath)
	result.WriteString(headerStyle.Render(header))
	result.WriteString("\n\n")
	
	// Error message
	if m.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)
		result.WriteString(errorStyle.Render("⚠️ " + m.errorMessage))
		result.WriteString("\n\n")
	}
	
	// Folder tree
	if len(m.visibleNodes) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
		result.WriteString(emptyStyle.Render("No folders found"))
	} else {
		// Render visible portion of tree
		start := m.viewport.offset
		end := start + m.viewport.size
		if end > len(m.visibleNodes) {
			end = len(m.visibleNodes)
		}
		
		for i := start; i < end; i++ {
			node := m.visibleNodes[i]
			isSelected := i == m.cursor
			line := RenderTreeLine(node, isSelected, m.width-2)
			result.WriteString(line)
			result.WriteString("\n")
		}
	}
	
	// Footer with stats and instructions
	result.WriteString("\n")
	result.WriteString(m.renderFooter())
	
	// Confirmation dialog
	if m.confirmMode {
		result.WriteString("\n")
		result.WriteString(m.renderConfirmDialog())
	}
	
	return result.String()
}

// renderFooter renders the footer with stats and controls
func (m *BrowserModel) renderFooter() string {
	var result strings.Builder
	
	// Current selection stats
	if m.showStats {
		currentNode := m.getCurrentNode()
		if currentNode != nil && currentNode.IsDir {
			statsStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3B82F6")).
				BorderTop(true).
				BorderStyle(lipgloss.NormalBorder()).
				Padding(1, 0)
			
			stats := fmt.Sprintf("📊 Selected: %s | Files: %s | Size: %s",
				currentNode.Name,
				FormatCount(currentNode.FileCount),
				FormatSize(currentNode.Size))
			
			result.WriteString(statsStyle.Render(stats))
			result.WriteString("\n")
		}
	}
	
	// Instructions
	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true)
	
	instructions := "↑↓: navigate • ←→: collapse/expand • Space: select • C: confirm • S: toggle stats • R: refresh"
	result.WriteString(instructionStyle.Render(instructions))
	
	return result.String()
}

// renderConfirmDialog renders the confirmation dialog
func (m *BrowserModel) renderConfirmDialog() string {
	currentNode := m.getCurrentNode()
	if currentNode == nil {
		return ""
	}
	
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#F59E0B")).
		Background(lipgloss.Color("#FEF3C7")).
		Foreground(lipgloss.Color("#92400E")).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)
	
	message := fmt.Sprintf("Select folder '%s'?\n\nThis will scan %s files (%s) and generate context.\n\nPress Y to confirm, N to cancel.",
		currentNode.Name,
		FormatCount(currentNode.FileCount),
		FormatSize(currentNode.Size))
	
	return dialogStyle.Render(message)
}

// GetSelectedFolder returns the currently selected folder
func (m *BrowserModel) GetSelectedFolder() *FolderNode {
	return m.tree.GetSelectedNode()
}

// SetSize updates the browser dimensions
func (m *BrowserModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.refreshView()
}

// GetTree returns the underlying folder tree
func (m *BrowserModel) GetTree() *FolderTree {
	return m.tree
}

// handleBrowserMsg processes browser-specific messages
func (m *BrowserModel) handleBrowserMsg(msg BrowserMsg) (*BrowserModel, tea.Cmd) {
	switch msg.Type {
	case "folder_selected":
		// Handle folder selection
		return m, nil
	case "refresh":
		return m.handleRefresh()
	}
	
	return m, nil
}