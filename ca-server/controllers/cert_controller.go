package controllers

import (
	"ca-server/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type CertController struct{}

func (c *CertController) GetCert(ctx *gin.Context) {
	// Logic to get certificate
}

func (c *CertController) ListCerts(ctx *gin.Context) {
	// Logic to list certificates
}

func (c *CertController) SignCSR(ctx *gin.Context) {
}

// CreateCA create a certificate authority
// a CA should include a private key and a certificate (public key) which is self-signed
func (c *CertController) CreateCA(ctx *gin.Context) {
	caPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate CA key"})
		return
	}

	// generate a self-signed certificate
	caTmpl := &x509.Certificate{
		Subject:               pkix.Name{Country: []string{"VN"}, Organization: []string{"Private Homelab"}, CommonName: "My Homelab"},
		SerialNumber:          utils.NewSerialNum(), // random serial number
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10), // valid for 10 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	if caPriv == nil {
		ctx.JSON(500, gin.H{"error": "CA private key is nil"})
		return
	}
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, caPriv.Public(), caPriv)
	if err != nil {
	caPrivDER, err := x509.MarshalPKCS8PrivateKey(caPriv)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to marshal CA private key"})
		return
	}
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to marshal CA private key"})
		return
	}

	// Encode the certificate to PEM format for easy storage and transfer
		caCertPEM := pem.EncodeToMemory(&pem.Block{Bytes: caCertDER, Type: "CERTIFICATE"})
	// Encode the private key to PEM format to ensure compatibility with other tools and systems
		caPrivPEM := pem.EncodeToMemory(&pem.Block{Bytes: caPrivDER, Type: "EC Private Key"})
	ctx.JSON(200, gin.H{
		"certPEM": caCertPEM,
		"keyPEM":  caPrivPEM,
	}) // the client can then use base64 decode to view the PEM files
}

func (c *CertController) CreateKey(ctx *gin.Context) {
	// Logic to create a key
	// use RSA for now
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate key"})
		return
	}

	// convert to PEM format
	privatekeyDER := x509.MarshalPKCS1PrivateKey(privatekey)
	privatekeyBlock := pem.Block{
		Type: "RSA PRIVATE KEY",
		// WARN: why doesn't openssl RSA understand the headers?
		Bytes: privatekeyDER,
	}

	privakeyPEM := pem.EncodeToMemory(&privatekeyBlock)
	base64PEM := base64.StdEncoding.EncodeToString(privakeyPEM)

	fmt.Printf("Generated PEM successfully\n")
	ctx.JSON(200, gin.H{"pem": string(base64PEM), "base64_encoded": true})
}

func NewCertController() *CertController {
	return &CertController{}
}
