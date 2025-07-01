package config

import (
	"fmt"
	"os"

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
	// Tìm .env file từ current directory
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Nếu không có .env, thử tìm ở parent directories
		for i := 0; i < 3; i++ { // Tìm tối đa 3 level lên
			envPath = filepath.Join("..", envPath)
			if _, err := os.Stat(envPath); err == nil {
				break
			}
		}
	}

	// Load .env file (không báo lỗi nếu không tìm thấy)
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("💡 .env file not found at %s, using system environment variables\n", envPath)
	} else {
		fmt.Printf("✅ Loaded .env file from %s\n", envPath)
	}

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
