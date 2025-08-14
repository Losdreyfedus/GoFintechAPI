package api

import (
	"net/http"
	"time"

	"backend_path/internal/api/handler"
	mw "backend_path/internal/api/middleware"
	"backend_path/internal/user"
	"backend_path/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(userService user.UserService, jwtService *jwt.JWTService) http.Handler {
	r := chi.NewRouter()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, jwtService)

	// Set service dependencies for handlers
	handler.SetUserService(userService)
	// Note: SetTransactionService should be called from main.go when transactionService is available

	// Ortak middleware'ler
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.ErrorHandlingMiddleware)                 // Error handling and panic recovery
	r.Use(mw.PrometheusMiddleware)                    // Prometheus metrics
	r.Use(mw.RateLimitMiddleware(100, 1*time.Minute)) // Rate limiting: 100 requests per minute
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Metrics endpoint for Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Örnek endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is up!"))
	})

	// Auth route grubu
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
	})

	// User route grubu (korumalı)
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(mw.AuthMiddleware(jwtService))
		r.Use(mw.RoleMiddleware("admin")) // örnek: sadece admin erişebilir
		r.Get("/", handler.ListUsers)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})

	// Transaction route grubu (korumalı)
	r.Route("/api/v1/transactions", func(r chi.Router) {
		r.Use(mw.AuthMiddleware(jwtService))
		r.Post("/credit", handler.Credit)
		r.Post("/debit", handler.Debit)
		r.Post("/transfer", handler.Transfer)
		r.Get("/history", handler.TransactionHistory)
		r.Get("/{id}", handler.GetTransaction)
	})

	// Balance route grubu (korumalı)
	r.Route("/api/v1/balances", func(r chi.Router) {
		r.Use(mw.AuthMiddleware(jwtService))
		r.Get("/current", handler.CurrentBalance)
		r.Get("/historical", handler.HistoricalBalance)
		r.Get("/at-time", handler.BalanceAtTime)
	})

	return r
}
