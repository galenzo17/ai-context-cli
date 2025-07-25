package app

import (
	"ai-context-cli/internal/context"
	tea "github.com/charmbracelet/bubbletea"
)

// startFolderScan starts scanning a specific folder
func (m Model) startFolderScan(folderPath string) tea.Cmd {
	return func() tea.Msg {
		// Create scanner with folder-specific config
		config := context.DefaultScanConfig(folderPath)
		scanner := context.NewProjectScanner(config)
		
		// Perform the scan
		result, err := scanner.Scan()
		if err != nil {
			return ScanCompleteMsg{Error: err}
		}
		
		return ScanCompleteMsg{Result: result}
	}
}