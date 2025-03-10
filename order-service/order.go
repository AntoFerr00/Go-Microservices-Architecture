package order

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"Go_Microservices_Architecture/project-root/common"
)

// Order represents an order record.
type Order struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Product   string `json:"product"`
	Quantity  int    `json:"quantity"`
	OrderDate string `json:"order_date"`
}

// OrderWithUser combines Order with the user's name.
type OrderWithUser struct {
	Order
	UserName string `json:"user_name"`
}

// User is used to decode responses from the user-service.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetOrdersHandler handles GET requests to /orders.
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllOrders(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllOrders retrieves orders from the database and augments them with user data.
func getAllOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := common.DB.Query("SELECT id, user_id, product, quantity, order_date FROM orders")
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Product, &o.Quantity, &o.OrderDate); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	var results []OrderWithUser
	for _, o := range orders {
		user, err := getUserFromUserService(o.UserID)
		userName := ""
		if err == nil && user != nil {
			userName = user.Name
		}
		results = append(results, OrderWithUser{
			Order:    o,
			UserName: userName,
		})
	}

	writeJSON(w, results)
}

// getUserFromUserService makes a service-to-service call to get the user data.
func getUserFromUserService(userID int) (*User, error) {
	url := fmt.Sprintf("http://localhost:8081/users?id=%d", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer secret-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respBody string
		if bodyBytes, err := io.ReadAll(resp.Body); err == nil {
			respBody = string(bodyBytes)
		}
		log.Printf("User service returned status %d: %s", resp.StatusCode, respBody)
		return nil, fmt.Errorf("user service returned status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// writeJSON writes the data as JSON.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
