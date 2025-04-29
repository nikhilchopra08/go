package main

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go-api/models" // Import your models package for the User struct and methods
)

var db *sql.DB

// Initialize the database and create the users table if not already present
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./names.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Updated SQL schema with password field
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,  -- Ensures unique emails
		password TEXT NOT NULL       -- Store the hashed password
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

// AddUser adds a new user to the database after hashing the password
func AddUser(user *models.User) error {
	// Check if the email already exists
	var existingEmail string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
	if err != sql.ErrNoRows {
		if err == nil {
			// Email already exists
			return fmt.Errorf("email already in use")
		}
		log.Println("Error checking email:", err)
		return err
	}

	// Insert the new user (hashed password will be stored)
	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with that email
		}
		log.Println("Error querying user:", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a user by their unique ID
func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with that ID
		}
		log.Println("Error querying user:", err)
		return nil, err
	}
	return &user, nil
}
