package main

import (
	"fmt"
	"net/http"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the file received from the client
	// You can process or save the file as needed
	fmt.Println("File received from client")
}

func main() {
	// Register the fileHandler function to handle file uploads
	http.HandleFunc("/upload", fileHandler)

	// Start the server to listen for file uploads from the client
	http.ListenAndServe(":8080", nil)
}
