package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Build-with-Go/go-starter/internal/config"
	"github.com/Build-with-Go/go-starter/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	logger *logger.Logger
	router *chi.Mux
	server *http.Server
}

// New creates a new HTTP server
func New(cfg *config.Config, log *logger.Logger) *Server {
	s := &Server{
		config: cfg,
		logger: log,
		router: chi.NewRouter(),
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures server middleware
func (s *Server) setupMiddleware() {
	// Core middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	// Logging middleware
	s.router.Use(s.loggingMiddleware())

	// Timeout middleware
	s.router.Use(middleware.Timeout(time.Duration(s.config.Server.ReadTimeout) * time.Second))

	// CORS middleware (basic setup)
	s.router.Use(middleware.AllowContentType("application/json"))
	s.router.Use(middleware.AllowContentType("text/plain"))
}

// loggingMiddleware creates a custom logging middleware
func (s *Server) loggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			next.ServeHTTP(ww, r)

			// Log request
			duration := time.Since(start).Milliseconds()
			s.logger.HTTPRequest(
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				ww.Status(),
				duration,
			)
		})
	}
}

// setupRoutes configures server routes
func (s *Server) setupRoutes() {
	// Health check routes
	s.router.Get("/healthz", s.handleHealthz)
	s.router.Get("/ready", s.handleReady)

	// API routes (placeholder for future expansion)
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", s.handleAPIRoot)
	})

	// Catch-all route
	s.router.NotFound(s.handleNotFound)
	s.router.MethodNotAllowed(s.handleMethodNotAllowed)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         s.config.GetAddr(),
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}

	s.logger.Startup("go-starter", "1.0.0", s.config.GetAddr(), s.config.Logger.Level)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Shutdown("Graceful shutdown initiated")

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.GracefulShutdown("http server", err)
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.GracefulShutdown("http server", nil)
	return nil
}

// handleHealthz handles basic health check
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"go-starter"}`))
}

// handleReady handles readiness probe
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual readiness checks (database connectivity, etc.)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","service":"go-starter"}`))
}

// handleAPIRoot handles API root endpoint
func (s *Server) handleAPIRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"go-starter API v1.0.0","endpoints":["/healthz","/ready"]}`))
}

// handleNotFound handles 404 errors
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"not found","path":"` + r.URL.Path + `"}`))
}

// handleMethodNotAllowed handles 405 errors
func (s *Server) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error":"method not allowed","method":"` + r.Method + `"}`))
}
