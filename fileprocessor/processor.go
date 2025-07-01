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
	config *config.Config
	files  []models.CodeFile
}

func New(cfg *config.Config) *FileProcessor {
	return &FileProcessor{
		config: cfg,
		files:  make([]models.CodeFile, 0),
	}
}

func (fp *FileProcessor) ScanDirectory(rootDir string) ([]models.CodeFile, error) {
	fmt.Printf("🔍 Scanning for .cs and .dart files in: %s\n", rootDir)

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
	if !fp.config.SupportedExtensions[ext] {
		return nil
	}

	if err := fp.processFile(path, ext); err != nil {
		fmt.Printf("❌ Error processing %s: %v\n", path, err)
	} else {
		fmt.Printf("📄 Added: %s\n", filepath.Base(path))
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
