<div align="center">
  <picture>
    <img src="docs/logo.png" width="256" alt="Git Review Logo"/>
  </picture>

  # Git Review 🚀

  ![Go Version](https://img.shields.io/badge/Go-1.20+-blue)
  ![License](https://img.shields.io/badge/License-MIT-green)
</div>

## Overview

Git Review is a powerful command-line tool designed to streamline code reviews by helping developers easily extract and analyze changes between Git commits.

## Prerequisites 🔍

- Git (version 2.25 or higher)
- Go (version 1.20 or higher)
- Operating Systems: Windows, macOS, Linux

## Installation 🛠️

To install Git Review, run:

```bash
go install github.com/gifflet/git-review@latest
```

The binary will be placed in:
- **Linux/macOS**: `$HOME/go/bin/git-review`
- **Windows**: `%USERPROFILE%\go\bin\git-review.exe`

### Configure Git Review as a Global Command

After installation, you must configure `git-review` as a global Git command for easier usage:

**Linux/macOS**:
```bash
git config --global alias.review '!"$HOME/go/bin/git-review"'
```

**Windows**:
```powershell
git config --global alias.review "!%USERPROFILE%/go/bin/git-review.exe"
```

Once configured, you can run `git review` from anywhere within your repositories.

## Quick Start 🚀

```bash
# Basic usage
git review -i <initial_commit>
```

## Features ✨

- 📂 Compare two commits and get a comprehensive list of modified files
- 💾 Save diffs of modified files to a specified output directory
- 🌿 Flexible branch comparison with optional main branch filtering

## Usage 📚

### Basic Command Structure

```bash
git review -i <initial_commit> [-f <final_commit>] [options]
```

### Available Options

- `-i, --initial <commit>`: **Required**. Starting commit hash for comparison
- `-f, --final <commit>`: **Optional**. Ending commit hash (defaults to `HEAD`)
- `-m, --main-branch <branch_name>`: **Optional**. Main branch for refined comparisons
- `-p, --project-path <path>`: **Optional**. Project directory path (defaults to current directory)
- `-o, --output-dir <path>`: **Optional**. Output directory for diff files (defaults to `git-review`)

### Detailed Usage Scenarios

#### 1. Basic Commit Comparison
```powershell
git review -i c9286370
```

#### 2. Comparing Two Specific Commits
```powershell
git review -i c9286370 -f 78094299
```

#### 3. Comparing with Main Branch Context
```powershell
git review -i c9286370 -f 78094299 -m main
```

#### 4. Specifying Custom Project and Output Paths
```powershell
git review -i c9286370 -f 78094299 \
  -m main \
  -p "C:\Projects\MyRepo" \
  -o "./review-output"
```

## Troubleshooting 🛠

### Common Issues

- **Permission Denied**: Ensure you have read/write permissions in the project directory
- **Commit Not Found**: Verify commit hashes are correct and exist in your repository
- **Large Repository**: For repositories with extensive history, the process might take longer

## Performance Considerations ⚡

- Recommended for repositories up to 10GB in size
- Processing time increases with number of files and commit complexity
- Large binary files might slow down diff generation

## Contributing 🤝

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-improvement`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-improvement`)
5. Open a Pull Request

## License 📝

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Happy Reviewing! 🎉**
