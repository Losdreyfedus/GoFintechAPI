# Go Fintech Backend

**Go Fintech Backend** is a scalable, modular backend application designed for modern financial transactions and digital banking operations. Built with Go, it supports microservice-ready architecture, real-time balance tracking, event sourcing, and production-grade observability through Prometheus and Grafana.

---

## 🔑 Key Features

- **User Management** – Registration, login, and role-based authorization
- **Transaction Processing** – Credit, debit, and transfer operations
- **Real-Time Balance Tracking** – Current and historical balance queries
- **Event Sourcing & Scheduling** – Reliable banking operations
- **Redis Caching** – Performance boost with in-memory caching
- **Monitoring & Metrics** – Prometheus endpoints and Grafana dashboards
- **Dockerized** – Easy local and production deployment

---

## ⚙️ Tech Stack

- **Go 1.24+**
- **Chi Router** – Lightweight HTTP routing
- **Microsoft SQL Server** – Primary relational database
- **Redis** – Caching and session management
- **Prometheus** – Metrics collection and alerting
- **Grafana** – Metrics visualization and dashboards
- **Docker & Docker Compose** – Containerization and orchestration

---

## 🚀 Getting Started

### 📦 Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose
- Git

### ⚙️ Setup Steps

```bash
# Clone the repository
git clone <your-repo-url>
cd GoProject

# Download Go dependencies
go mod download

# Build the project
go build ./cmd/main.go

# Run all services with Docker
docker-compose up --build
```

---

### 🔌 Services
| Service      | Address                      |
|--------------|-----------------------------|
| API          | http://localhost:8080       |
| Prometheus   | http://localhost:9090       |
| Grafana      | http://localhost:3000 (admin/admin) |
| SQL Server   | localhost:1433              |
| Redis        | localhost:6379              |

---

## 📡 API Endpoints

### 🔐 Authentication
- `POST /api/v1/auth/register` – Register a new user
- `POST /api/v1/auth/login` – Authenticate user
- `POST /api/v1/auth/refresh` – Refresh JWT token

### 👤 User Management (Admin Only)
- `GET /api/v1/users` – List all users
- `GET /api/v1/users/{id}` – Get user details
- `PUT /api/v1/users/{id}` – Update user
- `DELETE /api/v1/users/{id}` – Delete user

### 💳 Transactions
- `POST /api/v1/transactions/credit` – Add funds
- `POST /api/v1/transactions/debit` – Withdraw funds
- `POST /api/v1/transactions/transfer` – Transfer funds
- `GET /api/v1/transactions/history` – View transaction history
- `GET /api/v1/transactions/{id}` – Transaction details

### 💰 Balance
- `GET /api/v1/balances/current` – Get current balance
- `GET /api/v1/balances/historical` – View past balances
- `GET /api/v1/balances/at-time` – Balance at a specific timestamp

### 📈 Monitoring
- `GET /metrics` – Prometheus metrics endpoint

---

## 📁 Project Structure
```
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── api/                 # HTTP API layer
│   │   ├── handler/         # Route handlers
│   │   ├── middleware/      # Auth, logging, etc.
│   │   └── router.go        # API route definitions
│   ├── balance/             # Balance logic
│   ├── config/              # Environment config
│   ├── domain/              # Domain models
│   ├── metrics/             # Prometheus instrumentation
│   ├── process/             # Business logic layer
│   └── user/                # User operations
├── migrations/              # Database schema migrations
├── pkg/
│   ├── database/            # SQL Server setup
│   └── redis/               # Redis client
├── Dockerfile               # Application Docker image
├── docker-compose.yml       # Multi-service orchestration
└── prometheus.yml           # Prometheus config
```

---

## 🧪 Testing
```bash
go test ./...
```

## 🧹 Linting
```bash
golangci-lint run
```

---

## 🤝 Contributing
Contributions are welcome! Feel free to open issues or submit pull requests. For major changes, please open an issue first to discuss what you’d like to change.
