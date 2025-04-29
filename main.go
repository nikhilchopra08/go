package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	initDB()

	// Set up HTTP routes
	http.HandleFunc("/signup", SignupHandler)  // Route for user signup
	http.HandleFunc("/login", LoginHandler)    // Route for user login

	// Start the server on port 8080
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
