package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// Import the shared DB package
	"Go_Microservices_Architecture/project-root/common"
)

// Simple struct representing a user in our system
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Middleware that checks for a simple Bearer token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// Expect "Bearer secret-token" for this example
		if token != "Bearer secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Connect to DB (change port or credentials as needed)
	common.ConnectToDatabase("root:@tcp(127.0.0.1:3307)/go_microservices_db")

	// Wrap our /users endpoint with authentication middleware
	http.Handle("/users", authMiddleware(http.HandlerFunc(userHandler)))

	log.Println("Starting User Service on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// userHandler routes GET, POST, PUT, DELETE methods to different functions
func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			getAllUsers(w, r)
		} else {
			getUserByID(w, r, idParam)
		}
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		idParam := r.URL.Query().Get("id")
		updateUser(w, r, idParam)
	case http.MethodDelete:
		idParam := r.URL.Query().Get("id")
		deleteUser(w, r, idParam)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// =============== Handlers for each CRUD operation ===============

// GET all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := common.DB.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	writeJSON(w, users)
}

// GET a user by ID
func getUserByID(w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var user User
	err = common.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name)
	if err != nil {
		http.Error(w, "User not found or error querying database: "+err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, user)
}

// POST - create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := common.DB.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
	if err != nil {
		http.Error(w, "Error inserting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = int(lastID)

	writeJSON(w, user)
}

// PUT - update an existing user by ID
func updateUser(w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = common.DB.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, id)
	if err != nil {
		http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id
	writeJSON(w, user)
}

// DELETE - remove a user by ID
func deleteUser(w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = common.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// Utility to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
