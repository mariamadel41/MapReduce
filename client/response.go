package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	clientListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to start client listener: %v", err)
	}
	defer clientListener.Close()

	fmt.Println("Waiting for master to connect...")
	masterConn, err := clientListener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept master connection: %v", err)
	}
	defer masterConn.Close()

	buffer := make([]byte, 1024)
	n, err := masterConn.Read(buffer)
	if err != nil {
		log.Fatalf("Failed to read result from master: %v", err)
	}

	result := string(buffer[:n])
	fmt.Printf("Received result from master: %s\n", result)
}
