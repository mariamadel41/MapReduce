package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func divideFile(fileData []byte) ([]byte, []byte) {
	// Calculate the length of the file data
	fileLength := len(fileData)

	// Divide the file data into two equal-sized chunks
	halfLength := fileLength / 2
	chunk1 := fileData[:halfLength]
	chunk2 := fileData[halfLength:]

	// Return the two chunks
	return chunk1, chunk2
}

func sendToSlave(addr string, chunk []byte) error {
	// Connect to the slave with a timeout of 5 seconds
	conn, err := net.DialTimeout("tcp", addr, 500*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the chunk to the slave
	_, err = conn.Write(chunk)
	if err != nil {
		return err
	}

	return nil
}

func fileHandler() error {
	// Open the sequence.txt file
	file, err := os.Open("sequence.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file data into memory
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Divide the file data into two chunks
	chunk1, chunk2 := divideFile(fileData)

	// Send chunk1 to Slave 1
	err = sendToSlave("192.168.1.137:8081", chunk1)
	if err != nil {
		return err
	}

	// Send chunk2 to Slave 2
	err = sendToSlave("192.168.1.109:8082", chunk2)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Handle file uploads and sending chunks to slaves
	err := fileHandler()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("File received and chunks sent to slaves successfully")
}
