package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CertController struct{}

func (c *CertController) GetCert(ctx *gin.Context) {
	// Logic to get certificate
}

func (c *CertController) ListCerts(ctx *gin.Context) {
	// Logic to list certificates
}

func (c *CertController) CreateCert(ctx *gin.Context) {
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
		Headers: map[string]string{
			"CA-SERVER": "localhost",
		},
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
