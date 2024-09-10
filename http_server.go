package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair("dev_server.crt", "dev_server.key")
	if err != nil {
		log.Fatalf("failed to load server certificate and key: %v", err)
	}

	// Create a TLS config with client certificate verification
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    loadCA("dev_client.crt"),
	}

	// Create a server with TLS config
	server := &http.Server{
		Addr:      ":8801",
		Handler:   http.HandlerFunc(handler),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server listening on port 8801")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// handler handles different HTTP request methods.
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodPut:
		handlePut(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled GET request")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled POST request")
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled PUT request")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handled DELETE request")
}

func loadCA(caCertPath string) *x509.CertPool {
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	return certPool
}
