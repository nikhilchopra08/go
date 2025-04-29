package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go-api/models"
)

var db *sql.DB

// Initialize the database and create the users and products tables
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Users table
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	// Products table
	createProductsTable := `CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		image TEXT,
		price REAL,
		discounted_price REAL,
		rating REAL,
		tags TEXT, -- Stored as JSON array string
		is_new INTEGER,
		is_bestseller INTEGER
	);`

	_, err = db.Exec(createProductsTable)
	if err != nil {
		log.Fatal("Error creating products table:", err)
	}
}

// AddUser adds a new user to the database after checking for email uniqueness
func AddUser(user *models.User) error {
	var existingEmail string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email).Scan(&existingEmail)
	if err != sql.ErrNoRows {
		if err == nil {
			return fmt.Errorf("email already in use")
		}
		log.Println("Error checking email:", err)
		return err
	}

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}
	return nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println("Error fetching user:", err)
		return nil, err
	}
	return &user, nil
}

// AddProduct inserts a new product into the database
func AddProduct(product *models.Product) error {
	tagsJSON, err := json.Marshal(product.Tags)
	if err != nil {
		return fmt.Errorf("failed to encode tags: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO products 
		(id, name, description, image, price, discounted_price, rating, tags, is_new, is_bestseller)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.Name, product.Description, product.Image,
		product.Price, product.DiscountedPrice, product.Rating, string(tagsJSON),
		boolToInt(product.IsNew), boolToInt(product.IsBestseller),
	)
	if err != nil {
		log.Println("Error inserting product:", err)
		return err
	}
	return nil
}

// Helper: convert bool to int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// GetAllProducts retrieves all products from the database
func GetAllProducts() ([]models.Product, error) {
	rows, err := db.Query("SELECT id, name, description, image, price, discounted_price, rating, tags, is_new, is_bestseller FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		var tagsStr string
		var isNew, isBestseller int

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.DiscountedPrice, &p.Rating, &tagsStr, &isNew, &isBestseller)
		if err != nil {
			return nil, err
		}

		// Decode tags JSON array
		err = json.Unmarshal([]byte(tagsStr), &p.Tags)
		if err != nil {
			p.Tags = []string{} // fallback to empty
		}

		p.IsNew = isNew == 1
		p.IsBestseller = isBestseller == 1

		products = append(products, p)
	}

	return products, nil
}
