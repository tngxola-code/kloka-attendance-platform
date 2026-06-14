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

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jackc/pgx/v5/pgxpool"

    "kloka-attendance-platform/internal/domain/tenants"
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
    dbURL := flag.String("db", "", "PostgreSQL connection string (required)")
    platformKey := flag.String("platform-key", "", "Internal platform key for tenant creation")
    flag.Parse()

    if *dbURL == "" {
        slog.Error("DATABASE_URL is required")
        os.Exit(1)
    }

    ctx := context.Background()
    pool, err := pgxpool.New(ctx, *dbURL)
    if err != nil {
        slog.Error("failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer pool.Close()

    // Create tenants table if not exists
    if _, err := pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS tenants (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name TEXT NOT NULL,
            billing_email TEXT,
            country_code CHAR(2) NOT NULL DEFAULT 'ZA',
            status TEXT NOT NULL DEFAULT 'active',
            jwt_signing_secret TEXT NOT NULL,
            tenant_key_hash TEXT NOT NULL UNIQUE,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        );
        CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status);
    `); err != nil {
        slog.Error("failed to create tenants table", "error", err)
        os.Exit(1)
    }

    // Tenants wiring
    tenantRepo := tenants.NewRepository(pool)
    tenantSvc := tenants.NewService(tenantRepo)
    tenantHandler := tenants.NewHandler(tenantSvc, *platformKey)

    r := chi.NewRouter()
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(metrics.RequestMetrics)
    r.Use(versionHeaderMiddleware)

    // Root system endpoints
    r.Get("/health/live", system.LivenessHandler)
    r.Get("/health/ready", system.ReadinessHandler(nil))
    r.Get("/version", system.VersionHandler)
    r.Get("/openapi.yaml", system.OpenAPISpecHandler(openapiYAML))
    r.Get("/openapi.json", system.OpenAPIJSONHandler(openapiJSON))
    r.Handle("/metrics", metrics.MetricsHandler())

    // Versioned API
    r.Route(apiVersion, func(r chi.Router) {
        // System
        r.Get("/health/live", system.LivenessHandler)
        r.Get("/health/ready", system.ReadinessHandler(nil))
        r.Get("/version", system.VersionHandler)
        r.Get("/openapi.yaml", system.OpenAPISpecHandler(openapiYAML))
        r.Get("/openapi.json", system.OpenAPIJSONHandler(openapiJSON))
        r.Handle("/metrics", metrics.MetricsHandler())

        // Tenants
        r.Post("/tenants", tenantHandler.CreateTenant)
        r.Get("/tenants", tenantHandler.ListTenants)
        r.Get("/tenants/{id}", tenantHandler.GetTenant)
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
