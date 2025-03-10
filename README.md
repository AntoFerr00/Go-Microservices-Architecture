## Go Microservices Architecture

This project is a sample microservices architecture implemented in Go. It demonstrates how to build a system with multiple independent services that communicate via HTTP. The project includes:

- User Service: Provides CRUD operations for managing user data in a MySQL database with simple token-based authentication.
- Order Service: Manages orders and demonstrates service-to-service communication by fetching user details from the User Service.
- Common Package: Contains shared functionality (e.g., database connection logic).

# Features

- CRUD Operations: Create, read, update, and delete users.
- Service-to-Service Communication: Order Service retrieves user data from the User Service.
- Basic Authentication: Simple token-based protection on endpoints.
- MySQL Integration: Data persistence via a MySQL database.
- Docker Compose Setup: Run MySQL and both services together using Docker Compose.
- Environment Variables: Configure DSN and service ports dynamically.

# Prerequisites

- Go (v1.20 or later recommended)
- Docker (for running MySQL and services)
- MySQL database (automated via Docker Compose)

# Project Structure

Go_Microservices_Architecture/
├── docker-compose.yml
├── go.mod
├── go.sum
├── mysql-init/
│   └── init.sql         # SQL script to create tables and seed sample data
├── project-root/
│   └── common/
│       └── database.go  # Shared DB connection logic
├── order-service/
│   ├── Dockerfile
│   ├── main.go          # Service entrypoint; wires HTTP routing for orders
│   └── order.go         # Business logic for orders
└── user-service/
    ├── Dockerfile
    ├── main.go          # Service entrypoint; wires HTTP routing for users
    └── user.go          # Business logic for users

# Setup
1. Clone the Repository:

```bash
git clone <repository-url>
cd Go_Microservices_Architecture
```

2. Build and Run Services with Docker Compose:

The provided docker-compose.yml defines three services: MySQL, User Service, and Order Service. It also mounts the SQL initialization scripts.

- To build the images:

```bash
docker-compose build
```
- To start all services:

```bash
docker-compose up
```
Note:

- The MySQL container maps container port 3306 to host port 3308 (if port 3307 is in use, adjust as needed).
- Environment variables (e.g. MYSQL_DSN and SERVICE_PORT) are defined in the Compose file and used by the services.

3. Database Initialization (Optional):

The mysql-init/init.sql file automatically creates the required tables and seeds sample data when the MySQL container starts (if a new volume is used). The script includes:

```bash
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT,
  product VARCHAR(255),
  quantity INT,
  order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users (name) VALUES ('Alice'), ('Robert'), ('Charlie'), ('Diana');
INSERT INTO orders (user_id, product, quantity) VALUES (1, 'Laptop', 1), (2, 'Smartphone', 2);
```

# Running Services Manually

If you prefer to run the services without Docker Compose:

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
Make sure you have a running MySQL instance and that the environment variables (e.g., MYSQL_DSN) are configured appropriately.

# Usage
- User Service Endpoint:

Access at: http://localhost:8081/users

This endpoint is protected. Include the header:

```bash
Authorization: Bearer secret-token
```

Supported methods: GET (to list or get a specific user), POST, PUT, and DELETE.

- Order Service Endpoint:

Access at: http://localhost:8082/orders

The Order Service retrieves orders from the database and fetches user details from the User Service.

# Environment Variables
Both services use environment variables for configuration. For example, in the Docker Compose file:

- MYSQL_DSN: Connection string for MySQL (e.g., root:root@tcp(mysql_db:3306)/go_microservices_db).

- SERVICE_PORT: The port on which each service listens (8081 for User Service and 8082 for Order Service).

# License
This project is licensed under the MIT License. See the LICENSE file for details.