# Golang Inventory System

A robust, hexagonal architecture-based backend of inventory management system built with Go. This application provides a RESTful API for managing users, items, and borrowing records.

## 🚀 Tech Stack

-   **Language:** [Go 1.25+](https://go.dev/)
-   **Framework:** [Fiber v2](https://gofiber.io/)
-   **Database:** [PostgreSQL](https://www.postgresql.org/)
-   **ORM:** [GORM](https://gorm.io/)
-   **Caching:** [Redis](https://redis.io/)
-   **Authentication:** JWT
-   **Testing:** Testify, GoMock
-   **Containerization:** Docker & Docker Compose

## 🏗 Architecture

This project follows the **Hexagonal Architecture** (also known as Ports and Adapters), which enforces a clear separation of concerns:

-   **`internal/core`**: Contains the business logic and domain entities (`User`, `Item`, `Borrowing`). This layer is independent of external frameworks or databases.
-   **`internal/adapter`**: implemenets the interfaces defined by the core to interact with the outside world (e.g., Database repositories, HTTP handlers).
-   **`cmd`**: Entry point of the application.

## ✨ Features

-   **User Management**: Registration, Login, and Profile management.
-   **Item Management**: CRUD operations for inventory items.
-   **Borrowing System**: Track item lending and returns.
-   **Rate Limiting**: Redis-based rate limiting for API endpoints.

## 🛠 Getting Started

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/nitikhon/golang-inventory-system.git
    cd golang-inventory-system
    ```

2.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Environment Setup**
    Ensure you have a `.env` file in the root directory. You can copy the example `.env.example` file.

### Running the Application

**Using Docker (Recommended for dependencies)**

Start the database and redis services:
```bash
make docker-run
# OR
docker-compose up -d
```

**Run the Server**

Standard run:
```bash
make run
```

If you're using WSL, you can use the following command to run the server in development mode with hot reload (requires [Air](https://github.com/air-verse/air)):
```bash
make dev
```

### Testing

Run unit tests:
```bash
make test
```

Run integration tests:
```bash
make test-integration
```