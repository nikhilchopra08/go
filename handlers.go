package main

import (
	"database/sql"
	"encoding/json"
	"go-api/models"
	"net/http"
)

// SignupHandler handles the signup functionality
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var user models.User

		// Decode JSON body into User struct
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Hash the password before storing it
		err = user.HashPassword()
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Check if the email is already taken
		var existingEmail string
		err = db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
		if err != sql.ErrNoRows {
			http.Error(w, "Email already taken", http.StatusConflict)
			return
		}

		// Insert the new user into the database
		_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
		if err != nil {
			http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
			return
		}

		// Respond with the newly created user (without password)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		})
	} else {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}
}
