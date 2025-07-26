# OpenBao Playground

A simple setup to run OpenBao in Docker for learning and experimentation.

## What is OpenBao?

OpenBao is an open-source secrets management platform that helps you securely store and access secrets like API keys, passwords, certificates, and more. It's a fork of HashiCorp Vault focused on community governance and open development.

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- curl (for health checks)

### Running OpenBao

1. **Start OpenBao**:
   ```bash
   ./run-openbao.sh
   ```

2. **Access the Web UI**:
   Open http://localhost:8200 in your browser

3. **Login**:
   - Token: `myroot`

### Alternative: Manual Docker Commands

```bash
# Start with docker-compose
docker-compose up -d

# Stop
docker-compose down

# View logs
docker-compose logs -f
```

## Basic Usage

### CLI Setup

Install the OpenBao CLI:
```bash
# Download from https://github.com/openbao/openbao/releases
# Or use the container:
docker exec -it openbao-dev openbao
```

### Environment Variables

```bash
export OPENBAO_ADDR='http://localhost:8200'
export OPENBAO_TOKEN='myroot'
```

### Basic Commands

```bash
# Check status
openbao status

# Enable key-value secrets engine
openbao secrets enable -path=secret kv-v2

# Store a secret
openbao kv put secret/myapp db_password="supersecret"

# Retrieve a secret
openbao kv get secret/myapp

# List secrets
openbao kv list secret/
```

## Learning Resources

- [OpenBao Documentation](https://openbao.org/docs/)
- [API Reference](https://openbao.org/api-docs/)
- [GitHub Repository](https://github.com/openbao/openbao)

## Development Notes

- This setup uses **development mode** - data is stored in memory and lost on restart
- For production, configure persistent storage and proper authentication
- The root token `myroot` is for development only

## Troubleshooting

### OpenBao won't start
```bash
# Check if port 8200 is in use
lsof -i :8200

# View container logs
docker-compose logs openbao
```

### Permission issues
```bash
# Ensure the script is executable
chmod +x run-openbao.sh
```

## Stopping OpenBao

```bash
docker-compose down
```

To remove all data:
```bash
docker-compose down -v
```