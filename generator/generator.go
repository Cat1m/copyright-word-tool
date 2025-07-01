package generator

import (
	"copyright-code-word/config"
	"copyright-code-word/models"
	"copyright-code-word/paginator"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

type DocumentGenerator struct {
	config    *config.Config
	paginator *paginator.Paginator
}

func New(cfg *config.Config) *DocumentGenerator {
	return &DocumentGenerator{
		config:    cfg,
		paginator: paginator.New(cfg),
	}
}

func (dg *DocumentGenerator) InitializeLicense() error {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return err
	}

	err = license.SetMeteredKey(apiKey)
	if err != nil {
		return fmt.Errorf("license error: %v", err)
	}

	fmt.Println("✅ License activated successfully!")
	return nil
}

func (dg *DocumentGenerator) GenerateDocuments(files []models.CodeFile) error {
	if len(files) == 0 {
		return fmt.Errorf("no .cs or .dart files found")
	}

	totalPages := dg.paginator.CalculateTotalPages(files)
	dg.printStatistics(files, totalPages)

	if totalPages <= 100 {
		fmt.Printf("✅ ≤100 pages - Creating full document\n")
		return dg.createFullDocument(files)
	} else {
		fmt.Printf("⚠️  >100 pages - Creating 2 documents:\n")
		fmt.Printf("   - Full: %d pages\n", totalPages)
		fmt.Printf("   - Shortened: %d pages\n", dg.config.TargetPages)

		if err := dg.createFullDocument(files); err != nil {
			return err
		}
		return dg.createShortenedDocument(files)
	}
}

func (dg *DocumentGenerator) createFullDocument(files []models.CodeFile) error {
	doc := document.New()
	defer doc.Close()

	dg.setupPage(doc)
	dg.addAllFiles(doc, files)

	return dg.saveDocument(doc, "full_optimized")
}

func (dg *DocumentGenerator) createShortenedDocument(files []models.CodeFile) error {
	doc := document.New()
	defer doc.Close()

	dg.setupPage(doc)

	firstSection, middleStart, middleEnd, lastStart, totalLines := dg.paginator.CalculateContentSections(files)

	fmt.Printf("📝 Shortened sections:\n")
	fmt.Printf("   - Total content: %d lines\n", totalLines)
	fmt.Printf("   - First: lines 1-%d\n", firstSection)
	fmt.Printf("   - Middle: lines %d-%d\n", middleStart+1, middleEnd)
	fmt.Printf("   - Last: lines %d-%d\n", lastStart+1, totalLines)

	dg.addContentByLineRange(doc, files, 0, firstSection-1)
	dg.addContentByLineRange(doc, files, middleStart, middleEnd-1)
	dg.addContentByLineRange(doc, files, lastStart, totalLines-1)

	return dg.saveDocument(doc, "shortened_optimized")
}

func (dg *DocumentGenerator) setupPage(doc *document.Document) {
	section := doc.BodySection()
	section.SetPageSizeAndOrientation(
		measurement.Distance(210)*measurement.Millimeter,
		measurement.Distance(297)*measurement.Millimeter,
		wml.ST_PageOrientationPortrait,
	)
}

func (dg *DocumentGenerator) addAllFiles(doc *document.Document, files []models.CodeFile) {
	currentPageLines := 0

	for i, file := range files {
		fileHeaderLines := dg.config.CompactHeaderLines
		fileSeparatorLines := dg.config.FileSeparatorLines
		if i == len(files)-1 {
			fileSeparatorLines = 0
		}
		totalFileLinesNeeded := fileHeaderLines + len(file.Lines) + fileSeparatorLines

		// Smart page break
		if currentPageLines > dg.config.MinLinesForPageBreak &&
			currentPageLines+totalFileLinesNeeded > dg.config.LinesPerPage {

			breakPara := doc.AddParagraph()
			breakRun := breakPara.AddRun()
			breakRun.AddPageBreak()
			currentPageLines = 0

			fmt.Printf("🔄 Smart page break before %s\n", file.FileName)
		}

		dg.addFileToDocument(doc, file, i+1)
		currentPageLines += totalFileLinesNeeded

		if i < len(files)-1 {
			dg.addCompactFileSeparator(doc)
		}

		if currentPageLines >= dg.config.LinesPerPage {
			currentPageLines = currentPageLines % dg.config.LinesPerPage
		}
	}
}

func (dg *DocumentGenerator) addContentByLineRange(doc *document.Document, files []models.CodeFile, globalStartLine, globalEndLine int) {
	currentGlobalLine := 0

	for i, file := range files {
		fileStartLine := currentGlobalLine
		fileHeaderLines := dg.config.CompactHeaderLines
		fileContentLines := len(file.Lines)
		fileSeparatorLines := dg.config.FileSeparatorLines
		if i == len(files)-1 {
			fileSeparatorLines = 0
		}

		fileEndLine := fileStartLine + fileHeaderLines + fileContentLines + fileSeparatorLines - 1

		if fileEndLine >= globalStartLine && fileStartLine <= globalEndLine {
			fileLocalStartLine := 0
			if globalStartLine > fileStartLine+fileHeaderLines {
				fileLocalStartLine = globalStartLine - fileStartLine - fileHeaderLines
			}

			fileLocalEndLine := len(file.Lines) - 1
			if globalEndLine < fileStartLine+fileHeaderLines+fileContentLines-1 {
				fileLocalEndLine = globalEndLine - fileStartLine - fileHeaderLines
			}

			if globalStartLine <= fileStartLine+fileHeaderLines {
				dg.addCompactFileHeader(doc, file, i+1)
			}

			if fileLocalStartLine <= fileLocalEndLine && fileLocalEndLine >= 0 && fileLocalStartLine < len(file.Lines) {
				dg.addFileContentRange(doc, file, max(0, fileLocalStartLine), min(len(file.Lines)-1, fileLocalEndLine))
			}

			if i < len(files)-1 && globalEndLine >= fileEndLine-fileSeparatorLines {
				dg.addCompactFileSeparator(doc)
			}
		}

		currentGlobalLine = fileEndLine + 1
	}
}

func (dg *DocumentGenerator) addFileToDocument(doc *document.Document, file models.CodeFile, fileNumber int) {
	dg.addCompactFileHeader(doc, file, fileNumber)
	dg.addFileContentRange(doc, file, 0, len(file.Lines)-1)
}

func (dg *DocumentGenerator) addCompactFileHeader(doc *document.Document, file models.CodeFile, fileNumber int) {
	fileHeader := doc.AddParagraph()
	fileRun := fileHeader.AddRun()
	fileRun.AddText(fmt.Sprintf("📄 %s (%s, %d lines)",
		file.FileName,
		strings.ToUpper(file.Extension[1:]),
		len(file.Lines)))
	fileRun.Properties().SetBold(true)
	fileRun.Properties().SetSize(11)
	fileRun.Properties().SetColor(color.Blue)

	doc.AddParagraph()
}

func (dg *DocumentGenerator) addCompactFileSeparator(doc *document.Document) {
	separatorPara := doc.AddParagraph()
	separatorRun := separatorPara.AddRun()
	separatorRun.AddText(strings.Repeat("─", 60))
	separatorRun.Properties().SetSize(8)
	separatorRun.Properties().SetColor(color.LightGray)
}

func (dg *DocumentGenerator) addFileContentRange(doc *document.Document, file models.CodeFile, startLine, endLine int) {
	if startLine < 0 {
		startLine = 0
	}
	if endLine >= len(file.Lines) {
		endLine = len(file.Lines) - 1
	}

	for lineNum := startLine; lineNum <= endLine; lineNum++ {
		line := file.Lines[lineNum]

		codePara := doc.AddParagraph()

		// Line number
		lineNumRun := codePara.AddRun()
		lineNumRun.AddText(fmt.Sprintf("%4d │ ", lineNum+1))
		lineNumRun.Properties().SetFontFamily("Consolas")
		lineNumRun.Properties().SetSize(9)
		lineNumRun.Properties().SetColor(color.Gray)

		// Line content
		if len(line) > 120 {
			line = line[:120] + "..."
		}

		codeRun := codePara.AddRun()
		codeRun.AddText(line)
		codeRun.Properties().SetFontFamily("Consolas")
		codeRun.Properties().SetSize(9)
		codeRun.Properties().SetColor(color.Black)
	}
}

func (dg *DocumentGenerator) saveDocument(doc *document.Document, docType string) error {
	outputDir := "copyright_documents"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("source_code_%s_%s.docx", docType, timestamp)
	filepath := filepath.Join(outputDir, filename)

	if err := doc.SaveToFile(filepath); err != nil {
		return fmt.Errorf("failed to save document: %v", err)
	}

	fmt.Printf("✅ Created Word file: %s\n", filepath)
	return nil
}

func (dg *DocumentGenerator) printStatistics(files []models.CodeFile, totalPages int) {
	fmt.Printf("📊 Statistics (Optimized):\n")
	fmt.Printf("   - Files: %d\n", len(files))
	fmt.Printf("   - Total pages: %d (%d lines/page)\n", totalPages, dg.config.LinesPerPage)
	fmt.Printf("   - Details: ")
	for _, file := range files {
		fmt.Printf("%s(%dp) ", file.FileName, file.PageCount)
	}
	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
