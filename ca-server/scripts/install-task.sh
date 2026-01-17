#!/bin/bash

# install-task.sh - Install the latest Task CLI version using Go
# Usage: ./scripts/install-task.sh

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to get the latest version from GitHub API
get_latest_version() {
    local repo="go-task/task"
    local api_url="https://api.github.com/repos/${repo}/releases/latest"
    
    if command_exists curl; then
        curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget is available. Cannot fetch latest version."
        return 1
    fi
}

# Function to install Task using Go
install_task() {
    local version="$1"
    local package_url="github.com/go-task/task/v3/cmd/task@${version}"
    
    print_info "Installing Task CLI version ${version}..."
    
    # Install using go install
    if go install "$package_url"; then
        print_success "Task CLI ${version} installed successfully!"
        return 0
    else
        print_error "Failed to install Task CLI"
        return 1
    fi
}

# Function to verify installation
verify_installation() {
    if command_exists task; then
        local installed_version
        installed_version=$(task --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Task CLI is available: ${installed_version}"
        
        # Check if task is in PATH
        local task_path
        task_path=$(which task)
        print_info "Task CLI location: ${task_path}"
        
        return 0
    else
        print_error "Task CLI is not available in PATH"
        return 1
    fi
}

# Function to show Go binary path info
show_go_path_info() {
    if command_exists go; then
        local gopath
        local gobin
        gopath=$(go env GOPATH)
        gobin=$(go env GOBIN)
        
        print_info "Go environment:"
        print_info "  GOPATH: ${gopath}"
        print_info "  GOBIN: ${gobin:-${gopath}/bin}"
        
        # Check if Go bin directory is in PATH
        local go_bin_dir="${gobin:-${gopath}/bin}"
        if [[ ":$PATH:" == *":${go_bin_dir}:"* ]]; then
            print_success "Go binary directory is in PATH"
        else
            print_warning "Go binary directory is not in PATH"
            print_info "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            print_info "  export PATH=\"\$PATH:${go_bin_dir}\""
        fi
    fi
}

# Main function
main() {
    print_info "Task CLI Installation Script"
    print_info "============================="
    
    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed or not in PATH"
        print_error "Please install Go first: https://golang.org/dl/"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | awk '{print $3}')
    print_success "Go is available: ${go_version}"
    
    # Show Go path information
    show_go_path_info
    
    # Check if Task is already installed
    if command_exists task; then
        local current_version
        current_version=$(task --version 2>/dev/null | head -n1 || echo "unknown")
        print_warning "Task CLI is already installed: ${current_version}"
        
        # Ask if user wants to reinstall
        echo -n "Do you want to reinstall/update? (y/N): "
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled by user"
            exit 0
        fi
    fi
    
    # Get latest version
    print_info "Fetching latest Task CLI version..."
    local latest_version
    if ! latest_version=$(get_latest_version); then
        print_error "Failed to fetch latest version"
        print_info "Falling back to installing latest available version..."
        latest_version="latest"
    else
        print_success "Latest version found: ${latest_version}"
    fi
    
    # Install Task
    if install_task "$latest_version"; then
        print_success "Installation completed!"
        
        # Verify installation
        print_info "Verifying installation..."
        if verify_installation; then
            print_success "Task CLI is ready to use!"
            print_info ""
            print_info "Try running: task --version"
            print_info "Or in this project: task"
        else
            print_error "Installation verification failed"
            print_info "You may need to add Go's bin directory to your PATH"
            show_go_path_info
            exit 1
        fi
    else
        print_error "Installation failed"
        exit 1
    fi
}

# Run main function
main "$@"