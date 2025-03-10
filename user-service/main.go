package main

import (
	"Go_Microservices_Architecture/project-root/common"
	"Go_Microservices_Architecture/user-service"
	"fmt"
	"log"
	"net/http"
	"os"
	// assume authMiddleware is defined in this file or imported appropriately
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/go_microservices_db"
	}
	common.ConnectToDatabase(dsn)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	// Wrap the user handler with authentication middleware if needed.
	http.Handle("/users", authMiddleware(http.HandlerFunc(user.GetUsersHandler)))

	log.Printf("Starting User Service on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
