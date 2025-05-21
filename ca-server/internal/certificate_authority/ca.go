package ca

import (
	"crypto/x509"
	"fmt"
)

func CreatRootCA() {
}

func CreateNewCA(certConfig CertificateConfiguration) ([]string, *x509.Certificate, error) {
	checkInputError := false
	checkErorrs := make([]string, 0)
	var rootSlug string
	var caName string
	var rsaPrivateKeyPassword string

	if certConfig.Subject.CommonName == "" {
		checkInputError = true
		checkErorrs = append(checkErorrs, "missing common name field")
	} else {
		caName = certConfig.Subject.CommonName
		rootSlug = slugger(caName) // we do not really need to sligify it, just need for directory path building
	}

	if checkInputError {
		return checkErorrs, nil, fmt.Errorf("cert config error")
	}
	return []string{}, nil, nil
}
