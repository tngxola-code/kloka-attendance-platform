package tenants

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"kloka-attendance-platform/internal/httpx"
)

type Handler struct {
	svc     *MemoryService
	refresh RefreshStore
}

func NewHandler(svc *MemoryService, refresh RefreshStore) *Handler {
	return &Handler{svc: svc, refresh: refresh}
}

func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	providedKey := r.Header.Get("X-Platform-Key")
	var req CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
		return
	}
	if req.Name == "" {
		httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "name is required")
		return
	}
	if req.CountryCode == "" {
		req.CountryCode = "ZA"
	}

	resp, err := h.svc.CreateTenant(r.Context(), req, providedKey, "test-platform-key")
	if err != nil {
		httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}
	w.Header().Set("Location", "/api/v1/tenants/"+resp.TenantID)
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) GetTenant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant, err := h.svc.GetTenant(r.Context(), id)
	if err != nil {
		httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", err.Error())
		return
	}
	if tenant == nil {
		httpx.WriteProblem(w, r, http.StatusNotFound, "Not Found", "tenant not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tenant)
}

func (h *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	tenants, err := h.svc.ListTenants(r.Context(), limit, offset)
	if err != nil {
		httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":     tenants,
		"limit":    limit,
		"offset":   offset,
		"has_more": len(tenants) == limit,
	})
}
