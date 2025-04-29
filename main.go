package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	initDB()

	// Define the HTTP routes
	http.HandleFunc("/signup", SignupHandler)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
