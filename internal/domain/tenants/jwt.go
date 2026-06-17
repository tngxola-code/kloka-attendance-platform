package tenants

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

func generateTenantJWT(tenant *Tenant, secret []byte) (string, error) {
    now := time.Now()
    expiresIn := 86400 // 24 hours

    claims := jwt.MapClaims{
        "iss":       "kloka-tenant-auth",
        "sub":       tenant.ID.String(),
        "aud":       "kloka-api",
        "exp":       now.Add(time.Duration(expiresIn) * time.Second).Unix(),
        "iat":       now.Unix(),
        "nbf":       now.Unix(),
        "jti":       uuid.New().String(),
        "type":      "tenant",
        "tenant_id": tenant.ID.String(),
        "role":      "tenant",
        "scope": []string{
            "tenant:read",
            "tenant:write",
            "workers:read",
            "workers:write",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
