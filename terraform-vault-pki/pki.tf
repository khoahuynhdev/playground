# Enable PKI secrets engine
resource "vault_mount" "pki" {
  path        = var.pki_path
  type        = "pki"
  description = "PKI secrets engine for certificate management"

  default_lease_ttl_seconds = 3600
  max_lease_ttl_seconds     = 315360000 # 10 years
}

# Generate root CA certificate
resource "vault_pki_secret_backend_root_cert" "root_ca" {
  backend     = vault_mount.pki.path
  type        = "internal"
  common_name = var.ca_common_name
  ttl         = var.ca_ttl
  format      = "pem"
  key_type    = "rsa"
  key_bits    = 4096

  exclude_cn_from_sans = true
}

# Configure CA and CRL URLs
resource "vault_pki_secret_backend_config_urls" "config_urls" {
  backend = vault_mount.pki.path
  issuing_certificates = [
    "${var.vault_address}/v1/${var.pki_path}/ca"
  ]
  crl_distribution_points = [
    "${var.vault_address}/v1/${var.pki_path}/crl"
  ]
}

# Create a role for issuing certificates
resource "vault_pki_secret_backend_role" "server_role" {
  backend          = vault_mount.pki.path
  name             = "server"
  ttl              = 3600
  max_ttl          = 86400
  allow_localhost  = true
  allow_ip_sans    = true
  allow_any_name   = true
  enforce_hostnames = false
  key_type         = "rsa"
  key_bits         = 2048
  key_usage = [
    "DigitalSignature",
    "KeyAgreement",
    "KeyEncipherment"
  ]
  ext_key_usage = [
    "ServerAuth"
  ]
}

# Example: Issue a certificate
resource "vault_pki_secret_backend_cert" "example_cert" {
  backend     = vault_mount.pki.path
  name        = vault_pki_secret_backend_role.server_role.name
  common_name = "test.${var.ca_common_name}"
  
  alt_names = [
    "localhost",
    "test-server.${var.ca_common_name}"
  ]
  
  ip_sans = [
    "127.0.0.1"
  ]
}