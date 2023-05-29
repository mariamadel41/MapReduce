package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

func sendToSlave(addr string, chunk []byte, wg *sync.WaitGroup) error {
	defer wg.Done()

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
	file, err := os.Open("data.txt")
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

func handleSlaveConnection(conn net.Conn, wg *sync.WaitGroup, results chan<- map[string]int) {
	defer conn.Close()

	// Read the word count map from the slave
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Failed to read word count from slave: %v", err)
	}

	// Parse the word count map received from the slave
	wordCount := make(map[string]int)
	countStr := string(buffer[:n])
	lines := strings.Split(countStr, "\n")
	for _, line := range lines {
		if line != "" {
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				word := strings.TrimSpace(parts[0])
				countStr := strings.TrimSpace(parts[1])
				if countStr != "" {
					count, err := strconv.Atoi(countStr)
					if err != nil {
						log.Fatalf("Failed to parse word count from slave: %v", err)
					}
					wordCount[word] = count
				}
			}
		}
	}

	// Send the word count map to the results channel
	results <- wordCount
}

func saveWordCountToFile(wordCounts map[string]int) error {
	content := ""
	for word, count := range wordCounts {
		content += fmt.Sprintf("%s,%d\n", word, count)
	}

	err := ioutil.WriteFile("reduce_result.txt", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func sendResultToClient(clientAddr string, wordCounts map[string]int) error {
	clientConn, err := net.Dial("tcp", clientAddr)
	if err != nil {
		return err
	}
	defer clientConn.Close()

	for word, count := range wordCounts {
		countStr := fmt.Sprintf("%s,%d\n", word, count)
		_, err = clientConn.Write([]byte(countStr))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Handle file uploads and sending chunks to slaves
	chunk1, chunk2, err := fileHandler()
	if err != nil {
		log.Fatalf("Failed to handle file: %v", err)
	}

	// Create a WaitGroup to synchronize sending chunks to slaves
	var wg sync.WaitGroup
	wg.Add(2)

	// Send chunk1 to Slave 1
	go sendToSlave("192.168.0.102:8081", chunk1, &wg)

	// Send chunk2 to Slave 2
	go sendToSlave("192.168.0.102:8082", chunk2, &wg)

	// Wait for the chunk sending to complete
	wg.Wait()

	// Start listening for slave connections
	slaveListener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to start slave listener: %v", err)
	}
	defer slaveListener.Close()

	// Accept connections from two slaves
	fmt.Println("Waiting for slaves to connect...")
	slaveCount := 0
	var wgSlaves sync.WaitGroup
	results := make(chan map[string]int, 2)
	for slaveCount < 2 {
		slave, err := slaveListener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept slave connection: %v", err)
		}
		fmt.Printf("Slave %d connected\n", slaveCount+1)

		// Handle the slave connection concurrently
		wgSlaves.Add(1)
		go handleSlaveConnection(slave, &wgSlaves, results)
		slaveCount++
	}
	// Wait for all slave connections to finish
	wgSlaves.Wait()

	// Close the results channel after all slaves have sent their word counts
	close(results)

	// Combine the word count maps from all the slaves
	totalCount := make(map[string]int)
	for countMap := range results {
		for word, count := range countMap {
			totalCount[word] += count
		}
	}

	fmt.Println("Word count from all slaves:")
	for word, count := range totalCount {
		fmt.Printf("%s, %d\n", word, count)
	}

	// Save the word count result to a file
	err = saveWordCountToFile(totalCount)
	if err != nil {
		log.Fatalf("Failed to save word count to file: %v", err)
	}

	// Send the word count to the client
	clientAddr := "192.168.0.103:9000" // Replace with the actual client address

	err = sendResultToClient(clientAddr, totalCount)
	if err != nil {
		log.Fatalf("Failed to send word count to client: %v", err)
	}

	fmt.Println("Word count sent to the client")
}
