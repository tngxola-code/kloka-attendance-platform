package tenants

import (
    "time"

    "github.com/google/uuid"
)

type Tenant struct {
    ID                   uuid.UUID  `json:"id"`
    Name                 string     `json:"name"`
    LegalName            *string    `json:"legal_name,omitempty"`
    BillingEmail         *string    `json:"billing_email,omitempty"`
    TaxID                *string    `json:"tax_id,omitempty"`
    AddressLine1         *string    `json:"address_line1,omitempty"`
    AddressLine2         *string    `json:"address_line2,omitempty"`
    City                 *string    `json:"city,omitempty"`
    PostalCode           *string    `json:"postal_code,omitempty"`
    PrimaryContactName   *string    `json:"primary_contact_name,omitempty"`
    PrimaryContactEmail  *string    `json:"primary_contact_email,omitempty"`
    PrimaryContactPhone  *string    `json:"primary_contact_phone,omitempty"`
    CountryCode          string     `json:"country_code"`
    Timezone             string     `json:"timezone"`
    Locale               string     `json:"locale"`
    Industry             *string    `json:"industry,omitempty"`
    SubscriptionPlan     string     `json:"subscription_plan"`
    Settings             map[string]interface{} `json:"settings,omitempty"`
    Status               string     `json:"status"`
    CreatedAt            time.Time  `json:"created_at"`
    UpdatedAt            time.Time  `json:"updated_at"`
}

type CreateTenantRequest struct {
    Name                 string                  `json:"name" validate:"required"`
    LegalName            *string                 `json:"legal_name,omitempty"`
    BillingEmail         *string                 `json:"billing_email" validate:"omitempty,email"`
    TaxID                *string                 `json:"tax_id,omitempty"`
    AddressLine1         *string                 `json:"address_line1,omitempty"`
    AddressLine2         *string                 `json:"address_line2,omitempty"`
    City                 *string                 `json:"city,omitempty"`
    PostalCode           *string                 `json:"postal_code,omitempty"`
    PrimaryContactName   *string                 `json:"primary_contact_name,omitempty"`
    PrimaryContactEmail  *string                 `json:"primary_contact_email" validate:"omitempty,email"`
    PrimaryContactPhone  *string                 `json:"primary_contact_phone,omitempty"`
    CountryCode          string                  `json:"country_code" validate:"required,len=2,alpha"`
    Timezone             *string                 `json:"timezone,omitempty"`
    Locale               *string                 `json:"locale,omitempty"`
    Industry             *string                 `json:"industry,omitempty"`
    SubscriptionPlan     *string                 `json:"subscription_plan,omitempty"`
    Settings             map[string]interface{}  `json:"settings,omitempty"`
}

type CreateTenantResponse struct {
    TenantID             string                  `json:"tenant_id"`
    Name                 string                  `json:"name"`
    LegalName            *string                 `json:"legal_name,omitempty"`
    BillingEmail         *string                 `json:"billing_email,omitempty"`
    TaxID                *string                 `json:"tax_id,omitempty"`
    AddressLine1         *string                 `json:"address_line1,omitempty"`
    AddressLine2         *string                 `json:"address_line2,omitempty"`
    City                 *string                 `json:"city,omitempty"`
    PostalCode           *string                 `json:"postal_code,omitempty"`
    PrimaryContactName   *string                 `json:"primary_contact_name,omitempty"`
    PrimaryContactEmail  *string                 `json:"primary_contact_email,omitempty"`
    PrimaryContactPhone  *string                 `json:"primary_contact_phone,omitempty"`
    CountryCode          string                  `json:"country_code"`
    Timezone             string                  `json:"timezone"`
    Locale               string                  `json:"locale"`
    Industry             *string                 `json:"industry,omitempty"`
    SubscriptionPlan     string                  `json:"subscription_plan"`
    Settings             map[string]interface{}  `json:"settings,omitempty"`
    Status               string                  `json:"status"`
    TenantKey            string                  `json:"tenant_key"` // one‑time, shown only once
    CreatedAt            time.Time               `json:"created_at"`
}

type TenantResponse struct {
    ID                   string                  `json:"id"`
    Name                 string                  `json:"name"`
    LegalName            *string                 `json:"legal_name,omitempty"`
    BillingEmail         *string                 `json:"billing_email,omitempty"`
    TaxID                *string                 `json:"tax_id,omitempty"`
    AddressLine1         *string                 `json:"address_line1,omitempty"`
    AddressLine2         *string                 `json:"address_line2,omitempty"`
    City                 *string                 `json:"city,omitempty"`
    PostalCode           *string                 `json:"postal_code,omitempty"`
    PrimaryContactName   *string                 `json:"primary_contact_name,omitempty"`
    PrimaryContactEmail  *string                 `json:"primary_contact_email,omitempty"`
    PrimaryContactPhone  *string                 `json:"primary_contact_phone,omitempty"`
    CountryCode          string                  `json:"country_code"`
    Timezone             string                  `json:"timezone"`
    Locale               string                  `json:"locale"`
    Industry             *string                 `json:"industry,omitempty"`
    SubscriptionPlan     string                  `json:"subscription_plan"`
    Settings             map[string]interface{}  `json:"settings,omitempty"`
    Status               string                  `json:"status"`
    CreatedAt            time.Time               `json:"created_at"`
    UpdatedAt            time.Time               `json:"updated_at"`
}
