package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func formatCommitHash(hash string) string {
	if len(hash) < 7 {
		return hash
	}
	return hash[:7]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-review <initial_commit> [final_commit] [--main-branch <branch_name>] [--project-path <path>]")
		os.Exit(1)
	}

	initialCommit := formatCommitHash(os.Args[1])
	finalCommit := "HEAD"
	mainBranch := ""
	projectPath := "."

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

	// Change to project directory
	err = os.Chdir(absProjectPath)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
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

	baseDir := "git-review"
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		fmt.Printf("Error creating base directory: %v\n", err)
		os.Exit(1)
	}

	dirName := filepath.Join(baseDir, fmt.Sprintf("%s-%s", initialCommit, finalCommit))
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Get list of modified files
	var cmd *exec.Cmd
	if mainBranch != "" {
		// Get files that were modified in the commit range but not in main branch
		cmd = exec.Command("git", "diff", "--name-only",
			fmt.Sprintf("%s...%s", initialCommit, finalCommit),
			"--not", mainBranch)
		cmd.Dir = absProjectPath
	} else {
		cmd = exec.Command("git", "diff", "--name-only", initialCommit, finalCommit)
		cmd.Dir = absProjectPath
	}
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Git error: %s\n", string(exitErr.Stderr))
		}
		fmt.Printf("Error executing git diff: %v\n", err)
		os.Exit(1)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		fmt.Println("No modified files found")
		os.Exit(0)
	}

	// For each file, extract and save the diff
	for _, file := range files {
		var diffCmd *exec.Cmd
		if mainBranch != "" {
			// Get diff excluding changes from main branch
			diffCmd = exec.Command("git", "diff",
				fmt.Sprintf("%s...%s", initialCommit, finalCommit),
				"--not", mainBranch,
				"--", file)
			diffCmd.Dir = absProjectPath
		} else {
			diffCmd = exec.Command("git", "diff", initialCommit, finalCommit, "--", file)
			diffCmd.Dir = absProjectPath
		}
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
