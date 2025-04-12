package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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
		Headers: map[string]string{
			"CA-SERVER": "localhost",
		},
		Bytes: privatekeyDER,
	}

	privakeyPEM := pem.EncodeToMemory(&privatekeyBlock)

	fmt.Printf("Generated PEM successfully\n")
	ctx.JSON(200, gin.H{"pem": string(privakeyPEM)})
}

func NewCertController() *CertController {
	return &CertController{}
}
