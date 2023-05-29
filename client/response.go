package main

import (
	"fmt"
	"io/ioutil"
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

	// Receive the word count result file from the master
	resultData, err := ioutil.ReadAll(masterConn)
	if err != nil {
		log.Fatalf("Failed to read result from master: %v", err)
	}

	// Save the word count result to a file
	err = ioutil.WriteFile("word_count_result.txt", resultData, 0644)
	if err != nil {
		log.Fatalf("Failed to save result to file: %v", err)
	}

	// Display the word count result on the terminal
	fmt.Println("Word count result:")
	fmt.Println(string(resultData))
}
