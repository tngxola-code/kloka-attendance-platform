package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadyResponse struct {
	Status   string `json:"status"`
	Database bool   `json:"database"`
}

type VersionResponse struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "alive"})
}

func ReadinessHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbOk := false // default false when no DB or ping fails
		if db != nil {
			if err := db.PingContext(ctx); err == nil {
				dbOk = true
			}
		}

		status := http.StatusOK
		if !dbOk {
			status = http.StatusServiceUnavailable
		}

		respondJSON(w, status, ReadyResponse{
			Status:   map[bool]string{true: "ready", false: "not ready"}[dbOk],
			Database: dbOk,
		})
	}
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, VersionResponse{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	})
}

func OpenAPISpecHandler(specData []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write(specData)
	}
}

func OpenAPIJSONHandler(specData []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(specData)
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
