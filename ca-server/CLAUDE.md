# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is "xpki" (eXtensible PKI) - a Go-based Certificate Authority (CA) management system with both CLI and HTTP server components. The application provides PKI infrastructure with support for TLS and mTLS configurations.

## Architecture

The project follows a clean architecture pattern with:
- **CLI Interface**: Cobra-based command structure with `server` and `setup` subcommands
- **HTTP Server**: Gin framework with middleware, controllers, and routes
- **PKI Core**: Certificate Authority functionality in `internal/certificate_authority/`
- **Multi-server Support**: HTTP, TLS, and mTLS server modes

### Key Components

- `cmd/`: Cobra command definitions (server, setup)
- `internal/certificate_authority/`: Core CA functionality and types
- `controllers/`: HTTP request handlers (cert_controller, user_controller)
- `models/`: Data structures and in-memory store
- `middleware/`: Custom middleware (auth, logger)
- `routes/`: API route definitions
- `config/`: Application configuration management

## Development Commands

This project uses [Task](https://taskfile.dev) as a task runner. Install it first:

**Quick Installation (using provided script):**
```bash
# Run the provided installation script
./scripts/install-task.sh
```

**Manual Installation:**
```bash
# Install Task (various methods available)
go install github.com/go-task/task/v3/cmd/task@latest
# or: brew install go-task/tap/go-task
# or: curl -sL https://taskfile.dev/install.sh | sh
```

### Common Tasks
```bash
# Show all available tasks
task

# Development workflow
task workflow:dev              # Complete dev workflow (clean, deps, lint, test, build)
task workflow:setup            # Initial project setup

# Build tasks
task build                     # Build the application
task build:debug              # Build with debug symbols
task build:release            # Build optimized release binary

# Development
task dev                       # Run in development mode
task dev:setup                 # Setup PKI directory structure

# Testing
task test                      # Run all tests
task test:coverage            # Run tests with coverage
task test:race                # Run tests with race detection

# Code quality
task lint                      # Run linters
task lint:fix                 # Fix linting issues

# Dependencies
task deps                      # Download dependencies
task deps:update              # Update dependencies
```

### Server Tasks
```bash
# Start servers
task server                    # Start HTTP server
task server:dev               # Start server in development mode
task server:tls               # Start server with TLS
task server:mtls              # Start server with mTLS

# Test connectivity
task test:http                # Test HTTP server
task test:tls                 # Test TLS server
task test:mtls                # Test mTLS server
```

### PKI Tasks
```bash
# PKI management
task pki:setup                # Setup PKI directory structure
task pki:clean                # Clean PKI directory

# Certificate generation
task cert:ca                  # Generate CA certificate
task cert:server              # Generate server certificate
task cert:client              # Generate client certificate

# Clean up
task clean                    # Clean build artifacts
task clean:all                # Clean everything including PKI and certs
```

### Legacy Commands (without Task)
```bash
# Install dependencies
go mod tidy

# Build the application
go build -o xpki

# Run CLI commands
go run main.go setup              # Setup PKI directory structure
go run main.go server             # Start the HTTP/TLS/mTLS server
```

## Configuration

The application uses environment variables for configuration:
- `SERVER_PORT`: HTTP server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)
- `LOG_LEVEL`: Logging level

## Code Style (from .github/copilot-instructions.md)

- Use camelCase for variables, PascalCase for exported functions
- Group imports: standard library, external packages, local packages
- Include comments for exported functions and types
- Follow Go best practices for error handling
- Use dependency injection for services and repositories
- New route handlers go in appropriate controllers
- New middleware registered in main.go
- New environment variables added to config.go

## Server Modes

The application supports three server modes:
1. **HTTP**: Standard HTTP server (port 8080)
2. **TLS**: HTTPS server with server certificates
3. **mTLS**: Mutual TLS with client certificate validation

## Current Development Status

Based on recent commits, the project is in active development with:
- CLI setup command implementation
- PKI directory structure creation
- CA certificate generation (work in progress)
- Server configuration for multiple TLS modes