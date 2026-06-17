package workers

import (
    "encoding/json"
    "net/http"
    "os"
    "strconv"
    "time"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
    "kloka-attendance-platform/internal/httpx"
)

type Handler struct {
    svc *Service
}

func NewHandler(svc *Service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) CreateWorker(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
    var req CreateWorkerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    worker, err := h.svc.CreateWorker(r.Context(), tenantID, req)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", err.Error())
        return
    }
    w.Header().Set("Location", "/api/v1/workers/"+worker.ID)
    httpx.WriteJSON(w, http.StatusCreated, worker)
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid worker id")
        return
    }
    worker, err := h.svc.GetWorker(r.Context(), id, tenantID)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", err.Error())
        return
    }
    if worker == nil {
        httpx.WriteProblem(w, r, http.StatusNotFound, "Not Found", "worker not found")
        return
    }
    httpx.WriteJSON(w, http.StatusOK, worker)
}

func (h *Handler) UpdateWorker(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid worker id")
        return
    }
    var req UpdateWorkerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    worker, err := h.svc.UpdateWorker(r.Context(), id, tenantID, req)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", err.Error())
        return
    }
    httpx.WriteJSON(w, http.StatusOK, worker)
}

func (h *Handler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid worker id")
        return
    }
    if err := h.svc.DeleteWorker(r.Context(), id, tenantID); err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", err.Error())
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListWorkers(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
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
    workers, err := h.svc.ListWorkers(r.Context(), tenantID, limit, offset)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", err.Error())
        return
    }
    httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "data":     workers,
        "limit":    limit,
        "offset":   offset,
        "has_more": len(workers) == limit,
    })
}

func (h *Handler) LoginWorker(w http.ResponseWriter, r *http.Request) {
    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    worker, err := h.svc.Login(r.Context(), tenantID, req.Phone, req.Pin)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", err.Error())
        return
    }

    jwtSecret := []byte(os.Getenv("JWT_SECRET"))
    if len(jwtSecret) == 0 {
        jwtSecret = []byte("change-me-in-production")
    }

    // Generate refresh token
    refreshToken := uuid.New().String()
    workerRefreshStore.Store(worker.ID.String(), refreshToken)

    // Generate access token
    accessExpiresIn := 900
    now := time.Now()
    claims := jwt.MapClaims{
        "iss":             "kloka-worker-auth",
        "sub":             worker.ID.String(),
        "aud":             "kloka-api",
        "exp":             now.Add(time.Duration(accessExpiresIn) * time.Second).Unix(),
        "iat":             now.Unix(),
        "nbf":             now.Unix(),
        "jti":             uuid.New().String(),
        "type":            "worker",
        "tenant_id":       worker.TenantID.String(),
        "employment_type": worker.EmploymentType,
        "payment_schedule": worker.PaymentSchedule,
        "role":            "worker",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    accessToken, err := token.SignedString(jwtSecret)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "failed to generate token")
        return
    }

    httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "token_type":    "Bearer",
        "expires_in":    accessExpiresIn,
        "worker":        workerResponseFromWorker(worker),
    })
}

func (h *Handler) RefreshWorkerToken(w http.ResponseWriter, r *http.Request) {
    var req struct{ RefreshToken string `json:"refresh_token"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    if req.RefreshToken == "" {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "refresh_token required")
        return
    }

    // Validate refresh token
    workerID := ""
    workerRefreshStore.mu.RLock()
    for rt, wid := range workerRefreshStore.tokens {
        if rt == req.RefreshToken {
            workerID = wid
            break
        }
    }
    workerRefreshStore.mu.RUnlock()

    if workerID == "" {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "invalid refresh token")
        return
    }

    tenantIDStr := r.Context().Value("tenant_id").(string)
    tenantID, err := uuid.Parse(tenantIDStr)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid tenant")
        return
    }

    workerUUID, err := uuid.Parse(workerID)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "invalid worker ID")
        return
    }

    worker, err := h.svc.GetWorker(r.Context(), workerUUID, tenantID)
    if err != nil || worker == nil {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "worker not found")
        return
    }

    jwtSecret := []byte(os.Getenv("JWT_SECRET"))
    if len(jwtSecret) == 0 {
        jwtSecret = []byte("change-me-in-production")
    }

    now := time.Now()
    accessExpiresIn := 900
    claims := jwt.MapClaims{
        "iss":             "kloka-worker-auth",
        "sub":             worker.ID,
        "aud":             "kloka-api",
        "exp":             now.Add(time.Duration(accessExpiresIn) * time.Second).Unix(),
        "iat":             now.Unix(),
        "nbf":             now.Unix(),
        "jti":             uuid.New().String(),
        "type":            "worker",
        "tenant_id":       worker.TenantID,
        "employment_type": worker.EmploymentType,
        "payment_schedule": worker.PaymentSchedule,
        "role":            "worker",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    accessToken, err := token.SignedString(jwtSecret)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "failed to generate token")
        return
    }

    newRefreshToken := uuid.New().String()
    workerRefreshStore.Revoke(req.RefreshToken)
    workerRefreshStore.Store(workerID, newRefreshToken)

    httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{
        "access_token":  accessToken,
        "refresh_token": newRefreshToken,
        "token_type":    "Bearer",
        "expires_in":    accessExpiresIn,
    })
}

// WorkerRefreshStore definition (in-memory)
type workerRefreshStore struct {
    mu     sync.RWMutex
    tokens map[string]string
}

var workerRefreshStore = &workerRefreshStore{
    tokens: make(map[string]string),
}

func (s *workerRefreshStore) Store(workerID, refreshToken string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.tokens[refreshToken] = workerID
}

func (s *workerRefreshStore) Revoke(refreshToken string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.tokens, refreshToken)
}

func workerResponseFromWorker(w *Worker) *WorkerResponse {
    siteIDStr := ""
    if w.SiteID != nil {
        siteIDStr = w.SiteID.String()
    }
    return &WorkerResponse{
        ID:                w.ID.String(),
        TenantID:          w.TenantID.String(),
        EmployeeNumber:    w.EmployeeNumber,
        FullName:          w.FullName,
        Phone:             w.Phone,
        SiteID:            &siteIDStr,
        EmploymentType:    w.EmploymentType,
        PaymentSchedule:   w.PaymentSchedule,
        HourlyRate:        w.HourlyRate,
        ContractEndDate:   w.ContractEndDate,
        Status:            w.Status,
        DateOfBirth:       w.DateOfBirth,
        BiometricConsent:  w.BiometricConsent,
        ConsentTimestamp:  w.ConsentTimestamp,
        CreatedAt:         w.CreatedAt,
        UpdatedAt:         w.UpdatedAt,
    }
}
