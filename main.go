package main

import (
	"fmt"
	"os"

	"github.com/gifflet/git-review/cmd"
	"github.com/gifflet/git-review/internal/config"
)

// AppVersion is defined during compilation
var AppVersion = "dev"

func main() {
	// Initialize command flags before loading config
	if err := cmd.InitializeFlags(); err != nil {
		fmt.Printf("Error initializing flags: %v\n", err)
		os.Exit(1)
	}

	// Load configuration with project path
	cfg, err := config.LoadConfig(cmd.ProjectPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Set the application version
	cmd.AppVersion = AppVersion

	// Set config in cmd package
	cmd.SetConfig(cfg)

	// Execute the root command
	cmd.Execute()
}
