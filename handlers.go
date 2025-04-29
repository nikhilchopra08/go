package main

import (
	"database/sql"
	"encoding/json"
	"go-api/models"
	"net/http"
	"log"
)

// SignupHandler handles user signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var signupData struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Decode the signup data from the request body
		err := json.NewDecoder(r.Body).Decode(&signupData)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Create a new user
		user := &models.User{
			Name:     signupData.Name,
			Email:    signupData.Email,
			Password: signupData.Password,
		}

		// Hash the password before storing
		err = user.HashPassword()
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Add the user to the database
		err = AddUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})
	} else {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}
}

// LoginHandler handles user login and JWT token generation
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Decode the login credentials from the request body
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Retrieve the user from the database
		var user models.User
		err = db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", loginData.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Println("Error querying user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the password is correct
		if !user.CheckPassword(loginData.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := GenerateJWT(&user)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Respond with the JWT token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	} else {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}
}
