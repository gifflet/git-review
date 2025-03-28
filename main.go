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
	// Load configuration
	cfg, err := config.LoadConfig()
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
