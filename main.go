package main

import (
	"copyright-code-word/config"
	"copyright-code-word/fileprocessor"
	"copyright-code-word/generator"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// ✅ Load .env file trước khi làm gì khác
	if err := config.LoadEnv(); err != nil {
		fmt.Printf("⚠️ Warning: %v\n", err)
	}

	rootDir := os.Args[1]

	// Validate directory
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		fmt.Printf("❌ Directory does not exist: %s\n", rootDir)
		os.Exit(1)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize components
	fileProcessor := fileprocessor.New(cfg)
	docGenerator := generator.New(cfg)

	// Initialize license (sẽ tự động đọc từ .env)
	if err := docGenerator.InitializeLicense(); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	printHeader(rootDir, cfg)

	// Process files
	files, err := fileProcessor.ScanDirectory(rootDir)
	if err != nil {
		fmt.Printf("❌ Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	// Generate documents
	if err := docGenerator.GenerateDocuments(files); err != nil {
		fmt.Printf("❌ Error generating document: %v\n", err)
		os.Exit(1)
	}

	printFooter()
}

func printUsage() {
	fmt.Println("📝 Go Code to Word - Optimized (v2.0)")
	fmt.Println("")
	fmt.Println("Usage: go run main.go <directory_path>")
	fmt.Println("Example: go run main.go ./src")
	fmt.Println("")
	fmt.Println("Supported file types:")
	fmt.Println("  ✅ .cs (C#)")
	fmt.Println("  ✅ .dart (Dart)")
	fmt.Println("")
	fmt.Println("Setup API Key (choose one):")
	fmt.Println("  🔑 Environment variable:")
	fmt.Println("     Windows PS: $env:UNIDOC_LICENSE_API_KEY=\"your_key\"")
	fmt.Println("     Windows CMD: set UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("     Linux/Mac: export UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("")
	fmt.Println("  📄 Or create .env file:")
	fmt.Println("     UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("")
	fmt.Println("  🆓 Register free: https://cloud.unidoc.io")
}

func printHeader(rootDir string, cfg *config.Config) {
	fmt.Printf("🚀 Creating optimized Word document for copyright registration (v2.0)...\n")
	fmt.Printf("📁 Source directory: %s\n", rootDir)
	fmt.Printf("📝 Processing: .cs (C#) and .dart (Dart)\n")
	fmt.Printf("📖 Optimization: %d lines/page, page break threshold: %d lines\n",
		cfg.LinesPerPage, cfg.MinLinesForPageBreak)
	fmt.Printf("💡 Features: Compact header + minimal separator + smart page break\n")
	fmt.Println(strings.Repeat("=", 70))
}

func printFooter() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("✨ Completed! Check 'copyright_documents' directory\n")
	fmt.Printf("💡 Word files have been optimized - saves 40-60%% paper!\n")
	fmt.Printf("🎯 Smart page break has been applied\n")
}
