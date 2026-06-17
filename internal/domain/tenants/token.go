package tenants

import (
    "context"
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "encoding/json"
    "errors"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"

    "kloka-attendance-platform/internal/httpx"
)

const (
    accessTokenTTL  = 15 * time.Minute
    refreshTokenTTL = 30 * 24 * time.Hour
)

type TenantClaims struct {
    TenantID string   `json:"tenant_id"`
    Type     string   `json:"type"`
    Role     string   `json:"role"`
    Scope    []string `json:"scope"`
    jwt.RegisteredClaims
}

func defaultTenantScope() []string {
    return []string{
        "tenant:read", "tenant:write",
        "workers:read", "workers:write",
        "sites:read", "sites:write",
        "attendance:read", "reports:read",
    }
}

func buildAccessToken(tenantID string, secret []byte, now time.Time) (token, jti string, err error) {
    if len(secret) == 0 {
        return "", "", errors.New("empty signing secret")
    }
    jti = uuid.NewString()
    claims := TenantClaims{
        TenantID: tenantID,
        Type:     "tenant",
        Role:     "tenant",
        Scope:    defaultTenantScope(),
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "kloka-tenant-auth",
            Subject:   tenantID,
            Audience:  jwt.ClaimStrings{"kloka-api"},
            ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
            NotBefore: jwt.NewNumericDate(now),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        jti,
        },
    }
    t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := t.SignedString(secret)
    return signed, jti, err
}

type TokenRequest struct {
    TenantID  string `json:"tenant_id"`
    TenantKey string `json:"tenant_key"`
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
    var req TokenRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    if req.TenantID == "" || req.TenantKey == "" {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "tenant_id and tenant_key required")
        return
    }

    tenant, secret, err := h.svc.ValidateTenantKey(r.Context(), req.TenantID, req.TenantKey)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "authentication failed")
        return
    }
    if tenant == nil {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "invalid tenant credentials")
        return
    }

    resp, err := h.issueTokens(r.Context(), tenant.ID.String(), secret)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "failed to issue tokens")
        return
    }
    httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
    var req RefreshRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "invalid JSON")
        return
    }
    if req.RefreshToken == "" {
        httpx.WriteProblem(w, r, http.StatusBadRequest, "Bad Request", "refresh_token required")
        return
    }

    tenantID, err := h.refresh.Consume(r.Context(), req.RefreshToken)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "invalid or expired refresh token")
        return
    }

    secret, err := h.svc.SigningSecret(r.Context(), tenantID)
    if err != nil || len(secret) == 0 {
        httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "tenant not found")
        return
    }

    resp, err := h.issueTokens(r.Context(), tenantID, secret)
    if err != nil {
        httpx.WriteProblem(w, r, http.StatusInternalServerError, "Internal Error", "failed to issue tokens")
        return
    }
    httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) issueTokens(ctx context.Context, tenantID string, secret []byte) (TokenResponse, error) {
    now := time.Now()
    access, _, err := buildAccessToken(tenantID, secret, now)
    if err != nil {
        return TokenResponse{}, err
    }

    rawRefresh, err := newOpaqueToken()
    if err != nil {
        return TokenResponse{}, err
    }
    if err := h.refresh.Store(ctx, tenantID, rawRefresh, now.Add(refreshTokenTTL)); err != nil {
        return TokenResponse{}, err
    }

    return TokenResponse{
        AccessToken:  access,
        RefreshToken: rawRefresh,
        TokenType:    "Bearer",
        ExpiresIn:    int(accessTokenTTL.Seconds()),
    }, nil
}

func newOpaqueToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashOpaqueToken(raw string) (string, error) {
    h, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(h), nil
}

func compareOpaqueToken(hash, raw string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

func constantTimeEqual(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
