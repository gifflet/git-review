package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	initialCommit string
	finalCommit   string
	mainBranch    string
	projectPath   string
	outputDir     string

	// Application version
	AppVersion = "dev"

	// Root command
	rootCmd = &cobra.Command{
		Use:   "git-review",
		Short: "Git Review - Tool for analyzing differences between commits",
		Long: `Git Review is a powerful command-line tool designed to 
streamline code reviews by helping developers easily extract and analyze 
changes between Git commits.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no arguments, show help
			if len(args) == 0 && initialCommit == "" {
				cmd.Help()
				os.Exit(0)
			}

			// If there are positional arguments, use the first as initialCommit
			if len(args) > 0 && initialCommit == "" {
				initialCommit = args[0]
			}

			// If there's a second positional argument, use it as finalCommit
			if len(args) > 1 && finalCommit == "HEAD" {
				finalCommit = args[1]
			}

			// Execute the main function
			executeReview()
		},
		Version: AppVersion,
	}
)

// Execute adds all child commands to the root command and sets flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Flags
	rootCmd.Flags().StringVarP(&initialCommit, "initial", "i", "", "Starting commit hash for comparison (required)")
	rootCmd.Flags().StringVarP(&finalCommit, "final", "f", "HEAD", "Ending commit hash (defaults to HEAD)")
	rootCmd.Flags().StringVarP(&mainBranch, "main-branch", "m", "", "Main branch for refined comparisons")
	rootCmd.Flags().StringVarP(&projectPath, "project-path", "p", ".", "Project directory path")
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "git-review", "Output directory for diff files")

	// Mark initialCommit as required
	rootCmd.MarkFlagRequired("initial")
}
