package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"crypto/x509"
	"io/ioutil"
)

// Args holds the arguments for our RPC method
type Args struct {
	A, B int
}

func main() {
	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair("dev_client.crt", "dev_client.key")
	if err != nil {
		log.Fatalf("failed to load client certificate and key: %v", err)
	}

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      loadCA("dev_server.crt"),
	}

	// Connect to the RPC server with TLS
	client, err := tls.Dial("tcp", "localhost:8802", tlsConfig)
	if err != nil {
		log.Fatalf("Error connecting to RPC server: %v", err)
	}
	defer client.Close()

	// Create an RPC client
	rpcClient := rpc.NewClient(client)

	args := Args{A: 5, B: 10}
	var result int

	// Call the Multiply method on the server
	err = rpcClient.Call("Arith.Multiply", args, &result)
	if err != nil {
		log.Fatalf("Error calling RPC method: %v", err)
	}

	fmt.Printf("Result: %d * %d = %d\n", args.A, args.B, result)
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