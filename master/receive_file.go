package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the uploaded file from the client request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server to save the uploaded file
	dst, err := os.Create(fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the contents of the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response to the client indicating the successful file upload
	fmt.Fprint(w, "File received and saved successfully")
}

func main() {
	// Register the fileHandler function to handle file uploads
	http.HandleFunc("/upload", fileHandler)

	// Start the server to listen for file uploads from the client
	http.ListenAndServe(":8080", nil)
}
