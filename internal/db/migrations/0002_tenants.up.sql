CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    billing_email TEXT,
    country_code CHAR(2) NOT NULL DEFAULT 'ZA',
    status TEXT NOT NULL DEFAULT 'active',
    jwt_signing_secret TEXT NOT NULL,          -- server‑side secret, never exposed
    tenant_key_hash TEXT NOT NULL UNIQUE,      -- SHA‑256 of the one‑time tenant key
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_status ON tenants(status);
