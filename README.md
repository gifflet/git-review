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

## Quick Start 🚀

```bash
# Clone the repository
git clone https://github.com/gifflet/git-review.git

# Build the binary
cd git-review
go build .

# Basic usage
git-review <initial_commit> [final_commit]
```

## Features ✨

- 📂 Compare two commits and get a comprehensive list of modified files
- 💾 Save diffs of modified files to a specified output directory
- 🌿 Flexible branch comparison with optional main branch filtering

## Installation 🛠️

### Method 1: Build from Source

```bash
git clone https://github.com/gifflet/git-review.git
cd git-review
go build .
```

### Method 2: Go Install

```bash
go install github.com/gifflet/git-review@latest
```

## Usage 📖

### Basic Command Structure

```bash
git-review <initial_commit> [final_commit] [options]
```

### Comprehensive Options

- `<initial_commit>`: **Required**. Starting commit hash for comparison
- `[final_commit]`: **Optional**. Ending commit hash (defaults to `HEAD`)
- `--version, -v`: Optional. Display the current version of git-review
- `--main-branch <branch_name>`: Optional main branch for refined comparisons
- `--project-path <path>`: Optional project directory path
- `--output-dir <path>`: Optional output directory for diff files

### Detailed Usage Scenarios

#### 1. Basic Commit Comparison
```powershell
git-review c9286370 78094299
```

#### 2. Comparing with Main Branch Context
```powershell
git-review c9286370 78094299 --main-branch main
```

#### 3. Specifying Custom Project and Output Paths
```powershell
git-review c9286370 78094299 \
  --main-branch main \
  --project-path "C:\Projects\MyRepo" \
  --output-dir "./review-output"
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


## License 📄

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Happy Reviewing! 🎉**