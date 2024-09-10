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
	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair("dev_client.crt", "dev_client.key")
	if err != nil {
		log.Fatalf("failed to load client certificate and key: %v", err)
	}

	// Create a TLS config with client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      loadCA("dev_server.crt"),
	}

	// Create an HTTP client with TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make a request to the server
	resp, err := client.Get("https://localhost:8801")
	if err != nil {
		log.Fatalf("client failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Print the server response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}
	fmt.Printf("Server response: %s\n", body)
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
