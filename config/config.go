package config

import (
	"fmt"
	"os"
	"strings"

	"path/filepath"

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
	}
}

// ✅ Hàm mới: Load .env file
func LoadEnv() error {
	// Debug: Kiểm tra current working directory
	pwd, _ := os.Getwd()
	fmt.Printf("🔍 Current working directory: %s\n", pwd)

	envPath := ".env"
	fmt.Printf("🔍 Looking for .env at: %s\n", filepath.Join(pwd, envPath))

	// Đọc file .env và loại bỏ BOM nếu có
	if data, err := os.ReadFile(envPath); err == nil {
		// Loại bỏ UTF-8 BOM nếu có
		content := string(data)
		if strings.HasPrefix(content, "\ufeff") {
			content = strings.TrimPrefix(content, "\ufeff")
			fmt.Printf("🔧 Removed UTF-8 BOM from .env file\n")
		}

		// Tạo file tạm không có BOM
		tempFile := ".env.tmp"
		if err := os.WriteFile(tempFile, []byte(content), 0644); err == nil {
			defer os.Remove(tempFile) // Xóa file tạm sau khi dùng

			// Load file tạm
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
