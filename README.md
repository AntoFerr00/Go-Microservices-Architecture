# Go Microservices Architecture

This project is a sample microservices architecture implemented in Go. It demonstrates how to build a system with multiple independent services that communicate via HTTP. The project includes:

- **User Service**: Provides CRUD operations for managing user data in a MySQL database with simple token-based authentication.
- **Order Service**: Manages orders and demonstrates service-to-service communication by fetching user details from the User Service.
- **Common Package**: Contains shared functionality (e.g., database connection logic).

## Features

- **CRUD Operations**: Create, read, update, and delete users.
- **Service-to-Service Communication**: Order Service retrieves user data from User Service.
- **Basic Authentication**: Simple token-based protection on endpoints.
- **MySQL Integration**: Data persistence via a MySQL database (can be run locally using Docker).

## Prerequisites

- [Go](https://golang.org/dl/) (v1.20 or later recommended)
- [Docker](https://www.docker.com/get-started) (for running MySQL locally)
- MySQL database

## Setup

1. **Clone the Repository:**
   ```bash
   git clone <repository-url>
   cd Go_Microservices_Architecture
   ```
2. **Set Up the MySQL Database:**

- Run MySQL using Docker (example mapping container port 3306 to host port 3307):
   ```bash
    docker run --name local-mysql -e
   MYSQL_ALLOW_EMPTY_PASSWORD=yes -p 3307:3306 -d mysql:8.0
   ``` 
- Connect to MySQL and create the database:
sql
    ```bash
    CREATE DATABASE go_microservices_db;
    ```
- Create the required tables:
    ```bash
    USE go_microservices_db;

    CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    );

    CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    product VARCHAR(255),
    quantity INT,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    ```
3. **Populate the Database (Optional):**

    ```bash
    INSERT INTO users (name) VALUES ('Alice'), ('Robert'), ('Charlie'), ('Diana');
    INSERT INTO orders (user_id, product, quantity) VALUES (1, 'Laptop', 1), (2, 'Smartphone', 2);
    ```
4. **Run the Services:**

- User Service:
    ```bash
    cd user-service
    go run .
    ```
- Order Service:
    ```bash
    cd order-service
    go run .
    ```
## Usage
- User Service Endpoint:
Access at http://localhost:8081/users
Note: This endpoint is protected. Include the header:

    ```bash
    Authorization: Bearer secret-token
    ```
- Order Service Endpoint:
Access at http://localhost:8082/orders

## License
This project is licensed under the MIT License. See the LICENSE file for details.