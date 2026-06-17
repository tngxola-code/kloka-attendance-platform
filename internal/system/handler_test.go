package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLivenessHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	LivenessHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != "alive" {
		t.Errorf("expected status 'alive', got %q", resp.Status)
	}
}

func TestReadinessHandler_NoDB(t *testing.T) {
	// No database provided → returns not ready with 503 and database: false
	handler := ReadinessHandler(nil)
	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}

	var resp ReadyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Status != "not ready" {
		t.Errorf("expected status 'not ready', got %q", resp.Status)
	}
	if resp.Database != false {
		t.Error("expected database false when DB is nil")
	}
}

func TestVersionHandler(t *testing.T) {
	// Override version variables for test
	originalVersion := Version
	Version = "test-version"
	defer func() { Version = originalVersion }()

	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	VersionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp VersionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Version != "test-version" {
		t.Errorf("expected version 'test-version', got %q", resp.Version)
	}
}

func TestOpenAPISpecHandler(t *testing.T) {
	dummySpec := []byte("openapi: 3.0.3\ninfo:\n  title: Test")
	handler := OpenAPISpecHandler(dummySpec)
	req := httptest.NewRequest("GET", "/openapi.yaml", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/x-yaml" {
		t.Errorf("expected Content-Type application/x-yaml, got %q", ct)
	}
	if body := w.Body.String(); body != string(dummySpec) {
		t.Errorf("body mismatch: expected %q, got %q", dummySpec, body)
	}
}

func TestOpenAPIJSONHandler(t *testing.T) {
	dummyJSON := []byte(`{"openapi":"3.0.3","info":{"title":"Test"}}`)
	handler := OpenAPIJSONHandler(dummyJSON)
	req := httptest.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	if body := w.Body.String(); body != string(dummyJSON) {
		t.Errorf("body mismatch: expected %q, got %q", dummyJSON, body)
	}
}
