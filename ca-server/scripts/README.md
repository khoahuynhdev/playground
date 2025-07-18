# Scripts

This directory contains utility scripts for the xpki project.

## install-task.sh

A bash script to install the latest Task CLI version using Go.

### Features

- **Automatic Version Detection**: Fetches the latest version from GitHub API
- **Go Integration**: Uses `go install` for installation
- **Error Handling**: Comprehensive error checking and validation
- **PATH Verification**: Checks if Go binary directory is in PATH
- **Interactive**: Prompts before reinstalling if Task is already present
- **Colored Output**: User-friendly colored terminal output
- **Fallback Support**: Works with both curl and wget

### Usage

```bash
# Make sure the script is executable (already done)
chmod +x scripts/install-task.sh

# Run the installation script
./scripts/install-task.sh
```

### Prerequisites

- Go must be installed and available in PATH
- Internet connection to fetch the latest version
- Either curl or wget for API requests

### What it does

1. **Validates Environment**: Checks if Go is installed
2. **Shows Go Info**: Displays GOPATH, GOBIN, and PATH information
3. **Checks Existing Installation**: Warns if Task is already installed
4. **Fetches Latest Version**: Gets the latest version from GitHub API
5. **Installs Task**: Uses `go install` to install the latest version
6. **Verifies Installation**: Confirms Task is properly installed and accessible
7. **Provides Guidance**: Shows PATH configuration if needed

### Example Output

```
[INFO] Task CLI Installation Script
=============================
[SUCCESS] Go is available: go1.22.5
[INFO] Go environment:
[INFO]   GOPATH: /home/user/go
[INFO]   GOBIN: /home/user/go/bin
[SUCCESS] Go binary directory is in PATH
[INFO] Fetching latest Task CLI version...
[SUCCESS] Latest version found: v3.35.1
[INFO] Installing Task CLI version v3.35.1...
[SUCCESS] Task CLI v3.35.1 installed successfully!
[INFO] Verifying installation...
[SUCCESS] Task CLI is available: Task version: v3.35.1
[INFO] Task CLI location: /home/user/go/bin/task
[SUCCESS] Task CLI is ready to use!
[INFO] 
[INFO] Try running: task --version
[INFO] Or in this project: task
```

### Troubleshooting

**Go not found:**
- Install Go from https://golang.org/dl/
- Make sure Go is in your PATH

**Task not in PATH after installation:**
- Add Go's bin directory to your PATH:
  ```bash
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```
- Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.)

**Network issues:**
- Check internet connection
- Verify firewall settings allow HTTPS requests to api.github.com
- The script will fall back to installing "latest" if version fetch fails