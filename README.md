# Golang Inventory System

A robust, hexagonal architecture-based backend of inventory management system built with Go. This application provides a RESTful API for managing users, items, and borrowing records.

## 🚀 Tech Stack

### Backend
-   **Language:** [Go 1.25+](https://go.dev/)
-   **Framework:** [Fiber v2](https://gofiber.io/)
-   **Database:** [PostgreSQL](https://www.postgresql.org/)
-   **ORM:** [GORM](https://gorm.io/)
-   **Caching:** [Redis](https://redis.io/)
-   **Authentication:** JWT
-   **Testing:** Testify, GoMock
-   **Containerization:** Docker & Docker Compose

### Frontend
-   **Framework:** [React 19](https://react.dev/)
-   **Build Tool:** [Vite](https://vitejs.dev/)
-   **Language:** [TypeScript](https://www.typescriptlang.org/)
-   **Styling:** [Tailwind CSS v4](https://tailwindcss.com/)
-   **State Management:** [TanStack Query](https://tanstack.com/query/latest)
-   **Routing:** [React Router](https://reactrouter.com/)
-   **UI Libraries:** [Lucide React](https://lucide.dev/), SweetAlert2, React Hot Toast

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

2.  **Install Backend Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Install Frontend Dependencies**
    ```bash
    cd frontend
    npm install
    ```

4.  **Environment Setup**
    - **Backend:** Create a `.env` file in the root directory. You can copy `.env.example`.
    - **Frontend:** The frontend interacts with the backend at `http://localhost:8080` by default.

### Running the Application

**1. Start Backend Services (DB & Redis)**
```bash
make docker-run
# OR
docker-compose up -d
```

**2. Run Backend Server**

Standard run:
```bash
make run
```

With hot reload (requires [Air](https://github.com/air-verse/air)):
```bash
make dev
```

**3. Run Frontend Client**

Open a new terminal:
```bash
cd frontend
npm run dev
```

The application will be available at:
- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`

### Testing

Run unit tests:
```bash
make test
```

Run integration tests:
```bash
make test-integration
```