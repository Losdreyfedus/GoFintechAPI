package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend_path/pkg/logger"

	"github.com/rs/zerolog"
)

// Server represents an HTTP server with graceful shutdown
type Server struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

// NewServer creates a new HTTP server instance
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger.GetLogger(),
	}
}

// Start starts the server with graceful shutdown handling
func (s *Server) Start() error {
	// Create a channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests in a separate goroutine.
	go func() {
		s.logger.Info().Str("addr", s.httpServer.Addr).Msg("Server started")
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Create a channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a signal or an error.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error starting server: %w", err)

	case sig := <-shutdown:
		s.logger.Info().Str("signal", sig.String()).Msg("Start shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown the server.
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error().Err(err).Str("timeout", "30s").Msg("Could not stop server gracefully")
			if err := s.httpServer.Close(); err != nil {
				return fmt.Errorf("could not force close server: %w", err)
			}
		}
	}

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
