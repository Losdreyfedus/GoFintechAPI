package api

import (
	"net/http"

	"backend_path/internal/api/handler"
	mw "backend_path/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Ortak middleware'ler
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Örnek endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is up!"))
	})

	// Auth route grubu
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
	})

	// User route grubu (korumalı)
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		r.Use(mw.RoleMiddleware("admin")) // örnek: sadece admin erişebilir
		r.Get("/", handler.ListUsers)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})

	// Transaction route grubu
	r.Route("/api/v1/transactions", func(r chi.Router) {
		r.Post("/credit", handler.Credit)
		r.Post("/debit", handler.Debit)
		r.Post("/transfer", handler.Transfer)
		r.Get("/history", handler.TransactionHistory)
		r.Get("/{id}", handler.GetTransaction)
	})

	// Balance route grubu
	r.Route("/api/v1/balances", func(r chi.Router) {
		r.Get("/current", handler.CurrentBalance)
		r.Get("/historical", handler.HistoricalBalance)
		r.Get("/at-time", handler.BalanceAtTime)
	})

	return r
}
