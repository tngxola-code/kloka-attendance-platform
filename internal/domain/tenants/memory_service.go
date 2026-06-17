package tenants

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"

    "github.com/google/uuid"
)

type MemoryService struct {
    repo *MemoryRepository
}

func NewMemoryService(repo *MemoryRepository) *MemoryService {
    return &MemoryService{repo: repo}
}

func (s *MemoryService) CreateTenant(ctx context.Context, req CreateTenantRequest, platformKeyProvided, expectedPlatformKey string) (*CreateTenantResponse, error) {
    if expectedPlatformKey == "" || !constantTimeEqual(platformKeyProvided, expectedPlatformKey) {
        return nil, errors.New("unauthorized")
    }

    tenantKey, err := generateRandomHex(32)
    if err != nil {
        return nil, err
    }
    keyHash, err := hashOpaqueToken(tenantKey)
    if err != nil {
        return nil, err
    }

    signingSecret, err := generateRandomHex(32)
    if err != nil {
        return nil, err
    }

    tenant := &Tenant{
        ID:                   uuid.New(),
        Name:                 req.Name,
        LegalName:            req.LegalName,
        BillingEmail:         req.BillingEmail,
        TaxID:                req.TaxID,
        AddressLine1:         req.AddressLine1,
        AddressLine2:         req.AddressLine2,
        City:                 req.City,
        PostalCode:           req.PostalCode,
        PrimaryContactName:   req.PrimaryContactName,
        PrimaryContactEmail:  req.PrimaryContactEmail,
        PrimaryContactPhone:  req.PrimaryContactPhone,
        CountryCode:          req.CountryCode,
        Timezone:             "Africa/Johannesburg",
        Locale:               "en-ZA",
        Industry:             req.Industry,
        SubscriptionPlan:     "trial",
        Settings:             make(map[string]interface{}),
        Status:               "active",
        CreatedAt:            time.Now(),
        UpdatedAt:            time.Now(),
    }
    if req.Timezone != nil {
        tenant.Timezone = *req.Timezone
    }
    if req.Locale != nil {
        tenant.Locale = *req.Locale
    }
    if req.SubscriptionPlan != nil {
        tenant.SubscriptionPlan = *req.SubscriptionPlan
    }
    if req.Settings != nil {
        tenant.Settings = req.Settings
    }

    if err := s.repo.Create(ctx, tenant, signingSecret, string(keyHash)); err != nil {
        return nil, err
    }
    return &CreateTenantResponse{
        TenantID:  tenant.ID.String(),
        Name:      tenant.Name,
        TenantKey: tenantKey,
        CreatedAt: tenant.CreatedAt,
    }, nil
}

func (s *MemoryService) ValidateTenantKey(ctx context.Context, tenantID, key string) (*Tenant, []byte, error) {
    uid, err := uuid.Parse(tenantID)
    if err != nil {
        return nil, nil, nil
    }
    secret, ok, err := s.repo.ValidateKey(ctx, uid, key)
    if err != nil || !ok {
        return nil, nil, err
    }
    tenant, err := s.repo.GetByID(ctx, uid)
    if err != nil {
        return nil, nil, err
    }
    return tenant, secret, nil
}

func (s *MemoryService) SigningSecret(ctx context.Context, tenantID string) ([]byte, error) {
    uid, err := uuid.Parse(tenantID)
    if err != nil {
        return nil, nil
    }
    return s.repo.SigningSecret(ctx, uid)
}

func (s *MemoryService) GetTenant(ctx context.Context, id string) (*TenantResponse, error) {
    uid, err := uuid.Parse(id)
    if err != nil {
        return nil, errors.New("invalid tenant id")
    }
    tenant, err := s.repo.GetByID(ctx, uid)
    if err != nil || tenant == nil {
        return nil, nil
    }
    return &TenantResponse{
        ID:                   tenant.ID.String(),
        Name:                 tenant.Name,
        LegalName:            tenant.LegalName,
        BillingEmail:         tenant.BillingEmail,
        TaxID:                tenant.TaxID,
        AddressLine1:         tenant.AddressLine1,
        AddressLine2:         tenant.AddressLine2,
        City:                 tenant.City,
        PostalCode:           tenant.PostalCode,
        PrimaryContactName:   tenant.PrimaryContactName,
        PrimaryContactEmail:  tenant.PrimaryContactEmail,
        PrimaryContactPhone:  tenant.PrimaryContactPhone,
        CountryCode:          tenant.CountryCode,
        Timezone:             tenant.Timezone,
        Locale:               tenant.Locale,
        Industry:             tenant.Industry,
        SubscriptionPlan:     tenant.SubscriptionPlan,
        Settings:             tenant.Settings,
        Status:               tenant.Status,
        CreatedAt:            tenant.CreatedAt,
        UpdatedAt:            tenant.UpdatedAt,
    }, nil
}

func (s *MemoryService) ListTenants(ctx context.Context, limit, offset int) ([]TenantResponse, error) {
    tenants, err := s.repo.List(ctx, limit, offset)
    if err != nil {
        return nil, err
    }
    resp := make([]TenantResponse, len(tenants))
    for i, t := range tenants {
        resp[i] = TenantResponse{
            ID:                   t.ID.String(),
            Name:                 t.Name,
            LegalName:            t.LegalName,
            BillingEmail:         t.BillingEmail,
            TaxID:                t.TaxID,
            AddressLine1:         t.AddressLine1,
            AddressLine2:         t.AddressLine2,
            City:                 t.City,
            PostalCode:           t.PostalCode,
            PrimaryContactName:   t.PrimaryContactName,
            PrimaryContactEmail:  t.PrimaryContactEmail,
            PrimaryContactPhone:  t.PrimaryContactPhone,
            CountryCode:          t.CountryCode,
            Timezone:             t.Timezone,
            Locale:               t.Locale,
            Industry:             t.Industry,
            SubscriptionPlan:     t.SubscriptionPlan,
            Settings:             t.Settings,
            Status:               t.Status,
            CreatedAt:            t.CreatedAt,
            UpdatedAt:            t.UpdatedAt,
        }
    }
    return resp, nil
}

func generateRandomHex(bytes int) (string, error) {
    b := make([]byte, bytes)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}
