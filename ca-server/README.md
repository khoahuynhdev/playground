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

3. Create CA

```bash
curl -X POST http://localhost:8080/api/certs/ca
```

4. Create Server Cert

```bash
curl -X POST http://localhost:8080/api/certs/server
# copy the ca cert as ca-cert.pem and server credential as cert.pem and key.pem to the server/certs folder
```

5. Create client Cert

```bash
curl -X POST http://localhost:8080/api/certs/server
# copy the ca cert as ca-cert.pem and client credential as cert.pem and key.pem to the client/certs folder
```

6. Testing it

```bash
# in a terminal in server/
go run main.go

# in another terminal

# for tls
curl --cacert ca-cert.pem https://localhost:8080 -v
# for mtls
curl --cert cert.pem --key key.pem --cacert ca-cert.pem https://localhost:8443 -v
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

  - [x] validate signature

- [ ] Sign a CSR for new user
- [x] Example calling http with tls
- [x] Example calling http with mtls
- [ ] Write tests for tls and mtls case

### Notes

- quick way to generate CA certs from Go

```bash
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -ca -duration 87600h -host "localhost"
```

- Test TLS server curl

```bash
curl --cacert ./caCert.pem https://localhost:8080 -v
```

```bash
# Generate client private key
openssl ecparam -name prime256v1 -genkey -noout -out clientKey.pem

# Create client certificate signing request (CSR)
openssl req -new -key clientKey.pem -out clientReq.pem -subj "/CN=client"

# Sign the client certificate with CA
openssl x509 -req -in clientReq.pem -CA caCert.pem -CAkey caKey.pem -CAcreateserial -out clientCert.pem -days 365

curl --cert clientCert.pem --key clientKey.pem --cacert caCert.pem https://localhost:8443/api/resource
```
