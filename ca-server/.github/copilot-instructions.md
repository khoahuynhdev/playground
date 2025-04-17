<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Cert Manager Server Project

This is a Go HTTP server project using the Gin framework. It follows a clean architecture with controllers, models, middleware, routes, and configuration.

## Coding Style Guidelines

- Use camelCase for variable names and PascalCase for exported functions
- Group imports: standard library first, then external packages, then local packages
- Include comments for exported functions and types
- Follow Go best practices for error handling
- Use dependency injection for services and repositories

## Project Structure

- `controllers/`: Handle HTTP requests and responses
- `models/`: Define data structures
- `middleware/`: Implement request processing middleware
- `routes/`: Define API routes
- `config/`: Application configuration
- `main.go`: Application entry point

## Common Tasks

When working with this codebase, follow these patterns:

- New route handlers should be added to appropriate controller
- New middleware should be registered in main.go
- New environment variables should be added to config.go
- Follow RESTful API best practices for endpoint design