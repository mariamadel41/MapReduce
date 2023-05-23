package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func sendFileToMaster(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new HTTP POST request to the master endpoint
	request, err := http.NewRequest("POST", "http://localhost:8080/upload", file)
	if err != nil {
		return err
	}

	// Set the appropriate content type for file uploads
	request.Header.Set("Content-Type", "multipart/form-data")

	// Send the request to the master
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send file to master. Status: %d", response.StatusCode)
	}

	fmt.Println("File sent to master successfully")

	return nil
}

func receiveResponseFromMaster() {
	// Send a GET request to the master endpoint to receive the response
	response, err := http.Get("http://localhost:8082/result")
	if err != nil {
		fmt.Println("Failed to receive response from master:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	// Process the response as needed
	fmt.Println("Response received from master:", string(body))
}

func main() {
	// Specify the file to send to the master
	filename := "path/to/your/file.txt"

	// Send the file to the master
	err := sendFileToMaster(filename)
	if err != nil {
		fmt.Println("Failed to send file to master:", err)
		return
	}

	// Receive the response from the master
	receiveResponseFromMaster()
}
