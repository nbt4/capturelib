@echo off
echo ========================================
echo   CaptureLib Windows Builder
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed!
    echo.
    echo Please install Go from: https://go.dev/dl/
    echo.
    pause
    exit /b 1
)

echo [1/3] Checking Go version...
go version

echo.
echo [2/3] Downloading dependencies...
go mod download

echo.
echo [3/3] Building CaptureLib.exe...
go build -ldflags="-H windowsgui -s -w" -o CaptureLib.exe

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo   BUILD SUCCESSFUL!
    echo ========================================
    echo.
    echo CaptureLib.exe is ready to use!
    echo.
    dir CaptureLib.exe
) else (
    echo.
    echo ========================================
    echo   BUILD FAILED!
    echo ========================================
    echo.
    echo Check the error messages above.
)

echo.
pause
