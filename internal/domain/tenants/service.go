package tenants

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/google/uuid"
)

const keyEntropy = 32

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, req CreateTenantRequest, platformKeyProvided string, expectedPlatformKey string) (*CreateTenantResponse, error) {
    if platformKeyProvided != expectedPlatformKey {
        return nil, fmt.Errorf("invalid platform key")
    }

    // Apply defaults
    timezone := "Africa/Johannesburg"
    if req.Timezone != nil {
        timezone = *req.Timezone
    }
    locale := "en-ZA"
    if req.Locale != nil {
        locale = *req.Locale
    }
    subscriptionPlan := "trial"
    if req.SubscriptionPlan != nil {
        subscriptionPlan = *req.SubscriptionPlan
    }
    settings := make(map[string]interface{})
    if req.Settings != nil {
        settings = req.Settings
    }

    jwtSecret, err := generateRandomHex(32)
    if err != nil {
        return nil, fmt.Errorf("generate jwt secret: %w", err)
    }

    tenantKey, err := generateRandomHex(keyEntropy)
    if err != nil {
        return nil, fmt.Errorf("generate tenant key: %w", err)
    }

    hash := sha256.Sum256([]byte(tenantKey))
    tenantKeyHash := hex.EncodeToString(hash[:])

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
        Timezone:             timezone,
        Locale:               locale,
        Industry:             req.Industry,
        SubscriptionPlan:     subscriptionPlan,
        Settings:             settings,
        Status:               "active",
        CreatedAt:            time.Now(),
        UpdatedAt:            time.Now(),
    }

    if err := s.repo.Create(ctx, tenant, jwtSecret, tenantKeyHash); err != nil {
        return nil, err
    }

    return &CreateTenantResponse{
        TenantID:            tenant.ID.String(),
        Name:                tenant.Name,
        LegalName:           tenant.LegalName,
        BillingEmail:        tenant.BillingEmail,
        TaxID:               tenant.TaxID,
        AddressLine1:        tenant.AddressLine1,
        AddressLine2:        tenant.AddressLine2,
        City:                tenant.City,
        PostalCode:          tenant.PostalCode,
        PrimaryContactName:  tenant.PrimaryContactName,
        PrimaryContactEmail: tenant.PrimaryContactEmail,
        PrimaryContactPhone: tenant.PrimaryContactPhone,
        CountryCode:         tenant.CountryCode,
        Timezone:            tenant.Timezone,
        Locale:              tenant.Locale,
        Industry:            tenant.Industry,
        SubscriptionPlan:    tenant.SubscriptionPlan,
        Settings:            tenant.Settings,
        Status:              tenant.Status,
        TenantKey:           tenantKey,
        CreatedAt:           tenant.CreatedAt,
    }, nil
}

func (s *Service) GetTenant(ctx context.Context, id string) (*TenantResponse, error) {
    uuid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid tenant id")
    }
    tenant, err := s.repo.GetByID(ctx, uuid)
    if err != nil {
        return nil, err
    }
    if tenant == nil {
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

func (s *Service) ListTenants(ctx context.Context, limit, offset int) ([]TenantResponse, error) {
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
