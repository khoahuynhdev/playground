package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"ca-server/config"
	"ca-server/middleware"
	"ca-server/models"
	"ca-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Set gin mode
	gin.SetMode(cfg.Mode)

	// Create gin router with default middleware
	r := gin.New()

	// Add custom middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Initialize store
	store := models.NewMemoryStore()

	// Setup routes
	routes.SetupRoutes(r, store)

	// WaitGroup to track active servers
	var wg sync.WaitGroup

	// Start HTTP server if enabled
	if !cfg.TLSEnabled && !cfg.MTLSEnabled {
		serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
		wg.Add(1)
		go startHTTPServer(r, serverAddr, &wg)
	}

	// Start TLS server
	if cfg.TLSEnabled {
		tlsAddr := fmt.Sprintf(":%d", cfg.TLSServerPort)
		wg.Add(1)
		go startTLSServer(r, tlsAddr, cfg.TLSCertPath, cfg.TLSKeyPath, &wg)
	}

	// Start mTLS server
	if cfg.MTLSEnabled {
		mtlsAddr := fmt.Sprintf(":%d", cfg.MTLSServerPort)
		wg.Add(1)
		go startMTLSServer(r, mtlsAddr, cfg.TLSCertPath, cfg.TLSKeyPath, cfg.ClientCACertPath, &wg)
	}

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Wait for servers to complete
	wg.Wait()
	log.Println("All servers shutdown complete")
}

// startHTTPServer starts a regular HTTP server
func startHTTPServer(handler http.Handler, addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Startup server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}
}

// startTLSServer starts a TLS server without client certificate validation
func startTLSServer(handler http.Handler, addr, certPath, keyPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Startup server in a goroutine
	go func() {
		log.Printf("Starting TLS server on %s", addr)
		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
			log.Fatalf("TLS server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down TLS server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("TLS server shutdown error: %v", err)
	}
}

// startMTLSServer starts a server with mutual TLS (client certificate validation)
func startMTLSServer(handler http.Handler, addr, certPath, keyPath, clientCACertPath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Load CA cert for client certificate validation
	caCert, err := ioutil.ReadFile(clientCACertPath)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate to pool")
	}

	// Configure TLS with client certificate verification
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	// Startup server in a goroutine
	go func() {
		log.Printf("Starting mTLS server on %s", addr)
		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
			log.Fatalf("mTLS server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down mTLS server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("mTLS server shutdown error: %v", err)
	}
}
