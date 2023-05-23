package main

import (
	"fmt"
	"net/http"
)

func resultHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the result received from slaves
	// You can process or forward the result to the client as needed
	fmt.Println("Result received from slaves")
}

func main() {
	// Register the resultHandler function to handle the result
	http.HandleFunc("/result", resultHandler)

	// Start the server to listen for the result
	http.ListenAndServe(":8082", nil)
}
