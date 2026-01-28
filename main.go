package main

import (
	"log"
	"os"
	"path/filepath"
	
	"github.com/nbt4/capturelib/ui"
)

func main() {
	// Get config path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	
	configDir := filepath.Join(homeDir, ".capturelib")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	
	configPath := filepath.Join(configDir, "config.json")
	
	// Create and run app
	app, err := ui.NewApp(configPath)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	
	app.Run()
}
