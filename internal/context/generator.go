package context

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ContextSection represents a section of the generated context
type ContextSection struct {
	Title   string
	Content string
	Files   []string
}

// ContextResult represents the generated context
type ContextResult struct {
	ProjectName    string
	GeneratedAt    time.Time
	TotalFiles     int
	TotalSize      int64
	Sections       []ContextSection
	Summary        string
	TokenEstimate  int
}

// ContextGenerator generates comprehensive context from scan results
type ContextGenerator struct {
	maxFileSize     int64
	maxTotalSize    int64
	includeContent  bool
	includeSummary  bool
	priorityExtensions []string
}

// NewContextGenerator creates a new context generator
func NewContextGenerator() *ContextGenerator {
	return &ContextGenerator{
		maxFileSize:    50 * 1024,    // 50KB per file
		maxTotalSize:   10 * 1024 * 1024, // 10MB total
		includeContent: true,
		includeSummary: true,
		priorityExtensions: []string{
			".go", ".js", ".ts", ".py", ".java", ".c", ".cpp",
			".md", ".txt", ".json", ".yaml", ".yml",
		},
	}
}

// SetOptions configures the context generator
func (cg *ContextGenerator) SetOptions(maxFileSize, maxTotalSize int64, includeContent, includeSummary bool) {
	cg.maxFileSize = maxFileSize
	cg.maxTotalSize = maxTotalSize
	cg.includeContent = includeContent
	cg.includeSummary = includeSummary
}

// GenerateContext creates comprehensive context from scan results
func (cg *ContextGenerator) GenerateContext(scanResult *ScanResult, projectName string) (*ContextResult, error) {
	result := &ContextResult{
		ProjectName: projectName,
		GeneratedAt: time.Now(),
		TotalFiles:  scanResult.TotalFiles,
		TotalSize:   scanResult.TotalSize,
		Sections:    make([]ContextSection, 0),
	}
	
	// Generate project overview section
	result.Sections = append(result.Sections, cg.generateOverviewSection(scanResult))
	
	// Generate directory structure section
	result.Sections = append(result.Sections, cg.generateStructureSection(scanResult))
	
	// Generate file type analysis section
	result.Sections = append(result.Sections, cg.generateFileTypeSection(scanResult))
	
	// Generate file content sections (if enabled)
	if cg.includeContent {
		contentSections, err := cg.generateContentSections(scanResult)
		if err != nil {
			return nil, fmt.Errorf("failed to generate content sections: %w", err)
		}
		result.Sections = append(result.Sections, contentSections...)
	}
	
	// Generate summary
	if cg.includeSummary {
		result.Summary = cg.generateSummary(scanResult, result)
	}
	
	// Estimate tokens
	result.TokenEstimate = cg.estimateTokens(result)
	
	return result, nil
}

// generateOverviewSection creates the project overview section
func (cg *ContextGenerator) generateOverviewSection(scanResult *ScanResult) ContextSection {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf("# Project Overview\n\n"))
	content.WriteString(fmt.Sprintf("**Scan completed:** %s\n", scanResult.ScanDuration.Round(time.Millisecond)))
	content.WriteString(fmt.Sprintf("**Total files:** %d\n", scanResult.TotalFiles))
	content.WriteString(fmt.Sprintf("**Total directories:** %d\n", scanResult.TotalDirectories))
	content.WriteString(fmt.Sprintf("**Total size:** %s\n", FormatSize(scanResult.TotalSize)))
	content.WriteString(fmt.Sprintf("**Total lines:** %s\n", FormatNumber(scanResult.TotalLines)))
	content.WriteString(fmt.Sprintf("**Excluded files:** %d\n\n", scanResult.ExcludedFiles))
	
	// Top file extensions
	content.WriteString("## File Extensions\n\n")
	sortedExts := cg.sortExtensionsByCount(scanResult.Extensions)
	for i, ext := range sortedExts {
		if i >= 10 { // Show top 10
			break
		}
		name := ext.Extension
		if name == "" {
			name = "(no extension)"
		}
		content.WriteString(fmt.Sprintf("- **%s**: %d files\n", name, ext.Count))
	}
	content.WriteString("\n")
	
	// Largest files
	if len(scanResult.LargestFiles) > 0 {
		content.WriteString("## Largest Files\n\n")
		for i, file := range scanResult.LargestFiles {
			if i >= 5 { // Show top 5
				break
			}
			relativePath := cg.getRelativePath(file.Path)
			content.WriteString(fmt.Sprintf("- **%s**: %s (%d lines)\n", 
				relativePath, FormatSize(file.Size), file.Lines))
		}
		content.WriteString("\n")
	}
	
	return ContextSection{
		Title:   "Project Overview",
		Content: content.String(),
		Files:   []string{},
	}
}

// generateStructureSection creates the directory structure section
func (cg *ContextGenerator) generateStructureSection(scanResult *ScanResult) ContextSection {
	var content strings.Builder
	
	content.WriteString("# Directory Structure\n\n")
	content.WriteString("```\n")
	
	// Build directory tree
	tree := cg.buildDirectoryTree(scanResult.Files)
	content.WriteString(tree)
	
	content.WriteString("```\n\n")
	
	return ContextSection{
		Title:   "Directory Structure",
		Content: content.String(),
		Files:   []string{},
	}
}

// generateFileTypeSection creates the file type analysis section
func (cg *ContextGenerator) generateFileTypeSection(scanResult *ScanResult) ContextSection {
	var content strings.Builder
	
	content.WriteString("# File Type Analysis\n\n")
	
	// Group files by extension
	filesByExt := make(map[string][]FileInfo)
	for _, file := range scanResult.Files {
		ext := file.Extension
		if ext == "" {
			ext = "(no extension)"
		}
		filesByExt[ext] = append(filesByExt[ext], file)
	}
	
	// Sort extensions by priority and count
	sortedExts := cg.sortExtensionsByPriority(filesByExt)
	
	for _, ext := range sortedExts {
		files := filesByExt[ext]
		if len(files) == 0 {
			continue
		}
		
		content.WriteString(fmt.Sprintf("## %s Files (%d files)\n\n", ext, len(files)))
		
		// Calculate statistics
		totalSize := int64(0)
		totalLines := 0
		for _, file := range files {
			totalSize += file.Size
			totalLines += file.Lines
		}
		
		content.WriteString(fmt.Sprintf("- **Total size:** %s\n", FormatSize(totalSize)))
		if totalLines > 0 {
			content.WriteString(fmt.Sprintf("- **Total lines:** %s\n", FormatNumber(totalLines)))
		}
		
		// List files (limit to reasonable number)
		content.WriteString("- **Files:**\n")
		maxFiles := 20
		if len(files) > maxFiles {
			content.WriteString(fmt.Sprintf("  (Showing %d of %d files)\n", maxFiles, len(files)))
		}
		
		for i, file := range files {
			if i >= maxFiles {
				break
			}
			relativePath := cg.getRelativePath(file.Path)
			content.WriteString(fmt.Sprintf("  - %s", relativePath))
			if file.Size > 1024 {
				content.WriteString(fmt.Sprintf(" (%s)", FormatSize(file.Size)))
			}
			if file.Lines > 0 {
				content.WriteString(fmt.Sprintf(" - %d lines", file.Lines))
			}
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}
	
	return ContextSection{
		Title:   "File Type Analysis",
		Content: content.String(),
		Files:   []string{},
	}
}

// generateContentSections creates sections with actual file content
func (cg *ContextGenerator) generateContentSections(scanResult *ScanResult) ([]ContextSection, error) {
	var sections []ContextSection
	
	// Select files to include based on priority and size constraints
	selectedFiles := cg.selectFilesForContent(scanResult.Files)
	
	// Group files by type for better organization
	filesByType := make(map[string][]FileInfo)
	for _, file := range selectedFiles {
		ext := file.Extension
		if ext == "" {
			ext = "other"
		}
		filesByType[ext] = append(filesByType[ext], file)
	}
	
	// Generate content sections for each file type
	for ext, files := range filesByType {
		section, err := cg.generateFileContentSection(ext, files)
		if err != nil {
			return nil, err
		}
		if section.Content != "" {
			sections = append(sections, section)
		}
	}
	
	return sections, nil
}

// generateFileContentSection creates a section with file contents for a specific type
func (cg *ContextGenerator) generateFileContentSection(extension string, files []FileInfo) (ContextSection, error) {
	var content strings.Builder
	var includedFiles []string
	
	sectionTitle := fmt.Sprintf("%s Files Content", strings.ToUpper(strings.TrimPrefix(extension, ".")))
	if extension == "other" {
		sectionTitle = "Other Files Content"
	}
	
	content.WriteString(fmt.Sprintf("# %s\n\n", sectionTitle))
	
	for _, file := range files {
		// Check size constraints
		if file.Size > cg.maxFileSize {
			continue
		}
		
		relativePath := cg.getRelativePath(file.Path)
		content.WriteString(fmt.Sprintf("## %s\n\n", relativePath))
		
		// Read file content
		fileContent, err := cg.readFileContent(file.Path)
		if err != nil {
			content.WriteString(fmt.Sprintf("*Error reading file: %v*\n\n", err))
			continue
		}
		
		// Add file content with syntax highlighting hint
		language := cg.getLanguageFromExtension(file.Extension)
		content.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", language, fileContent))
		
		includedFiles = append(includedFiles, relativePath)
		
		// Check total size constraint
		if int64(content.Len()) > cg.maxTotalSize {
			content.WriteString("*Context truncated due to size limits*\n\n")
			break
		}
	}
	
	return ContextSection{
		Title:   sectionTitle,
		Content: content.String(),
		Files:   includedFiles,
	}, nil
}

// selectFilesForContent selects which files to include in the content sections
func (cg *ContextGenerator) selectFilesForContent(files []FileInfo) []FileInfo {
	var selected []FileInfo
	
	// Priority scoring for files
	type scoredFile struct {
		file  FileInfo
		score int
	}
	
	var scoredFiles []scoredFile
	
	for _, file := range files {
		score := cg.calculateFileScore(file)
		if score > 0 {
			scoredFiles = append(scoredFiles, scoredFile{file: file, score: score})
		}
	}
	
	// Sort by score (highest first)
	sort.Slice(scoredFiles, func(i, j int) bool {
		return scoredFiles[i].score > scoredFiles[j].score
	})
	
	// Select files within size constraints
	totalSize := int64(0)
	for _, sf := range scoredFiles {
		if totalSize+sf.file.Size > cg.maxTotalSize {
			break
		}
		if sf.file.Size > cg.maxFileSize {
			continue
		}
		selected = append(selected, sf.file)
		totalSize += sf.file.Size
	}
	
	return selected
}

// calculateFileScore calculates a priority score for a file
func (cg *ContextGenerator) calculateFileScore(file FileInfo) int {
	score := 0
	
	// Base score for being a text file
	if cg.isTextFile(file.Extension) {
		score += 10
	}
	
	// Priority extension bonus
	for i, ext := range cg.priorityExtensions {
		if file.Extension == ext {
			score += 50 - i // Higher score for earlier extensions
			break
		}
	}
	
	// Size penalty (prefer smaller files)
	if file.Size < 1024 {
		score += 5
	} else if file.Size < 10*1024 {
		score += 3
	} else if file.Size > 100*1024 {
		score -= 5
	}
	
	// Important file names
	baseName := strings.ToLower(filepath.Base(file.Path))
	importantNames := []string{
		"readme", "main", "index", "app", "config", "package",
		"makefile", "dockerfile", "docker-compose",
	}
	
	for _, name := range importantNames {
		if strings.Contains(baseName, name) {
			score += 20
			break
		}
	}
	
	return score
}

// Helper functions

func (cg *ContextGenerator) isTextFile(ext string) bool {
	textExtensions := []string{
		".txt", ".md", ".go", ".js", ".ts", ".py", ".java", ".c", ".cpp",
		".h", ".hpp", ".cs", ".rb", ".php", ".html", ".css", ".scss",
		".json", ".xml", ".yaml", ".yml", ".toml", ".ini", ".cfg",
		".sh", ".bat", ".ps1", ".sql", ".r", ".scala", ".kt", ".rs",
	}
	
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	
	return false
}

func (cg *ContextGenerator) getLanguageFromExtension(ext string) string {
	langMap := map[string]string{
		".go":   "go",
		".js":   "javascript",
		".ts":   "typescript",
		".py":   "python",
		".java": "java",
		".c":    "c",
		".cpp":  "cpp",
		".html": "html",
		".css":  "css",
		".json": "json",
		".yaml": "yaml",
		".yml":  "yaml",
		".md":   "markdown",
		".sh":   "bash",
		".sql":  "sql",
	}
	
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}

func (cg *ContextGenerator) readFileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

func (cg *ContextGenerator) getRelativePath(fullPath string) string {
	// Try to get relative path, fallback to basename
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, fullPath); err == nil {
			return rel
		}
	}
	return filepath.Base(fullPath)
}

type ExtensionCount struct {
	Extension string
	Count     int
}

func (cg *ContextGenerator) sortExtensionsByCount(extensions map[string]int) []ExtensionCount {
	var sorted []ExtensionCount
	for ext, count := range extensions {
		sorted = append(sorted, ExtensionCount{Extension: ext, Count: count})
	}
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})
	
	return sorted
}

func (cg *ContextGenerator) sortExtensionsByPriority(filesByExt map[string][]FileInfo) []string {
	var extensions []string
	for ext := range filesByExt {
		extensions = append(extensions, ext)
	}
	
	sort.Slice(extensions, func(i, j int) bool {
		extI, extJ := extensions[i], extensions[j]
		
		// Priority extensions first
		priorityI, priorityJ := -1, -1
		for idx, priorityExt := range cg.priorityExtensions {
			if extI == priorityExt {
				priorityI = idx
			}
			if extJ == priorityExt {
				priorityJ = idx
			}
		}
		
		if priorityI != -1 && priorityJ != -1 {
			return priorityI < priorityJ
		}
		if priorityI != -1 {
			return true
		}
		if priorityJ != -1 {
			return false
		}
		
		// Then by count
		return len(filesByExt[extI]) > len(filesByExt[extJ])
	})
	
	return extensions
}

func (cg *ContextGenerator) buildDirectoryTree(files []FileInfo) string {
	// Build a simple directory tree representation
	var tree strings.Builder
	
	// Get unique directories
	dirs := make(map[string]bool)
	for _, file := range files {
		dir := filepath.Dir(file.Path)
		dirs[dir] = true
	}
	
	// Convert to sorted slice
	var sortedDirs []string
	for dir := range dirs {
		sortedDirs = append(sortedDirs, dir)
	}
	sort.Strings(sortedDirs)
	
	// Simple tree representation (first few levels)
	for i, dir := range sortedDirs {
		if i > 50 { // Limit output
			tree.WriteString("... (truncated)\n")
			break
		}
		
		relativePath := cg.getRelativePath(dir)
		depth := strings.Count(relativePath, string(filepath.Separator))
		indent := strings.Repeat("  ", depth)
		
		tree.WriteString(fmt.Sprintf("%s%s/\n", indent, filepath.Base(relativePath)))
	}
	
	return tree.String()
}

func (cg *ContextGenerator) generateSummary(scanResult *ScanResult, result *ContextResult) string {
	var summary strings.Builder
	
	summary.WriteString("## Context Summary\n\n")
	summary.WriteString(fmt.Sprintf("This context contains information about a project with %d files ", scanResult.TotalFiles))
	summary.WriteString(fmt.Sprintf("totaling %s across %d directories. ", FormatSize(scanResult.TotalSize), scanResult.TotalDirectories))
	
	if len(scanResult.Extensions) > 0 {
		// Find most common extension
		maxExt := ""
		maxCount := 0
		for ext, count := range scanResult.Extensions {
			if count > maxCount {
				maxCount = count
				maxExt = ext
			}
		}
		
		if maxExt != "" {
			name := maxExt
			if name == "" {
				name = "files without extension"
			}
			summary.WriteString(fmt.Sprintf("The project primarily consists of %s files (%d files). ", name, maxCount))
		}
	}
	
	summary.WriteString(fmt.Sprintf("The context includes %d sections with detailed information about the project structure and contents.", len(result.Sections)))
	
	return summary.String()
}

func (cg *ContextGenerator) estimateTokens(result *ContextResult) int {
	totalChars := 0
	
	for _, section := range result.Sections {
		totalChars += len(section.Content)
	}
	
	totalChars += len(result.Summary)
	
	// Rough estimate: 4 characters per token
	return totalChars / 4
}

func FormatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
}