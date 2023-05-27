package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
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
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
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

func fileHandler() ([]byte, []byte, error) {
	// Open the sequence.txt file
	file, err := os.Open("sequence.txt")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Read the file data into memory
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	// Divide the file data into two chunks
	chunk1, chunk2 := divideFile(fileData)

	return chunk1, chunk2, nil
}

func handleSlaveConnection(conn net.Conn, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	defer conn.Close()

	// Read the character count from the slave
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Failed to read character count from slave: %v", err)
	}

	// Parse the character count received from the slave
	count, err := strconv.Atoi(string(buffer[:n]))
	if err != nil {
		log.Fatalf("Failed to parse character count from slave: %v", err)
	}

	// Send the character count to the results channel
	results <- count
}

func main() {
	// Handle file uploads and sending chunks to slaves
	chunk1, chunk2, err := fileHandler()
	if err != nil {
		log.Fatalf("Failed to handle file: %v", err)
	}

	// Send chunk1 to Slave 1
	err = sendToSlave("192.168.0.113:8081", chunk1)
	if err != nil {
		log.Fatalf("Failed to send chunk1 to Slave 1: %v", err)
	}

	// Send chunk2 to Slave 2
	err = sendToSlave("192.168.0.110:8082", chunk2)
	if err != nil {
		log.Fatalf("Failed to send chunk2 to Slave 2: %v", err)
	}

	// Start listening for slave connections
	slaveListener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to start slave listener: %v", err)
	}
	defer slaveListener.Close()

	// Accept connections from two slaves
	fmt.Println("Waiting for slaves to connect...")
	slaveCount := 0
	var wg sync.WaitGroup
	results := make(chan int, 2)
	for slaveCount < 2 {
		slave, err := slaveListener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept slave connection: %v", err)
		}
		fmt.Printf("Slave %d connected\n", slaveCount+1)

		// Handle the slave connection concurrently
		wg.Add(1)
		go handleSlaveConnection(slave, &wg, results)
		slaveCount++
	}
	// Wait for all slave connections to finish
	wg.Wait()

	// Collect the results from the channel
	close(results)
	totalCount := 0
	for count := range results {
		totalCount += count
	}

	fmt.Printf("Total character count: %d\n", totalCount)

	// Send the result to the client
	clientAddr := "192.168.0.123:9000" // Replace with the actual client address

	clientConn, err := net.Dial("tcp", clientAddr)
	if err != nil {
		log.Fatalf("Failed to connect to client: %v", err)
	}
	defer clientConn.Close()

	resultStr := strconv.Itoa(totalCount)

	_, err = clientConn.Write([]byte(resultStr))
	if err != nil {
		log.Fatalf("Failed to send result to client: %v", err)
	}

	fmt.Println("Result sent to the client")
}
