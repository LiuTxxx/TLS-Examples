package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "log"
	"io/ioutil"
)

func main() {
    // Load client's certificate and private key
    clientCert, err := tls.LoadX509KeyPair("dev_client.crt", "dev_client.key")
    if err != nil {
        log.Fatalf("client: failed to load certificate: %v", err)
    }

    // Configure TLS with the certificate, CA pool, and verification options
    config := &tls.Config{
        Certificates: []tls.Certificate{clientCert},
        RootCAs:      loadCA("dev_server.crt"),
    }

    conn, err := tls.Dial("tcp", "localhost:8800", config)
    if err != nil {
        log.Fatalf("client: failed to connect: %v", err)
    }
    defer conn.Close()

    fmt.Println("Client: Connected to server")

    // Send some data to the server
    _, err = conn.Write([]byte("Hello 123!"))
    if err != nil {
        log.Fatalf("client: failed to write data: %v", err)
    }

    // Read response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
    if err != nil {
        log.Printf("client: failed to read response: %v", err)
        return
    }

    // Extract the message
    serverResponse := string(buffer[:n])
    fmt.Printf("Client: Received response: %s\n", serverResponse)
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

