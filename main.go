package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"kloka-attendance-platform/internal/auth"
	"kloka-attendance-platform/internal/domain/tenants"
	// "kloka-attendance-platform/internal/domain/workers"
	"kloka-attendance-platform/internal/metrics"
	"kloka-attendance-platform/internal/system"
)

//go:embed openapi.yaml
var openapiYAML []byte

//go:embed openapi.json
var openapiJSON []byte

const apiVersion = "/api/v1"

func main() {
	port := flag.String("port", "8080", "HTTP server port")
	flag.Parse()

	// In‑memory stores
	tenantRepo := tenants.NewMemoryRepository()
	tenantSvc := tenants.NewMemoryService(tenantRepo)
	refreshStore := tenants.NewMemoryRefreshStore()
	tenantHandler := tenants.NewHandler(tenantSvc, refreshStore)

	// Worker parts commented out for now
	// workerRepo := workers.NewMemoryRepository()
	// workerSvc := workers.NewService(workerRepo)
	// workerHandler := workers.NewHandler(workerSvc)

	// Secret lookup for JWT middleware
	secretLookup := func(ctx context.Context, tenantID string) ([]byte, error) {
		return tenantSvc.SigningSecret(ctx, tenantID)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(metrics.RequestMetrics)
	r.Use(versionHeaderMiddleware)

	// Root system endpoints (no auth)
	r.Get("/health/live", system.LivenessHandler)
	r.Get("/health/ready", system.ReadinessHandler(nil))
	r.Get("/version", system.VersionHandler)
	r.Get("/openapi.yaml", system.OpenAPISpecHandler(openapiYAML))
	r.Get("/openapi.json", system.OpenAPIJSONHandler(openapiJSON))
	r.Handle("/metrics", metrics.MetricsHandler())

	// Versioned API
	r.Route(apiVersion, func(r chi.Router) {
		// System (versioned)
		r.Get("/health/live", system.LivenessHandler)
		r.Get("/health/ready", system.ReadinessHandler(nil))
		r.Get("/version", system.VersionHandler)
		r.Get("/openapi.yaml", system.OpenAPISpecHandler(openapiYAML))
		r.Get("/openapi.json", system.OpenAPIJSONHandler(openapiJSON))
		r.Handle("/metrics", metrics.MetricsHandler())

		// Token endpoints (public)
		r.Post("/tenants/token", tenantHandler.Token)
		r.Post("/tenants/refresh", tenantHandler.Refresh)

		// Tenant provisioning (platform-key protected)
		r.Post("/tenants", tenantHandler.CreateTenant)

		// Authenticated tenant self-service
		r.Group(func(r chi.Router) {
			r.Use(auth.Authenticate(secretLookup))

			r.With(auth.RequireScope("tenant:read"), auth.RequireTenantMatch).
				Get("/tenants/{id}", tenantHandler.GetTenant)
		})

		// Worker endpoints commented out for now
		// r.Group(func(r chi.Router) {
		//     r.Use(auth.Authenticate(secretLookup))
		//     r.Use(auth.RequireScope("workers:read"))
		//     r.Get("/workers", workerHandler.ListWorkers)
		//     r.Post("/workers", workerHandler.CreateWorker)
		//     r.Get("/workers/{id}", workerHandler.GetWorker)
		//     r.Patch("/workers/{id}", workerHandler.UpdateWorker)
		//     r.Delete("/workers/{id}", workerHandler.DeleteWorker)
		// })

		// Worker login (commented out)
		// r.Post("/auth/worker/login", workerHandler.LoginWorker)
		// r.Post("/auth/worker/refresh", workerHandler.RefreshWorkerToken)
	})

	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("starting server", "port", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

func versionHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", "v1")
		next.ServeHTTP(w, r)
	})
}
