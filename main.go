package main

import (
	"github.com/gifflet/git-review/cmd"
)

// AppVersion is defined during compilation
var AppVersion = "dev"

func main() {
	// Set the application version
	cmd.AppVersion = AppVersion

	// Execute the root command
	cmd.Execute()
}
