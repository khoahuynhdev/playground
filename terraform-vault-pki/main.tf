terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 4.0"
    }
  }
}

provider "vault" {
  address = var.vault_address
  token   = var.vault_token
}

variable "vault_address" {
  description = "Vault server address"
  type        = string
  default     = "http://127.0.0.1:8200"
}

variable "vault_token" {
  description = "Vault authentication token"
  type        = string
  sensitive   = true
}

variable "pki_path" {
  description = "Path where PKI secrets engine will be mounted"
  type        = string
  default     = "pki"
}

variable "ca_common_name" {
  description = "Common name for the root CA certificate"
  type        = string
  default     = "example.com"
}

variable "ca_ttl" {
  description = "TTL for the root CA certificate"
  type        = string
  default     = "87600h" # 10 years
}