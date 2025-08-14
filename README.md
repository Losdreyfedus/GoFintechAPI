# Go Fintech Backend

**Go Fintech Backend** is a scalable, modular backend application designed for modern financial transactions and digital banking operations. Built with Go, it supports microservice-ready architecture, real-time balance tracking, event sourcing, and production-grade observability through Prometheus and Grafana.

---

## 🔑 Key Features

- **User Management** – Registration, login, and role-based authorization
- **Transaction Processing** – Credit, debit, and transfer operations with full business logic
- **Real-Time Balance Tracking** – Current and historical balance queries with thread-safe operations
- **Event Sourcing & Scheduling** – Reliable banking operations
- **Redis Caching** – Performance boost with in-memory caching
- **Rate Limiting** – IP-based rate limiting for API protection
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
| Jaeger       | http://localhost:16686      |
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

## 🧱 Database Schema

### `users` Table
| Column      | Type           | Description             |
|-------------|----------------|-------------------------|
| id          | INT (PK)       | Auto-incremented ID     |
| username    | NVARCHAR(100)  | Unique user identifier  |
| password    | NVARCHAR(100)  | Hashed user password    |
| created_at  | DATETIME       | Account creation date   |

Other core tables include: `transactions`, `balances`, `audit_logs`.

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

## 🚀 Recent Updates

### ✅ Completed Core Features (Latest)
- **TransactionService** - Full implementation of credit, debit, and transfer operations
- **BalanceService** - Thread-safe balance management with historical tracking
- **API Handlers** - Complete transaction and balance endpoint implementations
- **Rate Limiting** - Configurable IP-based rate limiting middleware
- **Repository Layer** - SQL implementations for transactions and balances
- **DTO Updates** - Added missing request/response structures

### 🔧 Technical Improvements
- Enhanced error handling with structured logging and standardized error responses
- Thread-safe balance operations using sync.RWMutex
- Proper validation for all transaction types using go-playground/validator
- Comprehensive API response formatting with trace IDs
- Service dependency injection in handlers
- **OpenTelemetry Integration** - Distributed tracing with Jaeger
- **Advanced Security Headers** - HSTS, CSP, XSS Protection, etc.

### 🚨 Critical Fixes Applied (Latest)
- **Auth Middleware** - Fixed JWT token validation and user context injection
- **Database Transactions** - Added SQL Server transaction atomicity (BEGIN TX/COMMIT/ROLLBACK)
- **User CRUD Handlers** - Completed all user management endpoints
- **SQL Server Compatibility** - Fixed parameter placeholders (@p1, @p2, etc.)
- **Database Indexes** - Added performance indexes for critical queries
- **Role-based Authorization** - Enhanced middleware for admin-only routes

### 🔒 Security Improvements
- Proper JWT token validation in all protected routes
- User context properly injected from JWT claims
- All transaction and balance endpoints now require authentication
- Role-based access control for admin operations
- **Security Headers** - Comprehensive security middleware
- **Input Validation** - Structured validation with custom rules
- **Rate Limiting** - Configurable per-minute request limits

### 🚀 New Enterprise Features
- **Distributed Tracing** - OpenTelemetry + Jaeger integration
- **Standardized Error Handling** - Trace ID, error codes, structured responses
- **Advanced Validation** - Custom validation rules for business logic
- **Security Hardening** - Production-ready security headers
- **Observability** - Enhanced monitoring and debugging capabilities

---

## 🤝 Contributing
Contributions are welcome! Feel free to open issues or submit pull requests. For major changes, please open an issue first to discuss what you’d like to change.
