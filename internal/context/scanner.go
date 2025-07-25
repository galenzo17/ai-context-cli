package context

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileInfo represents information about a scanned file
type FileInfo struct {
	Path         string
	Size         int64
	Lines        int
	Extension    string
	IsDirectory  bool
	ModTime      time.Time
	IsExcluded   bool
	ExcludeReason string
}

// ScanResult represents the result of a project scan
type ScanResult struct {
	TotalFiles      int
	TotalDirectories int
	TotalSize       int64
	TotalLines      int
	ExcludedFiles   int
	ScanDuration    time.Duration
	Files           []FileInfo
	Extensions      map[string]int
	LargestFiles    []FileInfo
}

// ScanConfig holds configuration for the scanner
type ScanConfig struct {
	RootPath        string
	ExcludePatterns []string
	ExcludeExtensions []string
	MaxDepth        int
	MaxFileSize     int64 // in bytes
	IncludeHidden   bool
	FollowSymlinks  bool
}

// DefaultScanConfig returns a sensible default configuration
func DefaultScanConfig(rootPath string) ScanConfig {
	return ScanConfig{
		RootPath: rootPath,
		ExcludePatterns: []string{
			"node_modules/**",
			".git/**",
			"vendor/**",
			"dist/**",
			"build/**",
			"*.log",
			"*.tmp",
			"*.cache",
			".DS_Store",
			"Thumbs.db",
		},
		ExcludeExtensions: []string{
			".exe", ".dll", ".so", ".dylib",
			".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico",
			".mp3", ".mp4", ".avi", ".mkv", ".mov",
			".zip", ".tar", ".gz", ".rar", ".7z",
			".pdf", ".doc", ".docx", ".xls", ".xlsx",
		},
		MaxDepth:       50,
		MaxFileSize:    10 * 1024 * 1024, // 10MB
		IncludeHidden:  false,
		FollowSymlinks: false,
	}
}

// ProjectScanner handles scanning project directories
type ProjectScanner struct {
	config   ScanConfig
	progress chan ScanProgress
	cancel   chan bool
}

// ScanProgress represents progress during scanning
type ScanProgress struct {
	CurrentFile     string
	ProcessedFiles  int
	TotalEstimated  int
	CurrentPhase    string
	ElapsedTime     time.Duration
}

// NewProjectScanner creates a new project scanner
func NewProjectScanner(config ScanConfig) *ProjectScanner {
	return &ProjectScanner{
		config:   config,
		progress: make(chan ScanProgress, 100),
		cancel:   make(chan bool, 1),
	}
}

// Scan performs a full project scan
func (ps *ProjectScanner) Scan() (*ScanResult, error) {
	startTime := time.Now()
	
	result := &ScanResult{
		Files:      make([]FileInfo, 0),
		Extensions: make(map[string]int),
	}
	
	// Send initial progress
	ps.sendProgress(ScanProgress{
		CurrentPhase: "Initializing scan...",
		ElapsedTime:  time.Since(startTime),
	})
	
	// First pass: count files for progress estimation
	estimatedFiles := ps.estimateFileCount()
	
	ps.sendProgress(ScanProgress{
		CurrentPhase:   fmt.Sprintf("Scanning %d estimated files...", estimatedFiles),
		TotalEstimated: estimatedFiles,
		ElapsedTime:    time.Since(startTime),
	})
	
	// Second pass: actual scanning
	err := ps.scanDirectory(ps.config.RootPath, 0, result, startTime, estimatedFiles)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	
	// Post-process results
	result.ScanDuration = time.Since(startTime)
	ps.processResults(result)
	
	ps.sendProgress(ScanProgress{
		CurrentPhase:   "Scan completed!",
		ProcessedFiles: result.TotalFiles,
		TotalEstimated: estimatedFiles,
		ElapsedTime:    result.ScanDuration,
	})
	
	return result, nil
}

// GetProgressChannel returns the progress channel
func (ps *ProjectScanner) GetProgressChannel() <-chan ScanProgress {
	return ps.progress
}

// Cancel stops the scanning process
func (ps *ProjectScanner) Cancel() {
	select {
	case ps.cancel <- true:
	default:
	}
}

// estimateFileCount provides a rough estimate of files to scan
func (ps *ProjectScanner) estimateFileCount() int {
	count := 0
	filepath.WalkDir(ps.config.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Continue on errors during estimation
		}
		
		if ps.shouldExcludePath(path, d.IsDir()) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		
		if !d.IsDir() {
			count++
		}
		
		// Limit estimation time
		if count > 10000 {
			return fs.SkipAll
		}
		
		return nil
	})
	
	return count
}

// scanDirectory recursively scans a directory
func (ps *ProjectScanner) scanDirectory(dirPath string, depth int, result *ScanResult, startTime time.Time, totalEstimated int) error {
	if depth > ps.config.MaxDepth {
		return nil
	}
	
	// Check for cancellation
	select {
	case <-ps.cancel:
		return fmt.Errorf("scan cancelled")
	default:
	}
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}
	
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		
		// Send progress update
		ps.sendProgress(ScanProgress{
			CurrentFile:    fullPath,
			ProcessedFiles: result.TotalFiles + result.ExcludedFiles,
			TotalEstimated: totalEstimated,
			CurrentPhase:   "Scanning files...",
			ElapsedTime:    time.Since(startTime),
		})
		
		fileInfo := ps.scanFile(fullPath, entry)
		
		if entry.IsDir() {
			result.TotalDirectories++
			if !fileInfo.IsExcluded {
				// Recurse into subdirectory
				err := ps.scanDirectory(fullPath, depth+1, result, startTime, totalEstimated)
				if err != nil {
					return err
				}
			}
		} else {
			if fileInfo.IsExcluded {
				result.ExcludedFiles++
			} else {
				result.TotalFiles++
				result.TotalSize += fileInfo.Size
				result.TotalLines += fileInfo.Lines
				result.Extensions[fileInfo.Extension]++
				result.Files = append(result.Files, fileInfo)
			}
		}
	}
	
	return nil
}

// scanFile scans an individual file
func (ps *ProjectScanner) scanFile(path string, entry fs.DirEntry) FileInfo {
	info, err := entry.Info()
	if err != nil {
		return FileInfo{
			Path:          path,
			IsDirectory:   entry.IsDir(),
			IsExcluded:    true,
			ExcludeReason: fmt.Sprintf("Cannot read file info: %v", err),
		}
	}
	
	fileInfo := FileInfo{
		Path:        path,
		Size:        info.Size(),
		IsDirectory: entry.IsDir(),
		ModTime:     info.ModTime(),
		Extension:   strings.ToLower(filepath.Ext(path)),
	}
	
	// Check exclusion rules
	if ps.shouldExcludePath(path, entry.IsDir()) {
		fileInfo.IsExcluded = true
		fileInfo.ExcludeReason = "Matches exclude pattern"
		return fileInfo
	}
	
	// Check file size limit
	if !entry.IsDir() && info.Size() > ps.config.MaxFileSize {
		fileInfo.IsExcluded = true
		fileInfo.ExcludeReason = fmt.Sprintf("File too large (%d bytes)", info.Size())
		return fileInfo
	}
	
	// Count lines for text files
	if !entry.IsDir() && ps.isTextFile(fileInfo.Extension) {
		lines, err := ps.countLines(path)
		if err == nil {
			fileInfo.Lines = lines
		}
	}
	
	return fileInfo
}

// shouldExcludePath checks if a path should be excluded
func (ps *ProjectScanner) shouldExcludePath(path string, isDir bool) bool {
	// Check hidden files/directories
	if !ps.config.IncludeHidden {
		if strings.HasPrefix(filepath.Base(path), ".") {
			return true
		}
	}
	
	// Check extension exclusions
	if !isDir {
		ext := strings.ToLower(filepath.Ext(path))
		for _, excludeExt := range ps.config.ExcludeExtensions {
			if ext == excludeExt {
				return true
			}
		}
	}
	
	// Check pattern exclusions
	for _, pattern := range ps.config.ExcludePatterns {
		// Handle directory patterns like "node_modules/**"
		if strings.Contains(pattern, "/**") {
			dirPattern := strings.TrimSuffix(pattern, "/**")
			if strings.Contains(path, dirPattern) {
				return true
			}
		}
		
		// Handle simple file patterns
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
		
		// Handle full path patterns
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
	}
	
	return false
}

// isTextFile determines if a file is likely a text file
func (ps *ProjectScanner) isTextFile(ext string) bool {
	textExtensions := []string{
		".txt", ".md", ".go", ".js", ".ts", ".py", ".java", ".c", ".cpp",
		".h", ".hpp", ".cs", ".rb", ".php", ".html", ".css", ".scss",
		".json", ".xml", ".yaml", ".yml", ".toml", ".ini", ".cfg",
		".sh", ".bat", ".ps1", ".sql", ".r", ".scala", ".kt", ".rs",
		".jsx", ".tsx", ".vue", ".svelte", ".dart", ".swift", ".m",
	}
	
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	
	return ext == "" // Files without extension might be text
}

// countLines counts the number of lines in a file
func (ps *ProjectScanner) countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
		// Prevent very large files from slowing down the scan
		if lines > 100000 {
			break
		}
	}
	
	return lines, scanner.Err()
}

// processResults post-processes scan results
func (ps *ProjectScanner) processResults(result *ScanResult) {
	// Sort files by size (largest first) for LargestFiles
	sortedFiles := make([]FileInfo, len(result.Files))
	copy(sortedFiles, result.Files)
	
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Size > sortedFiles[j].Size
	})
	
	// Keep top 10 largest files
	maxLargest := 10
	if len(sortedFiles) < maxLargest {
		maxLargest = len(sortedFiles)
	}
	result.LargestFiles = sortedFiles[:maxLargest]
}

// sendProgress sends a progress update
func (ps *ProjectScanner) sendProgress(progress ScanProgress) {
	select {
	case ps.progress <- progress:
	default:
		// Don't block if channel is full
	}
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

// EstimateProcessingTime estimates how long it will take to process files
func EstimateProcessingTime(fileCount int) time.Duration {
	// Rough estimate: 1ms per file for scanning + processing
	baseTime := time.Duration(fileCount) * time.Millisecond
	
	// Add overhead for I/O and processing
	overhead := time.Duration(float64(baseTime) * 0.5)
	
	return baseTime + overhead
}