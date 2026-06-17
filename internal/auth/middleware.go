package auth

import (
    "context"
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/golang-jwt/jwt/v5"

    "kloka-attendance-platform/internal/httpx"
)

type SecretLookup func(ctx context.Context, tenantID string) ([]byte, error)

type Claims struct {
    TenantID string   `json:"tenant_id"`
    Role     string   `json:"role"`
    Scope    []string `json:"scope"`
    jwt.RegisteredClaims
}

func (c *Claims) HasScope(s string) bool {
    for _, v := range c.Scope {
        if v == s {
            return true
        }
    }
    return false
}

type ctxKey struct{}

var claimsKey ctxKey

const (
    expectedIssuer   = "kloka-tenant-auth"
    expectedAudience = "kloka-api"
    clockSkewLeeway  = 30 * time.Second
)

func Authenticate(lookup SecretLookup) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            raw, err := bearer(r)
            if err != nil {
                unauthorized(w, r)
                return
            }

            keyFunc := func(t *jwt.Token) (interface{}, error) {
                if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
                    return nil, errors.New("unexpected signing method")
                }
                claims, ok := t.Claims.(*Claims)
                if !ok || claims.TenantID == "" {
                    return nil, errors.New("missing tenant_id claim")
                }
                secret, err := lookup(r.Context(), claims.TenantID)
                if err != nil || len(secret) == 0 {
                    return nil, errors.New("unknown tenant")
                }
                return secret, nil
            }

            token, err := jwt.ParseWithClaims(raw, &Claims{}, keyFunc,
                jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
                jwt.WithIssuer(expectedIssuer),
                jwt.WithAudience(expectedAudience),
                jwt.WithExpirationRequired(),
                jwt.WithLeeway(clockSkewLeeway),
            )
            if err != nil {
                unauthorized(w, r)
                return
            }
            claims, ok := token.Claims.(*Claims)
            if !ok || !token.Valid {
                unauthorized(w, r)
                return
            }

            ctx := context.WithValue(r.Context(), claimsKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RequireScope(scope string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            c := ClaimsFromContext(r.Context())
            if c == nil || !c.HasScope(scope) {
                forbidden(w, r)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func RequireTenantMatch(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c := ClaimsFromContext(r.Context())
        if c == nil {
            unauthorized(w, r)
            return
        }
        pathID := chi.URLParam(r, "id")
        if pathID != "" && pathID != c.TenantID {
            httpx.WriteProblem(w, r, http.StatusNotFound, "Not Found", "tenant not found")
            return
        }
        next.ServeHTTP(w, r)
    })
}

func ClaimsFromContext(ctx context.Context) *Claims {
    c, _ := ctx.Value(claimsKey).(*Claims)
    return c
}

func bearer(r *http.Request) (string, error) {
    h := r.Header.Get("Authorization")
    const p = "Bearer "
    if len(h) <= len(p) || !strings.EqualFold(h[:len(p)], p) {
        return "", errors.New("no bearer token")
    }
    return strings.TrimSpace(h[len(p):]), nil
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
    httpx.WriteProblem(w, r, http.StatusUnauthorized, "Unauthorized", "invalid or missing token")
}

func forbidden(w http.ResponseWriter, r *http.Request) {
    httpx.WriteProblem(w, r, http.StatusForbidden, "Forbidden", "insufficient scope")
}
