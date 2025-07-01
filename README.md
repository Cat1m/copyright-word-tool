# 📝 Go Code to Word - Copyright Tool

> **Công cụ chuyển đổi source code C# và Dart sang Word Document để đăng ký bản quyền**

Tối ưu hóa không gian trang, tiết kiệm 40-60% giấy in với thuật toán phân trang thông minh!

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🎯 Tính năng chính

- ✅ **Hỗ trợ ngôn ngữ**: C# (.cs), Dart (.dart)
- 📄 **Tối ưu phân trang**: 70 dòng/trang với page break thông minh
- 📊 **Dual output**: File đầy đủ + file rút gọn (cho project >100 trang)
- 🎨 **Format đẹp**: Header compact, line numbers, màu sắc phân biệt
- 🚀 **Hiệu suất cao**: Scan nhanh, bỏ qua thư mục không cần thiết
- 🔒 **Bảo mật**: API key từ environment variable
- 📁 **Tự động tổ chức**: Output có timestamp, thư mục riêng

## 📋 Yêu cầu hệ thống

- **Go**: Version 1.21+ ([Tải tại đây](https://golang.org/dl/))
- **OS**: Windows, macOS, Linux
- **API Key**: UniDoc miễn phí ([Đăng ký tại đây](https://cloud.unidoc.io))

## 🚀 Hướng dẫn cài đặt

### Bước 1: Clone repository về máy
```bash
git clone https://github.com/your-username/copyright-word-tool.git
cd copyright-word-tool
```

### Bước 2: Cài đặt dependencies
```bash
go mod download
go mod tidy
```

### Bước 3: Lấy API Key miễn phí
1. 🌐 Truy cập [cloud.unidoc.io](https://cloud.unidoc.io)
2. 📝 Đăng ký account miễn phí (chỉ cần email)
3. 🔑 Copy API key từ dashboard
4. 💾 Lưu lại để dùng ở bước tiếp theo

### Bước 4: Setup API Key

**🎯 Cách 1: File .env (Đơn giản nhất)**
```bash
# Tạo file .env ở thư mục gốc
echo "UNIDOC_LICENSE_API_KEY=paste_api_key_của_bạn_vào_đây" > .env
```

**⚙️ Cách 2: Environment Variable**

**Windows PowerShell:**
```powershell
$env:UNIDOC_LICENSE_API_KEY="paste_api_key_của_bạn_vào_đây"
```

**Windows Command Prompt:**
```cmd
set UNIDOC_LICENSE_API_KEY=paste_api_key_của_bạn_vào_đây
```

**Linux/macOS:**
```bash
export UNIDOC_LICENSE_API_KEY=paste_api_key_của_bạn_vào_đây
```

### Bước 5: Kiểm tra cài đặt
```bash
# Test với thư mục hiện tại
go run main.go .

# Nếu thấy "✅ License activated successfully!" là thành công!
```

## 💻 Cách sử dụng

### Syntax cơ bản
```bash
go run main.go <đường_dẫn_thư_mục_source_code>
```

### 📋 Ví dụ thực tế

**Xử lý project Flutter:**
```bash
go run main.go D:\GitHub\flutter\MyFlutterApp
go run main.go /Users/john/Projects/flutter_app
```

**Xử lý project C#:**
```bash
go run main.go C:\Source\MyDotNetProject
go run main.go /home/user/dotnet-project
```

**Xử lý thư mục hiện tại:**
```bash
go run main.go .
```

**Xử lý thư mục con:**
```bash
go run main.go ./src
go run main.go ../OtherProject
```

## 📊 Kết quả và Output

### 📁 File output được tạo tại:
```
copyright_documents/
├── source_code_full_optimized_20250701_143022.docx      # File đầy đủ
└── source_code_shortened_optimized_20250701_143022.docx # File rút gọn (nếu >100 trang)
```

### 🎯 Logic tạo file:
- **≤100 trang**: Chỉ tạo file **full** (đầy đủ toàn bộ code)
- **>100 trang**: Tạo cả file **full** + **shortened** (75 trang: đầu + giữa + cuối)

### 📄 Nội dung file shortened:
- **25 trang đầu**: Code từ đầu project
- **25 trang giữa**: Code từ giữa project  
- **25 trang cuối**: Code từ cuối project
- **= 75 trang tổng cộng** (phù hợp đăng ký bản quyền)

## 🔧 Tùy chỉnh nâng cao

### Thay đổi cấu hình trong `config/config.go`:

```go
// Số dòng mỗi trang
LinesPerPage: 70,  // Tăng thành 80 để ít trang hơn

// Số trang file rút gọn
TargetPages: 75,   // Tăng thành 100 nếu muốn nhiều code hơn

// Ngưỡng page break thông minh
MinLinesForPageBreak: 45,  // Giảm xuống 30 để ít page break hơn
```

### Thêm hỗ trợ file type mới:
```go
SupportedExtensions: map[string]bool{
    ".cs":   true,  // C#
    ".dart": true,  // Dart
    ".go":   true,  // Go (thêm mới)
    ".py":   true,  // Python (thêm mới)
    ".js":   true,  // JavaScript (thêm mới)
    ".ts":   true,  // TypeScript (thêm mới)
    ".java": true,  // Java (thêm mới)
}
```

## 🐛 Xử lý lỗi thường gặp

### ❌ "UNIDOC_LICENSE_API_KEY environment variable is required"
**Nguyên nhân**: Chưa setup API key  
**Giải pháp**: Làm theo Bước 3 và 4 ở trên

### ❌ "Directory does not exist"
**Nguyên nhân**: Đường dẫn thư mục sai  
**Giải pháp**: Kiểm tra lại đường dẫn, dùng dấu `/` thay vì `\` trên Linux/Mac

### ❌ "License error: invalid API key"
**Nguyên nhân**: API key sai hoặc hết hạn  
**Giải pháp**: Lấy API key mới tại [cloud.unidoc.io](https://cloud.unidoc.io)

### ❌ "no .cs or .dart files found"
**Nguyên nhân**: Thư mục không có file C# hoặc Dart  
**Giải pháp**: Kiểm tra lại thư mục hoặc thêm file type mới vào config

## 📋 Demo Output

```
🚀 Creating optimized Word document for copyright registration (v2.0)...
📁 Source directory: D:\GitHub\flutter\MyApp
📝 Processing: .cs (C#) and .dart (Dart)
📖 Optimization: 70 lines/page, page break threshold: 45 lines
💡 Features: Compact header + minimal separator + smart page break
======================================================================
🔍 Scanning for .cs and .dart files in: D:\GitHub\flutter\MyApp
✅ License activated successfully!
📄 Added: main.dart
📄 Added: home_screen.dart
📄 Added: user_service.dart
🔄 Smart page break before user_service.dart
📊 Statistics (Optimized):
   - Files: 3
   - Total pages: 45 (70 lines/page)
   - Details: main.dart(12p) home_screen.dart(18p) user_service.dart(15p)
✅ ≤100 pages - Creating full document
✅ Created Word file: copyright_documents/source_code_full_optimized_20250701_143022.docx
======================================================================
✨ Completed! Check 'copyright_documents' directory
💡 Word files have been optimized - saves 40-60% paper!
🎯 Smart page break has been applied
```

## 🏗️ Build và phân phối

### Build executable file:
```bash
# Build cho Windows
go build -o copyright-tool.exe main.go

# Build cho Linux
GOOS=linux go build -o copyright-tool main.go

# Build cho macOS
GOOS=darwin go build -o copyright-tool main.go
```

### Sử dụng executable:
```bash
# Windows
./copyright-tool.exe D:\GitHub\flutter\MyProject

# Linux/Mac
./copyright-tool /home/user/my-project
```

## 🤝 Đóng góp

Nếu bạn muốn thêm tính năng hoặc sửa bug:
1. Fork repository này
2. Tạo branch mới: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Mở Pull Request

## 📜 License

Distributed under the MIT License. See `LICENSE` file for more information.

<!-- ## 🆘 Hỗ trợ

- 🐛 **Bug reports**: [Issues](https://github.com/your-username/copyright-word-tool/issues)
- 💡 **Feature requests**: [Discussions](https://github.com/your-username/copyright-word-tool/discussions)
- 📧 **Email**: your-email@domain.com

## 🏆 Tác giả

**Your Name** - [@your-twitter](https://twitter.com/your-twitter) - your-email@domain.com

Project Link: [https://github.com/your-username/copyright-word-tool](https://github.com/your-username/copyright-word-tool) -->

---

⭐ **Nếu tool này hữu ích, hãy cho 1 star để ủng hộ nhé!** ⭐