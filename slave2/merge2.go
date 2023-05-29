package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
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
	err = ioutil.WriteFile("slave2_chunk.txt", chunk, 0644)
	if err != nil {
		return err
	}

	return nil
}

func countCharacters(text string) int {
	return len(text)
}

func countWords(text string) map[string]int {
	words := strings.Fields(text)
	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	return wordCount
}

func saveWordCount(wordCounts map[string]int) error {
	content := ""
	for word, count := range wordCounts {
		content += fmt.Sprintf("%s,%d\n", word, count)
	}

	err := ioutil.WriteFile("word2_count.txt", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func sendCountToMaster() error {
	masterAddr := "192.168.0.103:8000" // Replace with the actual master address

	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Read the word count file
	content, err := ioutil.ReadFile("word2_count.txt")
	if err != nil {
		return err
	}

	// Send the word count file to the master
	_, err = conn.Write(content)
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

	fmt.Println("Chunk received and saved successfully on Slave")

	// Read the chunk from the file
	filename := "slave2_chunk.txt" // Replace with the actual filename

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	wordCounts := countWords(string(content))

	// Save the word count to a file
	err = saveWordCount(wordCounts)
	if err != nil {
		log.Fatalf("Failed to save word count: %v", err)
	}

	// Send the word count file to the master
	err = sendCountToMaster()
	if err != nil {
		log.Fatalf("Failed to send word count to master: %v", err)
	}
}
