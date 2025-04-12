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

- [] Gen a private/public key pairs for a user

  - [] validate signature

- [] Sign a CSR for new user

## Learning path

### Fundamentals

Generate RSA key pairs and save in PEM/DER formats
Generate ECDSA and Ed25519 key pairs
Basic encryption/decryption with different algorithms
Create and verify digital signatures

### Certificate Basics

Create self-signed certificates
Parse and display certificate details
Generate Certificate Signing Requests (CSRs)
Build a simple Certificate Authority
Implement certificate chain validation

### TLS Server Implementation

Create basic HTTPS server with Go
Implement mutual TLS (mTLS) authentication
Set up certificate pinning
Handle certificate revocation (CRL/OCSP)

### Certificate Management

Build API endpoints for certificate operations
Implement certificate rotation and renewal
Design secure certificate storage system
Create monitoring for certificate expiration

### Advanced Topics

Certificate Transparency integration
Automated certificate issuance with ACME
HSM integration for key protection
Implement OCSP stapling