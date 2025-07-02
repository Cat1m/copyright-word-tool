// processor.go - Enhanced version with file exclusion
package fileprocessor

import (
	"bufio"
	"copyright-code-word/config"
	"copyright-code-word/models"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileProcessor struct {
	config        *config.Config
	files         []models.CodeFile
	excludedCount int // ✅ Đếm số file bị exclude
}

func New(cfg *config.Config) *FileProcessor {
	return &FileProcessor{
		config:        cfg,
		files:         make([]models.CodeFile, 0),
		excludedCount: 0,
	}
}

func (fp *FileProcessor) ScanDirectory(rootDir string) ([]models.CodeFile, error) {
	fmt.Printf("🔍 Scanning for .cs and .dart files in: %s\n", rootDir)

	// ✅ In danh sách exclude để user biết
	fmt.Printf("🚫 File exclusion is enabled:\n")
	fp.config.PrintExcludeList()
	fmt.Println(strings.Repeat("-", 50))

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return fp.handleDirectory(d)
		}

		return fp.handleFile(path)
	})

	if err != nil {
		return nil, err
	}

	// Sort files by name
	sort.Slice(fp.files, func(i, j int) bool {
		return fp.files[i].FileName < fp.files[j].FileName
	})

	// ✅ In thống kê
	fp.printScanSummary()

	return fp.files, nil
}

func (fp *FileProcessor) handleDirectory(d fs.DirEntry) error {
	skipDirs := map[string]bool{
		"node_modules": true, ".git": true, "vendor": true,
		"target": true, "__pycache__": true, ".next": true,
		"build": true, "dist": true, "bin": true, "obj": true,
		".dart_tool": true, ".packages": true,
	}

	if skipDirs[d.Name()] {
		return filepath.SkipDir
	}
	return nil
}

func (fp *FileProcessor) handleFile(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	filename := filepath.Base(path)

	// ✅ Kiểm tra extension được hỗ trợ
	if !fp.config.SupportedExtensions[ext] {
		return nil
	}

	// ✅ Kiểm tra file có bị exclude không
	if fp.config.IsFileExcluded(filename) {
		fmt.Printf("🚫 Excluded: %s (sensitive file)\n", filename)
		fp.excludedCount++
		return nil
	}

	if err := fp.processFile(path, ext); err != nil {
		fmt.Printf("❌ Error processing %s: %v\n", path, err)
	} else {
		fmt.Printf("📄 Added: %s\n", filename)
	}

	return nil
}

func (fp *FileProcessor) processFile(filePath, ext string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string
	var content strings.Builder
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		content.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(lines) == 0 {
		fmt.Printf("⚠️  Skipped empty file: %s\n", filepath.Base(filePath))
		return nil
	}

	// Calculate page count
	totalLines := len(lines) + fp.config.CompactHeaderLines + fp.config.FileSeparatorLines
	pageCount := (totalLines + fp.config.LinesPerPage - 1) / fp.config.LinesPerPage
	if pageCount == 0 {
		pageCount = 1
	}

	fp.files = append(fp.files, models.CodeFile{
		FileName:  filepath.Base(filePath),
		Extension: ext,
		Lines:     lines,
		Content:   content.String(),
		PageCount: pageCount,
	})

	return nil
}

// ✅ Hàm in thống kê scan
func (fp *FileProcessor) printScanSummary() {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("📊 Scan Summary:\n")
	fmt.Printf("   ✅ Files included: %d\n", len(fp.files))
	fmt.Printf("   🚫 Files excluded: %d\n", fp.excludedCount)
	fmt.Printf("   📁 Total processed: %d\n", len(fp.files)+fp.excludedCount)

	if len(fp.files) > 0 {
		fmt.Printf("📋 Included files:\n")
		for _, file := range fp.files {
			fmt.Printf("   - %s (%d lines)\n", file.FileName, len(file.Lines))
		}
	}

	fmt.Println(strings.Repeat("=", 70))
}
