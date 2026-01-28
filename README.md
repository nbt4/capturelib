# 📦 CaptureLib

Desktop application for managing Capture (.c2o) stage design files with drag & drop support.

## ✨ Features

- 📁 **Library Management** - Scan and index your .c2o files
- 🔍 **Fast Search** - Real-time file search
- 🎯 **Drag & Drop** - Drag files directly into Capture
- ⚡ **Auto-Scan** - Automatically detect new files
- 🎨 **Dark/Light Theme** - Choose your preferred theme
- 💾 **SQLite Database** - Fast and efficient file indexing

## 🚀 Quick Start

### Download

**Windows:** Download the latest `CaptureLib.exe` from [Releases](https://github.com/nbt4/capturelib/releases)

### First Run

1. Launch `CaptureLib.exe`
2. Click the folder icon to select your Capture library folder
3. The app will scan and index all .c2o files
4. Use the search bar to find files
5. Drag files from the grid into Capture

## 🛠️ Building from Source

### Prerequisites

- Go 1.21 or higher
- Windows/Linux/macOS with GUI support

### Build Commands

```bash
# Clone repository
git clone https://github.com/nbt4/capturelib.git
cd capturelib

# Install dependencies
go mod tidy

# Build for current platform
go build -o capturelib

# Build for Windows (on Linux/Mac)
GOOS=windows GOARCH=amd64 go build -o CaptureLib.exe

# Or use fyne-cross for cross-platform builds
go install github.com/fyne-io/fyne-cross@latest
fyne-cross windows -icon assets/icon.png
```

## 📖 Usage

### Keyboard Shortcuts

- `Ctrl+F` - Focus search bar
- `Ctrl+R` - Refresh library
- `Ctrl+,` - Open settings

### Settings

- **Library Path** - Where your .c2o files are stored
- **Auto-Scan** - Automatically scan on startup
- **Scan Subdirectories** - Include subdirectories in scan
- **Theme** - Dark or light mode

## 🏗️ Architecture

```
capturelib/
├── main.go           # Application entry point
├── ui/              
│   └── app.go        # UI logic and components
├── core/            
│   ├── library.go    # Library management
│   ├── scanner.go    # File scanner
│   ├── database.go   # SQLite database
│   └── config.go     # Configuration
└── models/          
    └── capture_file.go  # Data models
```

## 🔧 Technology Stack

- **Language:** Go 1.21+
- **UI Framework:** [Fyne](https://fyne.io/) v2.7+
- **Database:** SQLite (modernc.org/sqlite)
- **Config:** JSON

## 📝 Configuration

Config file location: `~/.capturelib/config.json`

```json
{
  "library_path": "C:\\Users\\YourName\\Capture\\Library",
  "theme": "dark",
  "auto_scan": true,
  "scan_subdirectories": true,
  "window_width": 1200,
  "window_height": 800
}
```

## 🐛 Troubleshooting

### "No files found"
- Check that your library path is correct
- Ensure .c2o files exist in the folder
- Try manual refresh (⟳ button)

### "Scan failed"
- Check folder permissions
- Ensure the path exists
- Check logs in `~/.capturelib/`

### Drag & Drop not working
- Ensure Capture is running
- Try dragging to desktop first, then to Capture
- Check that files aren't locked

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

MIT License - See LICENSE file for details

## 🙏 Credits

Built for Tsunami Events UG by Noah Tielmann
Powered by [Fyne](https://fyne.io/)

## 📞 Support

Issues: https://github.com/nbt4/capturelib/issues
