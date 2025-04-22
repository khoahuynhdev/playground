package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// Load CA cert for server verification
	caCertPem, err := os.ReadFile("certs/ca-cert.pem")
	if err != nil {
		fmt.Printf("Error loading CA certificate: %v\n", err)
		return
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCertPem) {
		fmt.Println("Failed to append CA certificate to pool")
		return
	}

	// Load client certificate and key for mTLS
	clientCert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		fmt.Printf("Error loading client cert/key: %v\n", err)
		return
	}

	// Create HTTP client with mTLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      certPool,
				Certificates: []tls.Certificate{clientCert},
			},
		},
	}

	// Make request to server
	resp, err := client.Get("https://localhost:8443")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Printf("Response: %s\n", resp.Status)
		fmt.Printf("Body: %s\n", string(body))
	}
}
