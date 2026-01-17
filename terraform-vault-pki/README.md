# Terraform Vault PKI Playground

This playground demonstrates how to use Terraform to manage HashiCorp Vault's PKI (Public Key Infrastructure) secrets engine for certificate management.

## What's Included

This playground sets up:
- **PKI Secrets Engine**: Enables Vault's PKI capabilities
- **Root Certificate Authority**: Creates an internal root CA
- **Certificate Roles**: Defines policies for certificate issuance
- **Example Certificate**: Issues a sample server certificate
- **URL Configuration**: Sets up CA and CRL distribution points

## Prerequisites

1. **HashiCorp Vault**: Running Vault server (dev mode is fine for testing)
2. **Terraform**: Version 1.0 or later
3. **Vault CLI**: For initial setup and verification (optional)

## Getting Started

### 1. Start Vault Server

For development/testing, you can run Vault in dev mode:

```bash
vault server -dev
```

This starts Vault at `http://127.0.0.1:8200` with a root token displayed in the output.

### 2. Configure Variables

Copy the example variables file and update with your values:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:
```hcl
vault_address = "http://127.0.0.1:8200"
vault_token   = "your-root-token-here"
ca_common_name = "mycompany.com"
```

### 3. Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Plan the deployment
terraform plan

# Apply the configuration
terraform apply
```

### 4. Verify Setup

Check that the PKI engine is mounted:
```bash
vault secrets list
```

View the root CA certificate:
```bash
vault read pki/cert/ca
```

## What Gets Created

1. **PKI Mount**: Secrets engine at `/pki` path
2. **Root CA**: Self-signed root certificate (10-year TTL)
3. **Server Role**: Certificate role for server certificates
4. **Example Certificate**: Sample certificate for `test.example.com`

## Using the PKI

### Issue New Certificates

Once deployed, you can issue certificates using the `server` role:

```bash
vault write pki/issue/server common_name="app.example.com" ttl="24h"
```

### Revoke Certificates

```bash
vault write pki/revoke serial_number="39:dd:2e:90:b7:f8:95:00:2f:9b:37:f6:32:8d:2a:01:09:48:03:a5"
```

### View Certificate Revocation List

```bash
vault read pki/crl
```

## Configuration Options

### Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `vault_address` | Vault server URL | `http://127.0.0.1:8200` |
| `vault_token` | Vault authentication token | - |
| `pki_path` | PKI mount path | `pki` |
| `ca_common_name` | Root CA common name | `example.com` |
| `ca_ttl` | Root CA certificate TTL | `87600h` (10 years) |

### Outputs

The configuration provides these outputs:
- `root_ca_certificate`: Root CA certificate (PEM format)
- `example_certificate`: Example server certificate
- `example_private_key`: Example private key
- `pki_mount_path`: PKI secrets engine path
- `ca_issuing_url`: CA certificate download URL
- `crl_url`: Certificate revocation list URL

## Security Considerations

⚠️ **For Production Use:**

1. Use proper Vault authentication (not root tokens)
2. Configure appropriate certificate policies and roles
3. Set up proper CA certificate chain management
4. Implement certificate rotation strategies
5. Monitor certificate expiration
6. Secure private key storage

## Advanced Usage

### Multiple PKI Engines

You can create multiple PKI engines for different purposes:

```hcl
# Intermediate CA
resource "vault_mount" "intermediate_pki" {
  path = "intermediate-pki"
  type = "pki"
}
```

### Custom Certificate Roles

Create roles with specific constraints:

```hcl
resource "vault_pki_secret_backend_role" "client_role" {
  backend         = vault_mount.pki.path
  name           = "client"
  client_flag    = true
  server_flag    = false
  allow_any_name = false
  allowed_domains = ["client.example.com"]
}
```

## Cleanup

To destroy all resources:

```bash
terraform destroy
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Ensure Vault server is running
2. **Permission denied**: Verify your Vault token has sufficient privileges
3. **Path already in use**: Check if PKI engine is already mounted

### Useful Commands

```bash
# Check Vault status
vault status

# List all secrets engines
vault secrets list

# View PKI configuration
vault read pki/config/urls

# List certificate roles
vault list pki/roles
```

## Resources

- [Vault PKI Secrets Engine](https://www.vaultproject.io/docs/secrets/pki)
- [Terraform Vault Provider](https://registry.terraform.io/providers/hashicorp/vault/latest/docs)
- [PKI Best Practices](https://www.vaultproject.io/docs/secrets/pki/considerations)