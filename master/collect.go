package main

import (
	"fmt"
	"net/http"
)

func slaveHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the response from a slave
	// You can process the response as needed
	fmt.Println("Response received from slave")
}

func main() {
	// Register the slaveHandler function to handle slave responses
	http.HandleFunc("/response", slaveHandler)

	// Start the server to listen for slave responses
	http.ListenAndServe(":8081", nil)
}
