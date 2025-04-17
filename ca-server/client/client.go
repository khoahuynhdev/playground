package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

func main() {
	caCertPem, _ := os.ReadFile("caCert.pem")
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertPem)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
		},
	}

	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		fmt.Errorf("Error making request: %v", err)
	} else {
		fmt.Printf("Response: %s", resp.Status)
	}
}
