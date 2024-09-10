package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "log"
    "net"
	"io/ioutil"
)

func main() {
    // Load server's certificate and private key
    serverCert, err := tls.LoadX509KeyPair("dev_server.crt", "dev_server.key")
    if err != nil {
        log.Fatalf("server: failed to load certificate: %v", err)
    }

    // Configure TLS with the certificate, CA pool, and verification options
    config := &tls.Config{
        Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
        ClientCAs:    loadCA("dev_client.crt"),
    }

    listener, err := tls.Listen("tcp", ":8800", config)
    if err != nil {
        log.Fatalf("server: failed to start listener: %v", err)
    }
    defer listener.Close()

    fmt.Println("Server listening on port 8800")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("server: failed to accept connection: %v", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Server: Connection established")

    // Read the data from the client
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        log.Printf("server: failed to read data: %v", err)
        return
    }

    // Extract the message
    clientMessage := string(buffer[:n])
    fmt.Printf("Server: Received message: %s\n", clientMessage)

    // Prepare the response message
    response := "server received " + clientMessage
    _, err = conn.Write([]byte(response))
    if err != nil {
        log.Printf("server: failed to write response: %v", err)
    }
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
