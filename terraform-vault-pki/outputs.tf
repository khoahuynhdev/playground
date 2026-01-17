output "root_ca_certificate" {
  description = "Root CA certificate"
  value       = vault_pki_secret_backend_root_cert.root_ca.certificate
  sensitive   = true
}

output "example_certificate" {
  description = "Example server certificate"
  value       = vault_pki_secret_backend_cert.example_cert.certificate
  sensitive   = true
}

output "example_private_key" {
  description = "Example server private key"
  value       = vault_pki_secret_backend_cert.example_cert.private_key
  sensitive   = true
}

output "pki_mount_path" {
  description = "PKI secrets engine mount path"
  value       = vault_mount.pki.path
}

output "ca_issuing_url" {
  description = "CA certificate issuing URL"
  value       = "${var.vault_address}/v1/${var.pki_path}/ca"
}

output "crl_url" {
  description = "Certificate Revocation List URL"
  value       = "${var.vault_address}/v1/${var.pki_path}/crl"
}