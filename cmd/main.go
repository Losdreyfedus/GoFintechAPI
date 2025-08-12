package main

import (
	"backend_path/internal/api"
	"backend_path/internal/api/handler"
	"backend_path/internal/balance"
	"backend_path/internal/config"
	"backend_path/internal/transaction"
	"backend_path/internal/user"
	"backend_path/pkg/database"
	"backend_path/pkg/jwt"
	"backend_path/pkg/logger"
	"backend_path/pkg/server"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını otomatik yükle
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.InitLogger("info", cfg.Environment == "development")

	// Initialize database connection
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", err, nil)
	}
	defer db.Close()

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWTSecret, 1*time.Hour)

	// Initialize repositories
	userRepo := user.NewSQLRepository(db.DB)
	transactionRepo := transaction.NewSQLRepository(db.DB)
	balanceRepo := balance.NewSQLRepository(db.DB)

	// Initialize services
	userService := user.NewService(userRepo)
	transactionService := transaction.NewService(transactionRepo)
	balanceService := balance.NewService(balanceRepo)

	// Set service dependencies in handlers
	handler.SetTransactionService(transactionService)
	handler.SetBalanceService(balanceService)

	// Create router with dependencies
	router := api.NewRouter(userService, jwtService)

	// Create server
	srv := server.NewServer(":"+cfg.Port, router)

	// Start server with graceful shutdown
	logger.Info("Starting application", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
	})

	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed to start", err, nil)
	}
}
