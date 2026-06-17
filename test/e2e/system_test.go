package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kloka-attendance-platform/internal/metrics"
	"kloka-attendance-platform/internal/system"
)

func setupSystemServer() *httptest.Server {
	mux := http.NewServeMux()
	const apiVersion = "/api/v1"

	dummyYAML := []byte("openapi: 3.0.3\ninfo:\n  title: Test\n  version: 1.0.0")
	dummyJSON := []byte(`{"openapi":"3.0.3","info":{"title":"Test","version":"1.0.0"}}`)

	mux.HandleFunc("GET "+apiVersion+"/health/live", system.LivenessHandler)
	mux.HandleFunc("GET "+apiVersion+"/health/ready", system.ReadinessHandler(nil))
	mux.HandleFunc("GET "+apiVersion+"/version", system.VersionHandler)
	mux.HandleFunc("GET "+apiVersion+"/openapi.yaml", system.OpenAPISpecHandler(dummyYAML))
	mux.HandleFunc("GET "+apiVersion+"/openapi.json", system.OpenAPIJSONHandler(dummyJSON))
	mux.Handle("GET "+apiVersion+"/metrics", metrics.MetricsHandler())

	return httptest.NewServer(mux)
}

func TestSystemHealthLive(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/health/live")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "alive" {
		t.Errorf("expected status alive, got %q", body["status"])
	}
}

func TestSystemHealthReady(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/health/ready")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", resp.StatusCode)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "not ready" {
		t.Errorf("expected status 'not ready', got %q", body["status"])
	}
	if db, ok := body["database"].(bool); !ok || db != false {
		t.Errorf("expected database false, got %v", body["database"])
	}
}

func TestSystemVersion(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/version")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["version"] == "" {
		t.Error("version field missing or empty")
	}
}

func TestSystemOpenAPIYAML(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/x-yaml" {
		t.Errorf("expected application/x-yaml, got %q", ct)
	}
}

func TestSystemOpenAPIJSON(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
	var doc map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	info, ok := doc["info"].(map[string]interface{})
	if !ok {
		t.Fatal("missing info")
	}
	if title, ok := info["title"].(string); !ok || title != "Test" {
		t.Errorf("expected title 'Test', got %v", title)
	}
}

func TestSystemMetrics(t *testing.T) {
	srv := setupSystemServer()
	defer srv.Close()

	// Generate some metrics
	http.Get(srv.URL + "/api/v1/health/live")
	http.Get(srv.URL + "/api/v1/version")

	resp, err := http.Get(srv.URL + "/api/v1/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("expected Content-Type starting with text/plain, got %q", ct)
	}
	buf := make([]byte, 100)
	n, _ := resp.Body.Read(buf)
	if n == 0 {
		t.Error("metrics body empty")
	}
}
