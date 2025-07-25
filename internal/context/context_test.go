package context

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultScanConfig(t *testing.T) {
	config := DefaultScanConfig("/test/path")
	
	if config.RootPath != "/test/path" {
		t.Errorf("Expected RootPath '/test/path', got '%s'", config.RootPath)
	}
	
	if len(config.ExcludePatterns) == 0 {
		t.Error("Expected default exclude patterns")
	}
	
	if len(config.ExcludeExtensions) == 0 {
		t.Error("Expected default exclude extensions")
	}
	
	if config.MaxDepth != 50 {
		t.Errorf("Expected MaxDepth 50, got %d", config.MaxDepth)
	}
	
	if config.MaxFileSize != 10*1024*1024 {
		t.Errorf("Expected MaxFileSize 10MB, got %d", config.MaxFileSize)
	}
}

func TestProjectScannerCreation(t *testing.T) {
	config := DefaultScanConfig("/test")
	scanner := NewProjectScanner(config)
	
	if scanner == nil {
		t.Error("Expected non-nil scanner")
	}
	
	if scanner.config.RootPath != "/test" {
		t.Errorf("Expected config path '/test', got '%s'", scanner.config.RootPath)
	}
	
	if scanner.progress == nil {
		t.Error("Expected progress channel to be initialized")
	}
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}
	
	for _, tc := range testCases {
		result := FormatSize(tc.bytes)
		if result != tc.expected {
			t.Errorf("FormatSize(%d) = '%s', expected '%s'", tc.bytes, result, tc.expected)
		}
	}
}

func TestEstimateProcessingTime(t *testing.T) {
	// Test with different file counts
	testCases := []struct {
		fileCount int
		minTime   time.Duration
		maxTime   time.Duration
	}{
		{0, 0, time.Second},
		{100, time.Millisecond, time.Second},
		{1000, 10*time.Millisecond, 5*time.Second},
	}
	
	for _, tc := range testCases {
		duration := EstimateProcessingTime(tc.fileCount)
		if duration < tc.minTime || duration > tc.maxTime {
			t.Errorf("EstimateProcessingTime(%d) = %v, expected between %v and %v", 
				tc.fileCount, duration, tc.minTime, tc.maxTime)
		}
	}
}

func TestContextGeneratorCreation(t *testing.T) {
	generator := NewContextGenerator()
	
	if generator == nil {
		t.Error("Expected non-nil generator")
	}
	
	if generator.maxFileSize <= 0 {
		t.Error("Expected positive maxFileSize")
	}
	
	if generator.maxTotalSize <= 0 {
		t.Error("Expected positive maxTotalSize")
	}
	
	if !generator.includeContent {
		t.Error("Expected includeContent to be true by default")
	}
	
	if !generator.includeSummary {
		t.Error("Expected includeSummary to be true by default")
	}
}

func TestContextGeneratorSetOptions(t *testing.T) {
	generator := NewContextGenerator()
	
	generator.SetOptions(1024, 2048, false, false)
	
	if generator.maxFileSize != 1024 {
		t.Errorf("Expected maxFileSize 1024, got %d", generator.maxFileSize)
	}
	
	if generator.maxTotalSize != 2048 {
		t.Errorf("Expected maxTotalSize 2048, got %d", generator.maxTotalSize)
	}
	
	if generator.includeContent {
		t.Error("Expected includeContent to be false")
	}
	
	if generator.includeSummary {
		t.Error("Expected includeSummary to be false")
	}
}

func TestFileInfoCreation(t *testing.T) {
	fileInfo := FileInfo{
		Path:      "/test/file.go",
		Size:      1024,
		Lines:     50,
		Extension: ".go",
		ModTime:   time.Now(),
	}
	
	if fileInfo.Path != "/test/file.go" {
		t.Errorf("Expected path '/test/file.go', got '%s'", fileInfo.Path)
	}
	
	if fileInfo.Extension != ".go" {
		t.Errorf("Expected extension '.go', got '%s'", fileInfo.Extension)
	}
	
	if fileInfo.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", fileInfo.Size)
	}
}

func TestScanResultCreation(t *testing.T) {
	result := &ScanResult{
		TotalFiles:       10,
		TotalDirectories: 3,
		TotalSize:        5120,
		TotalLines:       500,
		Extensions:      make(map[string]int),
		Files:           make([]FileInfo, 0),
	}
	
	if result.TotalFiles != 10 {
		t.Errorf("Expected TotalFiles 10, got %d", result.TotalFiles)
	}
	
	if result.Extensions == nil {
		t.Error("Expected Extensions map to be initialized")
	}
	
	if result.Files == nil {
		t.Error("Expected Files slice to be initialized")
	}
}

func TestContextSectionCreation(t *testing.T) {
	section := ContextSection{
		Title:   "Test Section",
		Content: "Test content",
		Files:   []string{"file1.go", "file2.go"},
	}
	
	if section.Title != "Test Section" {
		t.Errorf("Expected title 'Test Section', got '%s'", section.Title)
	}
	
	if len(section.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(section.Files))
	}
}

func TestContextResultCreation(t *testing.T) {
	result := &ContextResult{
		ProjectName:   "Test Project",
		GeneratedAt:   time.Now(),
		TotalFiles:    10,
		Sections:      make([]ContextSection, 0),
		TokenEstimate: 1000,
	}
	
	if result.ProjectName != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", result.ProjectName)
	}
	
	if result.TokenEstimate != 1000 {
		t.Errorf("Expected token estimate 1000, got %d", result.TokenEstimate)
	}
	
	if result.Sections == nil {
		t.Error("Expected Sections slice to be initialized")
	}
}

// Integration tests that require file system operations
func TestScannerWithRealFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "context_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create some test files
	testFiles := []struct {
		path    string
		content string
	}{
		{"test.go", "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}"},
		{"README.md", "# Test Project\n\nThis is a test."},
		{"config.json", "{\"name\": \"test\"}"},
	}
	
	for _, tf := range testFiles {
		fullPath := filepath.Join(tempDir, tf.path)
		err := os.WriteFile(fullPath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}
	
	// Test scanning
	config := DefaultScanConfig(tempDir)
	scanner := NewProjectScanner(config)
	
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	
	if result.TotalFiles != 3 {
		t.Errorf("Expected 3 files, got %d", result.TotalFiles)
	}
	
	if len(result.Extensions) == 0 {
		t.Error("Expected some extensions to be found")
	}
	
	// Test context generation
	generator := NewContextGenerator()
	contextResult, err := generator.GenerateContext(result, "Test Project")
	if err != nil {
		t.Fatalf("Context generation failed: %v", err)
	}
	
	if contextResult.ProjectName != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", contextResult.ProjectName)
	}
	
	if len(contextResult.Sections) == 0 {
		t.Error("Expected some sections to be generated")
	}
	
	if contextResult.TokenEstimate <= 0 {
		t.Error("Expected positive token estimate")
	}
}

func TestScannerExcludePatterns(t *testing.T) {
	config := DefaultScanConfig("/test")
	scanner := NewProjectScanner(config)
	
	testCases := []struct {
		path      string
		isDir     bool
		shouldExclude bool
	}{
		{"node_modules/package.json", false, true},
		{".git/config", false, true},
		{"src/main.go", false, false},
		{".env", false, true},  // hidden file
		{"dist/bundle.js", false, true},
		{"README.md", false, false},
	}
	
	for _, tc := range testCases {
		result := scanner.shouldExcludePath(tc.path, tc.isDir)
		if result != tc.shouldExclude {
			t.Errorf("shouldExcludePath('%s', %v) = %v, expected %v", 
				tc.path, tc.isDir, result, tc.shouldExclude)
		}
	}
}

func TestGeneratorLanguageDetection(t *testing.T) {
	generator := NewContextGenerator()
	
	testCases := []struct {
		extension string
		expected  string
	}{
		{".go", "go"},
		{".js", "javascript"},
		{".ts", "typescript"},
		{".py", "python"},
		{".json", "json"},
		{".md", "markdown"},
		{".unknown", ""},
	}
	
	for _, tc := range testCases {
		result := generator.getLanguageFromExtension(tc.extension)
		if result != tc.expected {
			t.Errorf("getLanguageFromExtension('%s') = '%s', expected '%s'", 
				tc.extension, result, tc.expected)
		}
	}
}

func TestGeneratorTextFileDetection(t *testing.T) {
	generator := NewContextGenerator()
	
	testCases := []struct {
		extension string
		isText    bool
	}{
		{".go", true},
		{".txt", true},
		{".md", true},
		{".json", true},
		{".jpg", false},
		{".exe", false},
		{".pdf", false},
		{"", false}, // empty extension
	}
	
	for _, tc := range testCases {
		result := generator.isTextFile(tc.extension)
		if result != tc.isText {
			t.Errorf("isTextFile('%s') = %v, expected %v", 
				tc.extension, result, tc.isText)
		}
	}
}