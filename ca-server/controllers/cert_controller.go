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
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultValidDays = 365

type CertController struct{}

func (c *CertController) GetCert(ctx *gin.Context) {
	// Logic to get certificate
}

func (c *CertController) ListCerts(ctx *gin.Context) {
	// Logic to list certificates
}

func (c *CertController) SignCSR(ctx *gin.Context) {
}

func (c *CertController) CreateServerCert(ctx *gin.Context) {
	// Parse request body
	var req struct {
		CommonName  string   `json:"commonName"`
		DNSNames    []string `json:"dnsNames"`
		IPAddresses []string `json:"ipAddresses,omitempty"`
		ValidDays   int      `json:"validDays"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid request: " + err.Error(),
			"expected": gin.H{
				"commonName":  "string (required)",
				"dnsNames":    "array of strings (required)",
				"ipAddresses": "array of strings (optional)",
				"validDays":   "integer (required)",
			},
		})
		return
	}

	// Read CA key from file
	caKeyData, err := os.ReadFile("caKey.pem")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to read CA key: " + err.Error()})
		return
	}

	// Parse CA key
	block, _ := pem.Decode(caKeyData)
	if block == nil || block.Type != "EC Private Key" || block.Bytes == nil {
		ctx.JSON(500, gin.H{"error": "Failed to decode CA key PEM or invalid block bytes"})
		return
	}

	caPriv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to parse CA private key: " + err.Error()})
		return
	}

	// Read CA cert from file
	caCertData, err := os.ReadFile("caCert.pem")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to read CA certificate: " + err.Error()})
		return
	}

	// Parse CA certificate
	block, _ = pem.Decode(caCertData)
	if block == nil || block.Type != "CERTIFICATE" || block.Bytes == nil {
		ctx.JSON(500, gin.H{"error": "Failed to decode CA certificate PEM or invalid block bytes"})
		return
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to parse CA certificate: " + err.Error()})
		return
	}

	// Generate server private key with configurable curve
	curve := elliptic.P256() // Default curve
	if curveName := ctx.Query("curve"); curveName != "" {
		switch curveName {
		case "P384":
			curve = elliptic.P384()
		case "P521":
			curve = elliptic.P521()
		case "P256": // Explicitly handle P256
			curve = elliptic.P256()
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported curve: " + curveName})
			return
		}
	}

	serverPriv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate server key: " + err.Error()})
		return
	}

	// Parse IP addresses if provided
	var ips []net.IP
	for _, ip := range req.IPAddresses {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("Invalid IP address: %s", ip)})
			return
		}
		ips = append(ips, parsedIP)
	}

	// Create certificate template
	validDays := defaultValidDays
	if req.ValidDays > 0 {
		validDays = req.ValidDays
	}

	certTemplate := &x509.Certificate{
		Subject:      pkix.Name{CommonName: req.CommonName},
		SerialNumber: utils.NewSerialNum(),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * time.Duration(validDays)),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     req.DNSNames,
		IPAddresses:  ips,
	}

	// Sign the certificate
	certDER, err := x509.CreateCertificate(
		rand.Reader,
		certTemplate,
		caCert,
		serverPriv.Public(),
		caPriv,
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create certificate: " + err.Error()})
		return
	}

	// Encode server private key to PEM
	serverPrivDER, err := x509.MarshalPKCS8PrivateKey(serverPriv)
	if err != nil {
		if err == x509.IncorrectPasswordError || err == x509.ErrUnsupportedAlgorithm {
			ctx.JSON(500, gin.H{"error": "Unsupported key type for server private key: " + err.Error()})
		} else {
			ctx.JSON(500, gin.H{"error": "Failed to marshal server private key: " + err.Error()})
		}
		return
	}

	// Create PEM blocks
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverPrivDER})

	ctx.JSON(200, gin.H{
		"certPEM": string(certPEM),
		"keyPEM":  string(keyPEM),
	})
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
		ctx.JSON(500, gin.H{"error": "Failed to marshal CA cert"})
		return
	}
	caPrivDER, err := x509.MarshalPKCS8PrivateKey(caPriv)
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
