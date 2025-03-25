package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var AppVersion = "dev"

func showVersion() {
	fmt.Printf("git-review version %s\n", AppVersion)
}

func formatCommitHash(hash string) string {
	if len(hash) < 7 {
		return hash
	}
	return hash[:7]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-review <initial_commit> [final_commit] [--main-branch <branch_name>] [--project-path <path>] [--output-dir <path>]")
		os.Exit(1)
	}

	// Add version flag check
	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		showVersion()
		os.Exit(0)
	}

	initialCommit := formatCommitHash(os.Args[1])
	finalCommit := "HEAD"
	mainBranch := ""
	projectPath := "."
	outputDir := "git-review" // Default output directory

	// Parse arguments
	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--main-branch":
			if i+1 < len(os.Args) {
				mainBranch = os.Args[i+1]
				i++ // Skip next argument
			}
		case "--project-path":
			if i+1 < len(os.Args) {
				projectPath = os.Args[i+1]
				i++ // Skip next argument
			}
		case "--output-dir":
			if i+1 < len(os.Args) {
				outputDir = os.Args[i+1]
				i++ // Skip next argument
			}
		default:
			if finalCommit == "HEAD" {
				finalCommit = formatCommitHash(os.Args[i])
			}
		}
	}

	// Convert relative path to absolute if needed
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Printf("Error resolving project path: %v\n", err)
		os.Exit(1)
	}

	// If finalCommit is HEAD, we need to get the current hash
	if finalCommit == "HEAD" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = absProjectPath
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error getting HEAD hash: %v\n", err)
			os.Exit(1)
		}
		finalCommit = formatCommitHash(strings.TrimSpace(string(output)))
	}

	baseDir := outputDir

	// Check if a file exists with the same name as our intended directory
	fileInfo, err := os.Stat(baseDir)
	if err == nil {
		if !fileInfo.IsDir() {
			fmt.Printf("Error: '%s' already exists as a file. Please specify a different output directory using --output-dir\n", baseDir)
			os.Exit(1)
		}
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	dirName := filepath.Join(baseDir, fmt.Sprintf("%s-%s", initialCommit, finalCommit))
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Get list of modified files
	var files []string
	if mainBranch != "" {
		// First get all files modified between initial and final commit
		cmd := exec.Command("git", "diff", "--name-only",
			fmt.Sprintf("%s..%s", initialCommit, finalCommit))
		cmd.Dir = absProjectPath
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Git error: %s\n", string(exitErr.Stderr))
			}
			fmt.Printf("Error executing git diff: %v\n", err)
			os.Exit(1)
		}
		allFiles := strings.Split(strings.TrimSpace(string(output)), "\n")

		// For each file, check if it was also modified in the main branch
		for _, file := range allFiles {
			if file == "" {
				continue
			}

			// Check if this file has changes that are unique to our branch
			// by comparing our branch with the merge base against main branch
			mergeBaseCmd := exec.Command("git", "merge-base", mainBranch, finalCommit)
			mergeBaseCmd.Dir = absProjectPath
			mergeBase, err := mergeBaseCmd.Output()
			if err != nil {
				fmt.Printf("Error finding merge base: %v\n", err)
				continue
			}
			mergeBaseCommit := strings.TrimSpace(string(mergeBase))

			// Get changes between merge base and final commit
			diffCmd := exec.Command("git", "diff", "--name-only",
				fmt.Sprintf("%s..%s", mergeBaseCommit, finalCommit),
				"--", file)
			diffCmd.Dir = absProjectPath
			diff, err := diffCmd.Output()
			if err != nil {
				fmt.Printf("Error checking file %s: %v\n", file, err)
				continue
			}

			if strings.TrimSpace(string(diff)) != "" {
				// This file has changes that don't exist in main branch
				files = append(files, file)
			}
		}
	} else {
		// If no main branch specified, just get all changed files
		cmd := exec.Command("git", "diff", "--name-only", initialCommit, finalCommit)
		cmd.Dir = absProjectPath
		cmd.Stderr = os.Stderr
		output, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("Git error: %s\n", string(exitErr.Stderr))
			}
			fmt.Printf("Error executing git diff: %v\n", err)
			os.Exit(1)
		}
		allFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, file := range allFiles {
			if file != "" {
				files = append(files, file)
			}
		}
	}

	if len(files) == 0 {
		fmt.Println("No modified files found")
		os.Exit(0)
	}

	// For each file, extract and save the diff
	for _, file := range files {
		diffCmd := exec.Command("git", "diff", initialCommit, finalCommit, "--", file)
		diffCmd.Dir = absProjectPath
		diff, err := diffCmd.Output()
		if err != nil {
			fmt.Printf("Error getting diff for file %s: %v\n", file, err)
			continue
		}

		// Create output filename (replacing '/' with '_')
		safeFileName := strings.ReplaceAll(file, "/", "_")
		outputPath := filepath.Join(dirName, safeFileName+".diff")

		// Save diff to file
		err = os.WriteFile(outputPath, diff, 0644)
		if err != nil {
			fmt.Printf("Error saving diff for file %s: %v\n", file, err)
			continue
		}

		fmt.Printf("Diff saved for: %s\n", file)
	}

	fmt.Printf("\nAll diffs have been saved to directory: %s\n", dirName)
}
