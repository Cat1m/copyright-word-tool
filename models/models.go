package models

type CodeFile struct {
	FileName  string
	Extension string
	Lines     []string
	Content   string
	PageCount int
}

type PageRange struct {
	FileIndex int
	StartLine int
	EndLine   int
	Pages     int
}
