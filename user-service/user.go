package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Go_Microservices_Architecture/project-root/common"
)

// User represents a user record.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetUsersHandler routes requests for the /users endpoint.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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

// getAllUsers retrieves all users from the database.
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

// getUserByID retrieves a specific user by ID.
func getUserByID(w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}
	var user User
	err = common.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, user)
}

// createUser inserts a new user into the database.
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

// updateUser updates an existing user record.
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

// deleteUser removes a user record from the database.
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
	w.WriteHeader(http.StatusNoContent)
}

// writeJSON writes the given data as JSON.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
