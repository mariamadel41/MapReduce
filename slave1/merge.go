package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

func receiveChunk() error {
	// Listen for incoming connections on port 8081
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	defer ln.Close()

	// Accept incoming connection
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Read the chunk data from the connection
	chunk, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	// Save the received chunk to a file
	err = ioutil.WriteFile("slave1_chunk.txt", chunk, 0644)
	if err != nil {
		return err
	}

	return nil
}

func countCharacters(text string) int {
	return len(text)
}

func sendCountToMaster(count int) {
	masterAddr := "192.168.0.123:8000" // Replace with the actual master address

	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	countStr := strconv.Itoa(count)

	_, err = conn.Write([]byte(countStr))
	if err != nil {
		log.Fatalf("Failed to send character count to master: %v", err)
	}
}

func main() {
	// Receive the chunk from the master
	err := receiveChunk()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Chunk received and saved successfully on Slave 1")

	// Read the chunk from the file
	filename := "slave1_chunk.txt" // Replace with the actual filename

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	charCount := countCharacters(string(content))
	fmt.Printf("Total characters in file: %d\n", charCount)

	// Send the character count to the master
	sendCountToMaster(charCount)
}
