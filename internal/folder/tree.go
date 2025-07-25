package folder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// FolderNode represents a node in the folder tree
type FolderNode struct {
	Name       string
	Path       string
	IsDir      bool
	Size       int64
	FileCount  int
	DirCount   int
	ModTime    time.Time
	Children   []*FolderNode
	Parent     *FolderNode
	IsExpanded bool
	IsSelected bool
	Level      int
}

// FolderStats represents statistics for a folder
type FolderStats struct {
	TotalFiles       int
	TotalDirectories int
	TotalSize        int64
	LastModified     time.Time
	FileTypes        map[string]int
}

// FolderTree manages the folder tree structure and navigation
type FolderTree struct {
	root           *FolderNode
	currentPath    string
	selectedNode   *FolderNode
	expandedPaths  map[string]bool
	maxDepth       int
	showHidden     bool
	sortBy         SortType
}

// SortType defines how folders should be sorted
type SortType int

const (
	SortByName SortType = iota
	SortBySize
	SortByDate
	SortByType
)

// NewFolderTree creates a new folder tree
func NewFolderTree(rootPath string) (*FolderTree, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}
	
	// Check if path exists and is a directory
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access path: %w", err)
	}
	
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}
	
	tree := &FolderTree{
		currentPath:   absPath,
		expandedPaths: make(map[string]bool),
		maxDepth:      10,
		showHidden:    false,
		sortBy:        SortByName,
	}
	
	err = tree.buildTree()
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}
	
	return tree, nil
}

// buildTree constructs the folder tree structure
func (ft *FolderTree) buildTree() error {
	info, err := os.Stat(ft.currentPath)
	if err != nil {
		return err
	}
	
	ft.root = &FolderNode{
		Name:       filepath.Base(ft.currentPath),
		Path:       ft.currentPath,
		IsDir:      true,
		ModTime:    info.ModTime(),
		IsExpanded: true,
		Level:      0,
	}
	
	return ft.loadChildren(ft.root)
}

// loadChildren loads child nodes for a given directory
func (ft *FolderTree) loadChildren(node *FolderNode) error {
	if !node.IsDir || node.Level >= ft.maxDepth {
		return nil
	}
	
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return fmt.Errorf("cannot read directory %s: %w", node.Path, err)
	}
	
	node.Children = make([]*FolderNode, 0)
	
	for _, entry := range entries {
		// Skip hidden files/directories if not showing hidden
		if !ft.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		fullPath := filepath.Join(node.Path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue // Skip files we can't stat
		}
		
		child := &FolderNode{
			Name:       entry.Name(),
			Path:       fullPath,
			IsDir:      entry.IsDir(),
			Size:       info.Size(),
			ModTime:    info.ModTime(),
			Parent:     node,
			IsExpanded: ft.expandedPaths[fullPath],
			Level:      node.Level + 1,
		}
		
		// Calculate stats for directories
		if child.IsDir {
			ft.calculateStats(child)
			
			// Load children if expanded
			if child.IsExpanded {
				ft.loadChildren(child)
			}
		}
		
		node.Children = append(node.Children, child)
	}
	
	// Sort children
	ft.sortChildren(node.Children)
	
	return nil
}

// calculateStats calculates statistics for a directory
func (ft *FolderTree) calculateStats(node *FolderNode) {
	if !node.IsDir {
		return
	}
	
	stats, err := ft.GetFolderStats(node.Path)
	if err != nil {
		return
	}
	
	node.FileCount = stats.TotalFiles
	node.DirCount = stats.TotalDirectories
	node.Size = stats.TotalSize
}

// GetFolderStats calculates comprehensive statistics for a folder
func (ft *FolderTree) GetFolderStats(folderPath string) (*FolderStats, error) {
	stats := &FolderStats{
		FileTypes: make(map[string]int),
	}
	
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue on errors
		}
		
		// Skip hidden files if not showing hidden
		if !ft.showHidden && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		
		info, err := d.Info()
		if err != nil {
			return nil
		}
		
		if d.IsDir() {
			stats.TotalDirectories++
		} else {
			stats.TotalFiles++
			stats.TotalSize += info.Size()
			
			// Track file types
			ext := strings.ToLower(filepath.Ext(path))
			if ext == "" {
				ext = "(no extension)"
			}
			stats.FileTypes[ext]++
		}
		
		// Update last modified time
		if info.ModTime().After(stats.LastModified) {
			stats.LastModified = info.ModTime()
		}
		
		return nil
	})
	
	return stats, err
}

// sortChildren sorts child nodes based on the current sort type
func (ft *FolderTree) sortChildren(children []*FolderNode) {
	sort.Slice(children, func(i, j int) bool {
		a, b := children[i], children[j]
		
		// Directories first
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		
		switch ft.sortBy {
		case SortBySize:
			if a.Size != b.Size {
				return a.Size > b.Size
			}
		case SortByDate:
			if !a.ModTime.Equal(b.ModTime) {
				return a.ModTime.After(b.ModTime)
			}
		case SortByType:
			extA := filepath.Ext(a.Name)
			extB := filepath.Ext(b.Name)
			if extA != extB {
				return extA < extB
			}
		}
		
		// Default to name sorting
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}

// ExpandNode expands a directory node
func (ft *FolderTree) ExpandNode(node *FolderNode) error {
	if !node.IsDir || node.IsExpanded {
		return nil
	}
	
	node.IsExpanded = true
	ft.expandedPaths[node.Path] = true
	
	return ft.loadChildren(node)
}

// CollapseNode collapses a directory node
func (ft *FolderTree) CollapseNode(node *FolderNode) {
	if !node.IsDir || !node.IsExpanded {
		return
	}
	
	node.IsExpanded = false
	delete(ft.expandedPaths, node.Path)
	
	// Recursively collapse children
	for _, child := range node.Children {
		if child.IsDir {
			ft.CollapseNode(child)
		}
	}
}

// ToggleNode toggles the expansion state of a node
func (ft *FolderTree) ToggleNode(node *FolderNode) error {
	if !node.IsDir {
		return nil
	}
	
	if node.IsExpanded {
		ft.CollapseNode(node)
	} else {
		return ft.ExpandNode(node)
	}
	
	return nil
}

// SelectNode selects a node
func (ft *FolderTree) SelectNode(node *FolderNode) {
	// Deselect previous node
	if ft.selectedNode != nil {
		ft.selectedNode.IsSelected = false
	}
	
	// Select new node
	node.IsSelected = true
	ft.selectedNode = node
}

// GetSelectedNode returns the currently selected node
func (ft *FolderTree) GetSelectedNode() *FolderNode {
	return ft.selectedNode
}

// GetVisibleNodes returns all currently visible nodes in display order
func (ft *FolderTree) GetVisibleNodes() []*FolderNode {
	var nodes []*FolderNode
	ft.collectVisibleNodes(ft.root, &nodes)
	return nodes
}

// collectVisibleNodes recursively collects visible nodes
func (ft *FolderTree) collectVisibleNodes(node *FolderNode, nodes *[]*FolderNode) {
	*nodes = append(*nodes, node)
	
	if node.IsExpanded {
		for _, child := range node.Children {
			ft.collectVisibleNodes(child, nodes)
		}
	}
}

// SetSortType changes the sorting method
func (ft *FolderTree) SetSortType(sortType SortType) error {
	ft.sortBy = sortType
	return ft.refreshTree()
}

// SetShowHidden toggles hidden file/directory visibility
func (ft *FolderTree) SetShowHidden(show bool) error {
	ft.showHidden = show
	return ft.refreshTree()
}

// refreshTree rebuilds the tree with current settings
func (ft *FolderTree) refreshTree() error {
	return ft.buildTree()
}

// GetPath returns the current root path
func (ft *FolderTree) GetPath() string {
	return ft.currentPath
}

// NavigateToPath changes the root path
func (ft *FolderTree) NavigateToPath(newPath string) error {
	absPath, err := filepath.Abs(newPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("cannot access path: %w", err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}
	
	ft.currentPath = absPath
	ft.selectedNode = nil
	ft.expandedPaths = make(map[string]bool)
	
	return ft.buildTree()
}

// GetNodeByPath finds a node by its path
func (ft *FolderTree) GetNodeByPath(path string) *FolderNode {
	return ft.findNodeByPath(ft.root, path)
}

// findNodeByPath recursively searches for a node by path
func (ft *FolderTree) findNodeByPath(node *FolderNode, path string) *FolderNode {
	if node.Path == path {
		return node
	}
	
	for _, child := range node.Children {
		if found := ft.findNodeByPath(child, path); found != nil {
			return found
		}
	}
	
	return nil
}

// FormatSize formats a file size in human-readable format
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatCount formats a count with K/M suffixes
func FormatCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	} else if count < 1000000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	}
}

// RenderTreeLine renders a single line of the tree
func RenderTreeLine(node *FolderNode, isSelected bool, width int) string {
	var result strings.Builder
	
	// Build indentation
	indent := strings.Repeat("  ", node.Level)
	result.WriteString(indent)
	
	// Add expansion indicator for directories
	if node.IsDir {
		if node.IsExpanded {
			result.WriteString("▼ ")
		} else {
			result.WriteString("▶ ")
		}
	} else {
		result.WriteString("  ")
	}
	
	// Add icon
	if node.IsDir {
		result.WriteString("📁 ")
	} else {
		result.WriteString("📄 ")
	}
	
	// Add name
	name := node.Name
	if len(name) > 30 {
		name = name[:27] + "..."
	}
	result.WriteString(name)
	
	// Add stats for directories
	if node.IsDir && (node.FileCount > 0 || node.DirCount > 0) {
		stats := fmt.Sprintf(" (%s, %s files)", 
			FormatSize(node.Size), 
			FormatCount(node.FileCount))
		result.WriteString(stats)
	} else if !node.IsDir {
		size := fmt.Sprintf(" (%s)", FormatSize(node.Size))
		result.WriteString(size)
	}
	
	// Apply styling
	line := result.String()
	
	if isSelected {
		style := lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Width(width)
		return style.Render(line)
	} else {
		style := lipgloss.NewStyle().
			Width(width)
		if node.IsDir {
			style = style.Foreground(lipgloss.Color("#3B82F6"))
		} else {
			style = style.Foreground(lipgloss.Color("#6B7280"))
		}
		return style.Render(line)
	}
}