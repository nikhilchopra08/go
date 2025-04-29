package main

import (
	"database/sql"
	"encoding/json"
	"go-api/models"
	"net/http"
	"strings"
)

// SignupHandler creates a new user
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if email is already taken
	var existingUser models.User
	err = db.QueryRow("SELECT id, name, email FROM users WHERE email = ?", user.Email).Scan(&existingUser.ID, &existingUser.Name, &existingUser.Email)
	if err != sql.ErrNoRows {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// LoginHandler authenticates a user and generates a JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	var user models.User
	err = db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", loginData.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Validate the password
	if user.Password != loginData.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT
	token, err := GenerateJWT(&user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetAllUsersHandler retrieves all registered users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Query the database for all users
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// If no users found, return empty list
	if len(users) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.User{})
		return
	}

	// Return the list of users as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	// Convert tags array to comma-separated string
	tagsStr := strings.Join(product.Tags, ",")

	// Insert into DB
	_, err = db.Exec(`
		INSERT INTO products 
		(id, name, description, image, price, discounted_price, rating, tags, is_new, is_bestseller)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID, product.Name, product.Description, product.Image,
		product.Price, product.DiscountedPrice, product.Rating, tagsStr,
		boolToInt(product.IsNew), boolToInt(product.IsBestseller),
	)

	if err != nil {
		http.Error(w, "Error saving product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product created successfully"})
}


func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, name, description, image, price, discounted_price, rating, tags, is_new, is_bestseller
		FROM products`)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var tagsStr string
		var isNew, isBestseller int

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.DiscountedPrice,
			&product.Rating,
			&tagsStr,
			&isNew,
			&isBestseller,
		); err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}

		// Parse JSON string of tags into a slice
		if err := json.Unmarshal([]byte(tagsStr), &product.Tags); err != nil {
			product.Tags = []string{}
		}

		product.IsNew = isNew == 1
		product.IsBestseller = isBestseller == 1

		products = append(products, product)
	}

	// Respond with an empty list if no products found
	if len(products) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Product{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
