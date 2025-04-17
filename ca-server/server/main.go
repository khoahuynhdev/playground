package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

func getServerCerts() (tls.Certificate, error) {
	serverCert, _ := os.ReadFile("serverCert.pem")
	serverKey, _ := os.ReadFile("serverKey.pem")
	return tls.X509KeyPair(serverCert, serverKey)
}

func main() {
	cert, err := getServerCerts()
	if err != nil {
		panic(err)
	}

	serv := http.Server{
		Addr: "localhost:8443",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HTTPs server is working!")
		}),
	}

	serv.ListenAndServeTLS("", "")
}
