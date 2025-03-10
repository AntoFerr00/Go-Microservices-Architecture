package main

import (
	"Go_Microservices_Architecture/order-service"
	"Go_Microservices_Architecture/project-root/common"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/go_microservices_db"
	}
	common.ConnectToDatabase(dsn)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8082"
	}
	http.HandleFunc("/orders", order.GetOrdersHandler)

	log.Printf("Starting Order Service on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
