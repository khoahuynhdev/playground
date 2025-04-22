package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getServerCerts() (tls.Certificate, error) {
	serverCert, _ := os.ReadFile("certs/cert.pem")
	serverKey, _ := os.ReadFile("certs/key.pem")
	return tls.X509KeyPair(serverCert, serverKey)
}

func getCaCertPool() *x509.CertPool {
	caCertPem, _ := os.ReadFile("certs/ca-cert.pem")
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertPem)
	return certPool
}

func StartMTLS() *http.Server {
	cert, err := getServerCerts()
	if err != nil {
		panic(err)
	}

	tlsConfig := tls.Config{
		ClientCAs:    getCaCertPool(),
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	// Example mTLS with Gin https://www.bastionxp.com/blog/api-gateway-security-mtls-authentication/
	server := &http.Server{
		Addr:      "localhost:8443",
		TLSConfig: &tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HTTPS with mTLS is working!")
		}),
	}

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			fmt.Printf("MTLS server error: %v\n", err)
		}
	}()

	fmt.Println("MTLS server started on localhost:8443")
	return server
}

func StartTLS() *http.Server {
	cert, err := getServerCerts()
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: "localhost:8080",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HTTPS with TLS is working!")
		}),
	}

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			fmt.Printf("TLS server error: %v\n", err)
		}
	}()

	fmt.Println("TLS server started on localhost:8080")
	return server
}

func main() {
	// Create context that listens for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())

	// Start both servers and keep instances for shutdown
	mtlsServer := StartMTLS()
	tlsServer := StartTLS()

	// Wait for termination signal
	sig := <-sigChan
	fmt.Printf("Received signal %v, shutting down servers...\n", sig)

	// Cancel context and implement graceful shutdown
	cancel()

	// Proper server shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := mtlsServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down MTLS server: %v\n", err)
	}
	if err := tlsServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down TLS server: %v\n", err)
	}

	fmt.Println("Servers stopped")
}
