# Cert Manager HTTP Server

A scalable HTTP server template built with Go and the Gin framework.

## Project Structure

```
├── config/           # Application configuration
├── controllers/      # Request handlers
├── middleware/       # Custom middleware
├── models/           # Data models
├── routes/           # Route definitions
└── main.go           # Entry point
```

## Getting Started

### Prerequisites

- Go 1.16+

### Installation

1. Install dependencies:

```bash
go mod tidy
```

2. Run the server:

```bash
go run main.go
```

The server will start on port 8080 by default.

## Environment Variables

- `SERVER_PORT`: HTTP server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release) (default: debug)
- `LOG_LEVEL`: Logging level (default: info)

## API Endpoints

- `GET /`: Welcome message
- `GET /health`: Health check endpoint
- `GET /api/ping`: Ping endpoint

## Todo

- [x] Gen a private/public key pairs for a user

  - [] validate signature

- [] Sign a CSR for new user
- [x] Example calling http with tls

### Notes

- quick way to generate CA certs from Go

```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -ca -duration 87600h -host "localhost"
```
