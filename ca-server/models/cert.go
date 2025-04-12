package models

import (
	"crypto/x509"
	"encoding/pem"
	"time"
)

// Certificate represents an X.509 certificate
type Certificate struct {
	SerialNumber    string            `json:"serialNumber"`
	Subject         string            `json:"subject"`
	Issuer          string            `json:"issuer"`
	NotBefore       time.Time         `json:"notBefore"`
	NotAfter        time.Time         `json:"notAfter"`
	IsCA            bool              `json:"isCA"`
	SignatureAlg    string            `json:"signatureAlg"`
	PublicKeyAlg    string            `json:"publicKeyAlg"`
	RawCertificate  []byte            `json:"-"`
	PemEncodedCert  string            `json:"pemEncodedCert,omitempty"`
	X509Certificate *x509.Certificate `json:"-"`
}

// ParseCertificate creates a Certificate model from PEM encoded certificate data
func ParseCertificate(pemData []byte) (*Certificate, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, x509.CertificateInvalidError{}
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &Certificate{
		SerialNumber:    cert.SerialNumber.String(),
		Subject:         cert.Subject.CommonName,
		Issuer:          cert.Issuer.CommonName,
		NotBefore:       cert.NotBefore,
		NotAfter:        cert.NotAfter,
		IsCA:            cert.IsCA,
		SignatureAlg:    cert.SignatureAlgorithm.String(),
		PublicKeyAlg:    cert.PublicKeyAlgorithm.String(),
		RawCertificate:  block.Bytes,
		PemEncodedCert:  string(pemData),
		X509Certificate: cert,
	}, nil
}
