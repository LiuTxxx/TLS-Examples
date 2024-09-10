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

// Arith provides the arithmetic operations
type Arith struct{}

// Multiply performs multiplication of two integers
func (t *Arith) Multiply(args *Args, result *int) error {
	*result = args.A * args.B
	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)

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

	// Create a TLS listener
	listener, err := tls.Listen("tcp", ":8802", tlsConfig)
	if err != nil {
		log.Fatalf("Error starting TLS listener: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8802")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
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