package ui

import (
	"fmt"
	"path/filepath"
	"github.com/nbt4/capturelib/core"
	"github.com/nbt4/capturelib/models"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/storage"
)

type App struct {
	fyneApp fyne.App
	window  fyne.Window
	library *core.Library
	
	// UI Components
	searchEntry *widget.Entry
	fileGrid    *fyne.Container
	statusLabel *widget.Label
	
	// Current files
	allFiles      []*models.CaptureFile
	displayedFiles []*models.CaptureFile
}

// NewApp creates a new application
func NewApp(configPath string) (*App, error) {
	// Create library
	lib, err := core.NewLibrary(configPath)
	if err != nil {
		return nil, err
	}
	
	// Create Fyne app
	fyneApp := app.NewWithID("com.tsunami.capturelib")
	
	a := &App{
		fyneApp: fyneApp,
		library: lib,
	}
	
	// Setup UI
	a.setupUI()
	
	return a, nil
}

func (a *App) setupUI() {
	a.window = a.fyneApp.NewWindow("Capture Library Manager")
	
	// Create UI components
	a.createToolbar()
	a.createSearchBar()
	a.createFileGrid()
	a.createStatusBar()
	
	// Layout
	content := container.NewBorder(
		container.NewVBox(a.createToolbar(), a.createSearchBar()),
		a.statusLabel,
		nil,
		nil,
		container.NewVScroll(a.fileGrid),
	)
	
	a.window.SetContent(content)
	a.window.Resize(fyne.NewSize(
		float32(a.library.GetConfig().WindowWidth),
		float32(a.library.GetConfig().WindowHeight),
	))
	
	// Load files
	a.refreshFiles()
}

func (a *App) createToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			a.selectLibraryFolder()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			a.rescanLibrary()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			a.showSettings()
		}),
	)
}

func (a *App) createSearchBar() *fyne.Container {
	a.searchEntry = widget.NewEntry()
	a.searchEntry.SetPlaceHolder("Search files...")
	a.searchEntry.OnChanged = func(query string) {
		a.filterFiles(query)
	}
	
	return container.NewBorder(nil, nil, widget.NewLabel("🔍"), nil, a.searchEntry)
}

func (a *App) createFileGrid() {
	a.fileGrid = container.NewGridWrap(fyne.NewSize(200, 120))
}

func (a *App) createStatusBar() {
	a.statusLabel = widget.NewLabel("Ready")
}

func (a *App) refreshFiles() {
	files, err := a.library.GetAllFiles()
	if err != nil {
		a.showError(fmt.Sprintf("Failed to load files: %v", err))
		return
	}
	
	a.allFiles = files
	a.displayedFiles = files
	a.updateFileGrid()
	a.updateStatus()
}

func (a *App) updateFileGrid() {
	a.fileGrid.Objects = nil
	
	for _, file := range a.displayedFiles {
		card := a.createFileCard(file)
		a.fileGrid.Add(card)
	}
	
	a.fileGrid.Refresh()
}

func (a *App) createFileCard(file *models.CaptureFile) fyne.CanvasObject {
	// Create card content
	info := fmt.Sprintf("%s\n%s", 
		formatSize(file.Size),
		file.ModifiedAt.Format("02.01.2006"))
	
	card := widget.NewCard(
		file.Filename,
		info,
		widget.NewIcon(theme.DocumentIcon()),
	)
	
	// Wrap in draggable
	draggable := &DraggableCard{
		Card:     card,
		FilePath: file.Filepath,
	}
	
	return draggable
}

func (a *App) filterFiles(query string) {
	if query == "" {
		a.displayedFiles = a.allFiles
	} else {
		files, err := a.library.SearchFiles(query)
		if err != nil {
			a.showError(fmt.Sprintf("Search failed: %v", err))
			return
		}
		a.displayedFiles = files
	}
	
	a.updateFileGrid()
	a.updateStatus()
}

func (a *App) selectLibraryFolder() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			return
		}
		
		path := uri.Path()
		if err := a.library.SetLibraryPath(path); err != nil {
			a.showError(fmt.Sprintf("Failed to set library path: %v", err))
			return
		}
		
		a.refreshFiles()
		a.showInfo(fmt.Sprintf("Library set to: %s", path))
	}, a.window)
}

func (a *App) rescanLibrary() {
	a.statusLabel.SetText("Scanning...")
	
	count, err := a.library.Scan()
	if err != nil {
		a.showError(fmt.Sprintf("Scan failed: %v", err))
		return
	}
	
	a.refreshFiles()
	a.showInfo(fmt.Sprintf("Scanned %d files", count))
}

func (a *App) showSettings() {
	cfg := a.library.GetConfig()
	
	// Theme selection
	themeSelect := widget.NewSelect([]string{"dark", "light"}, func(value string) {
		cfg.Theme = value
	})
	themeSelect.SetSelected(cfg.Theme)
	
	// Auto-scan checkbox
	autoScanCheck := widget.NewCheck("Auto-scan on startup", func(checked bool) {
		cfg.AutoScan = checked
	})
	autoScanCheck.SetChecked(cfg.AutoScan)
	
	// Subdirectories checkbox
	subdirsCheck := widget.NewCheck("Scan subdirectories", func(checked bool) {
		cfg.ScanSubdirectories = checked
	})
	subdirsCheck.SetChecked(cfg.ScanSubdirectories)
	
	// Library path
	pathLabel := widget.NewLabel(fmt.Sprintf("Library: %s", cfg.LibraryPath))
	
	form := container.NewVBox(
		pathLabel,
		widget.NewLabel("Theme:"),
		themeSelect,
		autoScanCheck,
		subdirsCheck,
	)
	
	dialog.ShowCustomConfirm("Settings", "Save", "Cancel", form, func(save bool) {
		if save {
			// Save config
			configPath := filepath.Join(filepath.Dir(cfg.LibraryPath), "..", "config.json")
			if err := a.library.SaveConfig(configPath); err != nil {
				a.showError(fmt.Sprintf("Failed to save settings: %v", err))
			} else {
				a.showInfo("Settings saved")
			}
		}
	}, a.window)
}

func (a *App) updateStatus() {
	count, _ := a.library.GetFileCount()
	totalSize := int64(0)
	for _, f := range a.allFiles {
		totalSize += f.Size
	}
	
	a.statusLabel.SetText(fmt.Sprintf("📁 Library: %s | 📊 %d files | %s",
		a.library.GetConfig().LibraryPath,
		count,
		formatSize(totalSize),
	))
}

func (a *App) showError(message string) {
	dialog.ShowError(fmt.Errorf(message), a.window)
}

func (a *App) showInfo(message string) {
	dialog.ShowInformation("Info", message, a.window)
}

// Run starts the application
func (a *App) Run() {
	a.window.ShowAndRun()
}

// Helper functions
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// DraggableCard is a draggable file card
type DraggableCard struct {
	*widget.Card
	FilePath string
}

// Draggable interface implementation
func (d *DraggableCard) Dragged(e *fyne.DragEvent) {
	// This will be implemented in v2 with proper file drag support
}

func (d *DraggableCard) DragEnd() {
	// Drag end handling
}
