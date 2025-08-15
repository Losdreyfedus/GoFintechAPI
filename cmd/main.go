package main

import (
	"backend_path/internal/api"
	"backend_path/internal/api/handler"
	"backend_path/internal/auth"
	"backend_path/internal/balance"
	"backend_path/internal/config"
	"backend_path/internal/transaction"
	"backend_path/internal/user"
	"backend_path/pkg/cache"
	"backend_path/pkg/database"
	"backend_path/pkg/jwt"
	"backend_path/pkg/logger"
	"backend_path/pkg/server"
	"backend_path/pkg/tracing"
	"backend_path/pkg/validator"
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// .env dosyasını otomatik yükle
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.InitLogger("info", cfg.Environment == "development")

	// Initialize validator
	validator.InitValidator()

	// Initialize RBAC manager (for future use)
	_ = auth.NewRBACManager()

	// Initialize advanced cache (for future use)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_ = cache.NewAdvancedCache(redisClient, cache.CacheStrategy(cfg.CacheStrategy))

	// Initialize tracing
	if err := tracing.InitTracer("gofintech-backend", "1.0.0", cfg.JaegerURL); err != nil {
		logger.Error("Failed to initialize tracing", err, map[string]interface{}{
			"jaeger_url": cfg.JaegerURL,
		})
		// Don't fail startup if tracing fails
	}

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
	balanceService := balance.NewService(balanceRepo)
	transactionService := transaction.NewService(transactionRepo, balanceService)

	// Set service dependencies in handlers
	handler.SetUserService(userService)
	handler.SetTransactionService(transactionService)
	handler.SetBalanceService(balanceService)

	// Create router with dependencies
	router := api.NewRouter(userService, jwtService, cfg)

	// Create server
	srv := server.NewServer(":"+cfg.Port, router)

	// Start server with graceful shutdown
	logger.Info("Starting application", map[string]interface{}{
		"port":        cfg.Port,
		"environment": cfg.Environment,
		"jaeger_url":  cfg.JaegerURL,
		"rate_limit":  cfg.RateLimit,
	})

	// Graceful shutdown with tracing cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed to start", err, nil)
	}

	// Cleanup tracing
	if err := tracing.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown tracing", err, nil)
	}
}
