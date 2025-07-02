// config.go - Enhanced version with file exclusion
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	LinesPerPage         int
	TargetPages          int
	SectionPages         int
	MinLinesForPageBreak int
	CompactHeaderLines   int
	FileSeparatorLines   int
	SupportedExtensions  map[string]bool
	// ✅ Thêm chức năng exclude files
	ExcludeFiles    map[string]bool // Exclude exact filename
	ExcludePatterns []string        // Exclude by pattern (contains)
}

func LoadConfig() *Config {
	return &Config{
		LinesPerPage:         70,
		TargetPages:          75,
		SectionPages:         26,
		MinLinesForPageBreak: 45,
		CompactHeaderLines:   2,
		FileSeparatorLines:   1,
		SupportedExtensions: map[string]bool{
			".cs":   true, // C#
			".dart": true, // Dart
		},
		// ✅ Files cần loại bỏ (exact match)
		ExcludeFiles: map[string]bool{
			"program.cs":              true, // Main entry point với config
			"appsettings.json":        true, // Config files
			"appsettings.local.json":  true,
			"web.config":              true,
			"app.config":              true,
			"database.cs":             true, // Database configs
			"connectionstrings.cs":    true,
			"secrets.cs":              true, // Any secrets
			"apikeys.cs":              true,
			"main.dart":               true, // Flutter main với sensitive config
			"AppConstants.cs":         true,
			"Utility.cs":              true,
			"APIKeyCheckAttribute.cs": true,
			"AppController.cs":        true,
			"Enum.cs":                 true,
		},
		// ✅ Patterns cần loại bỏ (contains match)
		ExcludePatterns: []string{
			"secret",     // Bất kỳ file nào chứa "secret"
			"password",   // Bất kỳ file nào chứa "password"
			"apikey",     // API keys
			"config",     // General config files
			"setting",    // Settings files
			"credential", // Credentials
			"token",      // Token files
			"env",
			"passcode", // Environment files
		},
	}
}

// ✅ Hàm kiểm tra file có bị exclude không
func (c *Config) IsFileExcluded(filename string) bool {
	// Chuẩn hóa filename về lowercase
	lowerFilename := strings.ToLower(filename)

	// 1. Kiểm tra exact match
	if c.ExcludeFiles[lowerFilename] {
		return true
	}

	// 2. Kiểm tra patterns
	for _, pattern := range c.ExcludePatterns {
		if strings.Contains(lowerFilename, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}

// ✅ Hàm thêm file exclude runtime (nếu cần)
func (c *Config) AddExcludeFile(filename string) {
	c.ExcludeFiles[strings.ToLower(filename)] = true
}

// ✅ Hàm thêm pattern exclude runtime (nếu cần)
func (c *Config) AddExcludePattern(pattern string) {
	c.ExcludePatterns = append(c.ExcludePatterns, pattern)
}

// ✅ Hàm in danh sách exclude để debug
func (c *Config) PrintExcludeList() {
	fmt.Printf("🚫 Excluded files (exact match):\n")
	for file := range c.ExcludeFiles {
		fmt.Printf("   - %s\n", file)
	}

	fmt.Printf("🚫 Excluded patterns (contains):\n")
	for _, pattern := range c.ExcludePatterns {
		fmt.Printf("   - *%s*\n", pattern)
	}
}

// Các hàm khác giữ nguyên...
func LoadEnv() error {
	pwd, _ := os.Getwd()
	fmt.Printf("🔍 Current working directory: %s\n", pwd)

	envPath := ".env"
	fmt.Printf("🔍 Looking for .env at: %s\n", filepath.Join(pwd, envPath))

	if data, err := os.ReadFile(envPath); err == nil {
		content := string(data)
		if strings.HasPrefix(content, "\ufeff") {
			content = strings.TrimPrefix(content, "\ufeff")
			fmt.Printf("🔧 Removed UTF-8 BOM from .env file\n")
		}

		tempFile := ".env.tmp"
		if err := os.WriteFile(tempFile, []byte(content), 0644); err == nil {
			defer os.Remove(tempFile)

			if err := godotenv.Load(tempFile); err == nil {
				fmt.Printf("✅ Loaded .env file successfully\n")
				return nil
			}
		}
	}

	fmt.Printf("💡 .env file not found or couldn't load, using system environment variables\n")
	return nil
}

func GetAPIKey() (string, error) {
	apiKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("UNIDOC_LICENSE_API_KEY environment variable is required\n\n" +
			"💡 Solutions:\n" +
			"   1. Set environment variable:\n" +
			"      Windows PowerShell: $env:UNIDOC_LICENSE_API_KEY=\"your_key\"\n" +
			"      Windows CMD: set UNIDOC_LICENSE_API_KEY=your_key\n" +
			"      Linux/Mac: export UNIDOC_LICENSE_API_KEY=your_key\n\n" +
			"   2. Create .env file with:\n" +
			"      UNIDOC_LICENSE_API_KEY=your_key\n\n" +
			"🆓 Get free key at: https://cloud.unidoc.io")
	}
	return apiKey, nil
}
