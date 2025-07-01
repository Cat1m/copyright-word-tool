package paginator

import (
	"copyright-code-word/config"
	"copyright-code-word/models"
)

type Paginator struct {
	config *config.Config
}

func New(cfg *config.Config) *Paginator {
	return &Paginator{config: cfg}
}

func (p *Paginator) CalculateTotalPages(files []models.CodeFile) int {
	totalLines := 0

	for i, file := range files {
		// Header compact
		totalLines += p.config.CompactHeaderLines

		// File content
		totalLines += len(file.Lines)

		// Separator (except last file)
		if i < len(files)-1 {
			totalLines += p.config.FileSeparatorLines
		}
	}

	return (totalLines + p.config.LinesPerPage - 1) / p.config.LinesPerPage
}

func (p *Paginator) CalculatePageRanges(files []models.CodeFile) []models.PageRange {
	pageRanges := make([]models.PageRange, 0)
	currentPageLines := 0

	for fileIndex, file := range files {
		linesRemaining := len(file.Lines)
		startLine := 0

		for linesRemaining > 0 {
			linesToAdd := min(p.config.LinesPerPage-currentPageLines, linesRemaining)

			if currentPageLines > 0 &&
				linesToAdd < p.config.MinLinesForPageBreak &&
				linesRemaining > p.config.MinLinesForPageBreak {
				currentPageLines = 0
				linesToAdd = min(p.config.LinesPerPage, linesRemaining)
			}

			endLine := startLine + linesToAdd - 1
			if endLine >= len(file.Lines) {
				endLine = len(file.Lines) - 1
			}

			pageRanges = append(pageRanges, models.PageRange{
				FileIndex: fileIndex,
				StartLine: startLine,
				EndLine:   endLine,
				Pages:     1,
			})

			startLine = endLine + 1
			linesRemaining = len(file.Lines) - startLine
			currentPageLines = (currentPageLines + linesToAdd) % p.config.LinesPerPage
		}
	}

	return pageRanges
}

func (p *Paginator) CalculateContentSections(files []models.CodeFile) (firstSection, middleStart, middleEnd, lastStart, totalLines int) {
	totalLines = p.calculateTotalContentLines(files)
	linesPerSection := (p.config.TargetPages * p.config.LinesPerPage) / 3

	firstSection = min(linesPerSection, totalLines)
	middleStart = max(0, (totalLines/2)-(linesPerSection/2))
	middleEnd = min(totalLines, middleStart+linesPerSection)
	lastStart = max(0, totalLines-linesPerSection)

	return
}

func (p *Paginator) calculateTotalContentLines(files []models.CodeFile) int {
	totalLines := 0

	for i, file := range files {
		totalLines += p.config.CompactHeaderLines
		totalLines += len(file.Lines)

		if i < len(files)-1 {
			totalLines += p.config.FileSeparatorLines
		}
	}

	return totalLines
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
