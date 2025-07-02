// main.go - Enhanced version with file exclusion support
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

	// ✅ Xử lý arguments để thêm exclude files (nếu có)
	if len(os.Args) > 2 {
		handleAdditionalArgs(cfg, os.Args[2:])
	}

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

// ✅ Xử lý arguments để thêm exclude files
func handleAdditionalArgs(cfg *config.Config, args []string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--exclude=") {
			filename := strings.TrimPrefix(arg, "--exclude=")
			cfg.AddExcludeFile(filename)
			fmt.Printf("🚫 Added to exclude list: %s\n", filename)
		} else if strings.HasPrefix(arg, "--exclude-pattern=") {
			pattern := strings.TrimPrefix(arg, "--exclude-pattern=")
			cfg.AddExcludePattern(pattern)
			fmt.Printf("🚫 Added exclude pattern: *%s*\n", pattern)
		}
	}
}

func printUsage() {
	fmt.Println("📝 Go Code to Word - Optimized with File Exclusion (v2.1)")
	fmt.Println("")
	fmt.Println("Usage: go run main.go <directory_path> [options]")
	fmt.Println("Example: go run main.go ./src")
	fmt.Println("")
	fmt.Println("📂 Supported file types:")
	fmt.Println("  ✅ .cs (C#)")
	fmt.Println("  ✅ .dart (Dart)")
	fmt.Println("")
	fmt.Println("🚫 File Exclusion Options:")
	fmt.Println("  --exclude=filename          Exclude specific file (e.g., --exclude=program.cs)")
	fmt.Println("  --exclude-pattern=pattern   Exclude files containing pattern")
	fmt.Println("")
	fmt.Println("  Examples:")
	fmt.Println("    go run main.go ./src --exclude=program.cs --exclude=database.cs")
	fmt.Println("    go run main.go ./src --exclude-pattern=secret --exclude-pattern=config")
	fmt.Println("")
	fmt.Println("🚫 Default excluded files:")
	fmt.Println("  📄 Exact files: program.cs, appsettings.json, database.cs, secrets.cs...")
	fmt.Println("  🔍 Patterns: secret, password, apikey, config, setting, credential...")
	fmt.Println("")
	fmt.Println("🔑 Setup API Key (choose one):")
	fmt.Println("  📄 Create .env file:")
	fmt.Println("     UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("")
	fmt.Println("  🌍 Environment variable:")
	fmt.Println("     Windows PS: $env:UNIDOC_LICENSE_API_KEY=\"your_key\"")
	fmt.Println("     Windows CMD: set UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("     Linux/Mac: export UNIDOC_LICENSE_API_KEY=your_key")
	fmt.Println("")
	fmt.Println("  🆓 Register free: https://cloud.unidoc.io")
}

func printHeader(rootDir string, cfg *config.Config) {
	fmt.Printf("🚀 Creating optimized Word document with file exclusion (v2.1)...\n")
	fmt.Printf("📁 Source directory: %s\n", rootDir)
	fmt.Printf("📝 Processing: .cs (C#) and .dart (Dart)\n")
	fmt.Printf("📖 Optimization: %d lines/page, page break threshold: %d lines\n",
		cfg.LinesPerPage, cfg.MinLinesForPageBreak)
	fmt.Printf("🚫 File exclusion: enabled (%d files, %d patterns)\n",
		len(cfg.ExcludeFiles), len(cfg.ExcludePatterns))
	fmt.Printf("💡 Features: Compact header + minimal separator + smart page break + sensitive file filtering\n")
	fmt.Println(strings.Repeat("=", 70))
}

func printFooter() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("✨ Completed! Check 'copyright_documents' directory\n")
	fmt.Printf("💡 Word files have been optimized - saves 40-60%% paper!\n")
	fmt.Printf("🎯 Smart page break and sensitive file filtering applied\n")
	fmt.Printf("🔒 Sensitive files (config, secrets, etc.) were automatically excluded\n")
}
