package tenants

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/lib/pq"
)

type Repository struct {
    pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
    return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, tenant *Tenant, jwtSecret, keyHash string) error {
    query := `
        INSERT INTO tenants (
            id, name, legal_name, billing_email, tax_id,
            address_line1, address_line2, city, postal_code,
            primary_contact_name, primary_contact_email, primary_contact_phone,
            country_code, timezone, locale, industry, subscription_plan, settings,
            status, jwt_signing_secret, tenant_key_hash
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
        )
    `
    _, err := r.pool.Exec(ctx, query,
        tenant.ID, tenant.Name, tenant.LegalName, tenant.BillingEmail, tenant.TaxID,
        tenant.AddressLine1, tenant.AddressLine2, tenant.City, tenant.PostalCode,
        tenant.PrimaryContactName, tenant.PrimaryContactEmail, tenant.PrimaryContactPhone,
        tenant.CountryCode, tenant.Timezone, tenant.Locale, tenant.Industry, tenant.SubscriptionPlan,
        tenant.Settings, tenant.Status, jwtSecret, keyHash,
    )
    return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
    query := `
        SELECT
            id, name, legal_name, billing_email, tax_id,
            address_line1, address_line2, city, postal_code,
            primary_contact_name, primary_contact_email, primary_contact_phone,
            country_code, timezone, locale, industry, subscription_plan, settings,
            status, created_at, updated_at
        FROM tenants WHERE id = $1
    `
    var t Tenant
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &t.ID, &t.Name, &t.LegalName, &t.BillingEmail, &t.TaxID,
        &t.AddressLine1, &t.AddressLine2, &t.City, &t.PostalCode,
        &t.PrimaryContactName, &t.PrimaryContactEmail, &t.PrimaryContactPhone,
        &t.CountryCode, &t.Timezone, &t.Locale, &t.Industry, &t.SubscriptionPlan,
        &t.Settings, &t.Status, &t.CreatedAt, &t.UpdatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("query tenant: %w", err)
    }
    return &t, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int) ([]Tenant, error) {
    query := `
        SELECT
            id, name, legal_name, billing_email, tax_id,
            address_line1, address_line2, city, postal_code,
            primary_contact_name, primary_contact_email, primary_contact_phone,
            country_code, timezone, locale, industry, subscription_plan, settings,
            status, created_at, updated_at
        FROM tenants ORDER BY created_at DESC LIMIT $1 OFFSET $2
    `
    rows, err := r.pool.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tenants []Tenant
    for rows.Next() {
        var t Tenant
        if err := rows.Scan(
            &t.ID, &t.Name, &t.LegalName, &t.BillingEmail, &t.TaxID,
            &t.AddressLine1, &t.AddressLine2, &t.City, &t.PostalCode,
            &t.PrimaryContactName, &t.PrimaryContactEmail, &t.PrimaryContactPhone,
            &t.CountryCode, &t.Timezone, &t.Locale, &t.Industry, &t.SubscriptionPlan,
            &t.Settings, &t.Status, &t.CreatedAt, &t.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        tenants = append(tenants, t)
    }
    return tenants, rows.Err()
}
