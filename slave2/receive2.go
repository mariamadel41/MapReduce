package main

import (
	"fmt"
	"io/ioutil"
	"net"
	
)

func receiveChunk() error {
	// Listen for incoming connections on port 8081
	ln, err := net.Listen("tcp", ":8082")
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

func main() {
	// Receive the chunk from the master
	err := receiveChunk()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Chunk received and saved successfully on Slave 1")
}