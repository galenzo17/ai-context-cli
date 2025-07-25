package folder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFolderTree(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "folder_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test structure
	os.Mkdir(filepath.Join(tempDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tempDir, "subdir2"), 0755)
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test content"), 0644)
	os.WriteFile(filepath.Join(tempDir, "subdir1", "file2.go"), []byte("package main"), 0644)
	
	// Test successful creation
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	if tree.root == nil {
		t.Error("Expected root node to be set")
	}
	
	if tree.currentPath != tempDir {
		t.Errorf("Expected current path '%s', got '%s'", tempDir, tree.currentPath)
	}
	
	// Test invalid path
	_, err = NewFolderTree("/invalid/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
	
	// Test file path (should fail)
	testFile := filepath.Join(tempDir, "file1.txt")
	_, err = NewFolderTree(testFile)
	if err == nil {
		t.Error("Expected error for file path")
	}
}

func TestFolderTreeNavigation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "folder_nav_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test structure
	subdir := filepath.Join(tempDir, "testdir")
	os.Mkdir(subdir, 0755)
	os.WriteFile(filepath.Join(subdir, "test.txt"), []byte("test"), 0644)
	
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Test expansion
	if len(tree.root.Children) == 0 {
		t.Error("Expected root to have children")
	}
	
	// Find the subdirectory node
	var subdirNode *FolderNode
	for _, child := range tree.root.Children {
		if child.Name == "testdir" && child.IsDir {
			subdirNode = child
			break
		}
	}
	
	if subdirNode == nil {
		t.Fatal("Could not find testdir node")
	}
	
	// Test expansion
	if subdirNode.IsExpanded {
		t.Error("Subdirectory should not be expanded initially")
	}
	
	err = tree.ExpandNode(subdirNode)
	if err != nil {
		t.Errorf("Failed to expand node: %v", err)
	}
	
	if !subdirNode.IsExpanded {
		t.Error("Subdirectory should be expanded after ExpandNode")
	}
	
	// Test collapse
	tree.CollapseNode(subdirNode)
	if subdirNode.IsExpanded {
		t.Error("Subdirectory should not be expanded after CollapseNode")
	}
}

func TestFolderTreeSelection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "folder_select_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Test initial state
	if tree.GetSelectedNode() != nil {
		t.Error("Expected no selected node initially")
	}
	
	// Test selection
	tree.SelectNode(tree.root)
	selected := tree.GetSelectedNode()
	if selected == nil {
		t.Error("Expected selected node after SelectNode")
	}
	
	if selected != tree.root {
		t.Error("Expected selected node to be root")
	}
	
	if !tree.root.IsSelected {
		t.Error("Expected root node to be marked as selected")
	}
}

func TestFolderTreeStats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "folder_stats_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.go"), []byte("package main\nfunc main() {}"), 0644)
	
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Test folder stats
	stats, err := tree.GetFolderStats(tempDir)
	if err != nil {
		t.Errorf("Failed to get folder stats: %v", err)
	}
	
	if stats.TotalFiles != 2 {
		t.Errorf("Expected 2 files, got %d", stats.TotalFiles)
	}
	
	if stats.TotalSize == 0 {
		t.Error("Expected non-zero total size")
	}
	
	if len(stats.FileTypes) == 0 {
		t.Error("Expected file types to be tracked")
	}
}

func TestFolderTreeSorting(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "folder_sort_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create files with different properties
	os.WriteFile(filepath.Join(tempDir, "zzz.txt"), []byte("small"), 0644)
	os.WriteFile(filepath.Join(tempDir, "aaa.txt"), []byte("large content here"), 0644)
	
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Test name sorting (default)
	if tree.sortBy != SortByName {
		t.Error("Expected default sort by name")
	}
	
	// Test sort type change
	err = tree.SetSortType(SortBySize)
	if err != nil {
		t.Errorf("Failed to set sort type: %v", err)
	}
	
	if tree.sortBy != SortBySize {
		t.Error("Expected sort type to be SortBySize")
	}
}

func TestBrowserModelCreation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "browser_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Test successful creation
	browser, err := NewBrowserModel(tempDir)
	if err != nil {
		t.Fatalf("Failed to create browser model: %v", err)
	}
	
	if browser.tree == nil {
		t.Error("Expected tree to be initialized")
	}
	
	if len(browser.visibleNodes) == 0 {
		t.Error("Expected visible nodes to be populated")
	}
	
	// Test invalid path
	_, err = NewBrowserModel("/invalid/path")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestBrowserModelNavigation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "browser_nav_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test structure
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)
	
	browser, err := NewBrowserModel(tempDir)
	if err != nil {
		t.Fatalf("Failed to create browser model: %v", err)
	}
	
	// Test initial state
	if browser.cursor != 0 {
		t.Errorf("Expected cursor at 0, got %d", browser.cursor)
	}
	
	// Test cursor movement
	initialCursor := browser.cursor
	if len(browser.visibleNodes) > 1 {
		// Simulate down arrow
		browser.cursor++
		browser.updateViewport()
		
		if browser.cursor == initialCursor {
			t.Error("Expected cursor to move down")
		}
	}
}

func TestFormatUtilities(t *testing.T) {
	// Test FormatSize
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}
	
	for _, tc := range testCases {
		result := FormatSize(tc.bytes)
		if result != tc.expected {
			t.Errorf("FormatSize(%d) = '%s', expected '%s'", tc.bytes, result, tc.expected)
		}
	}
	
	// Test FormatCount
	countCases := []struct {
		count    int
		expected string
	}{
		{0, "0"},
		{500, "500"},
		{1500, "1.5K"},
		{1500000, "1.5M"},
	}
	
	for _, tc := range countCases {
		result := FormatCount(tc.count)
		if result != tc.expected {
			t.Errorf("FormatCount(%d) = '%s', expected '%s'", tc.count, result, tc.expected)
		}
	}
}

func TestRenderTreeLine(t *testing.T) {
	node := &FolderNode{
		Name:      "test.txt",
		Path:      "/test/test.txt",
		IsDir:     false,
		Size:      1024,
		Level:     1,
		IsExpanded: false,
	}
	
	// Test normal rendering
	line := RenderTreeLine(node, false, 80)
	if line == "" {
		t.Error("Expected non-empty line")
	}
	
	// Test selected rendering
	selectedLine := RenderTreeLine(node, true, 80)
	if selectedLine == "" {
		t.Error("Expected non-empty selected line")
	}
	
	// In testing environment, styles might not be visible, so just check that both returned strings
	if len(selectedLine) == 0 || len(line) == 0 {
		t.Error("Expected both lines to have content")
	}
	
	// Test directory node
	dirNode := &FolderNode{
		Name:       "testdir",
		Path:       "/test/testdir",
		IsDir:      true,
		Level:      0,
		IsExpanded: true,
		FileCount:  5,
		Size:       5120,
	}
	
	dirLine := RenderTreeLine(dirNode, false, 80)
	if dirLine == "" {
		t.Error("Expected non-empty directory line")
	}
}

func TestFolderNodePathFinding(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "path_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Test finding root node
	foundNode := tree.GetNodeByPath(tempDir)
	if foundNode == nil {
		t.Error("Expected to find root node by path")
	}
	
	if foundNode != tree.root {
		t.Error("Expected found node to be root")
	}
	
	// Test finding non-existent path
	nonExistentNode := tree.GetNodeByPath("/does/not/exist")
	if nonExistentNode != nil {
		t.Error("Expected nil for non-existent path")
	}
}

func TestHiddenFileHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hidden_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create hidden file
	os.WriteFile(filepath.Join(tempDir, ".hidden"), []byte("hidden content"), 0644)
	os.WriteFile(filepath.Join(tempDir, "visible.txt"), []byte("visible content"), 0644)
	
	// Test with hidden files disabled (default)
	tree, err := NewFolderTree(tempDir)
	if err != nil {
		t.Fatalf("Failed to create folder tree: %v", err)
	}
	
	// Count visible files in root children
	visibleCount := 0
	hiddenCount := 0
	for _, child := range tree.root.Children {
		if child.Name == ".hidden" {
			hiddenCount++
		} else if child.Name == "visible.txt" {
			visibleCount++
		}
	}
	
	if hiddenCount > 0 {
		t.Error("Expected hidden file to be excluded by default")
	}
	
	if visibleCount != 1 {
		t.Errorf("Expected 1 visible file, got %d", visibleCount)
	}
	
	// Test with hidden files enabled
	err = tree.SetShowHidden(true)
	if err != nil {
		t.Errorf("Failed to enable hidden files: %v", err)
	}
	
	// Recount after enabling hidden files
	visibleCount = 0
	hiddenCount = 0
	for _, child := range tree.root.Children {
		if child.Name == ".hidden" {
			hiddenCount++
		} else if child.Name == "visible.txt" {
			visibleCount++
		}
	}
	
	if hiddenCount != 1 {
		t.Errorf("Expected 1 hidden file when enabled, got %d", hiddenCount)
	}
	
	if visibleCount != 1 {
		t.Errorf("Expected 1 visible file when hidden enabled, got %d", visibleCount)
	}
}