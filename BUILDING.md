# 🔨 Building CaptureLib

## Windows (Super Easy!)

### Prerequisites
1. Install Go: https://go.dev/dl/ (Download the .msi installer)
2. Install Git: https://git-scm.com/download/win

### Build Steps

**Option A: One-Click Build** ⭐ EASIEST
1. Clone the repo:
   ```
   git clone https://github.com/nbt4/capturelib.git
   ```
2. Open the `capturelib` folder
3. **Double-click `build.bat`**
4. Done! `CaptureLib.exe` is ready

**Option B: Manual Build**
```cmd
git clone https://github.com/nbt4/capturelib.git
cd capturelib
go mod download
go build -ldflags="-H windowsgui -s -w" -o CaptureLib.exe
```

---

## Linux

### Prerequisites
```bash
sudo apt-get update
sudo apt-get install -y golang libgl1-mesa-dev xorg-dev
```

### Build
```bash
git clone https://github.com/nbt4/capturelib.git
cd capturelib
go mod download
go build -o capturelib
```

---

## macOS

### Prerequisites
```bash
brew install go
```

### Build
```bash
git clone https://github.com/nbt4/capturelib.git
cd capturelib
go mod download
go build -o capturelib
```

---

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Reinstall Go and make sure to check "Add to PATH"

### "fatal: unable to access ... SSL certificate problem"
- Git SSL issue
- Run: `git config --global http.sslVerify false` (temporary fix)

### Build takes forever
- First build downloads dependencies (~50MB)
- Subsequent builds are much faster

### "cannot find package"
- Run `go mod download` first
- Then try building again

---

## Build Flags Explained

- `-ldflags="-H windowsgui"` - Hide console window on Windows
- `-s -w` - Strip debug info (smaller .exe, ~30% reduction)
- `-o CaptureLib.exe` - Output file name

---

## File Size

Expected .exe size: **15-25 MB** (includes entire Go runtime + Fyne UI)

To make it smaller, use UPX:
```cmd
upx --best CaptureLib.exe
```
(Reduces to ~8-10 MB)
