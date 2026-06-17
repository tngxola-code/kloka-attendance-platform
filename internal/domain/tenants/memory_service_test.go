package tenants

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMemoryService_CreateTenant_Success(t *testing.T) {
    repo := NewMemoryRepository()
    svc := NewMemoryService(repo)

    req := CreateTenantRequest{
        Name:         "Test Tenant",
        BillingEmail: ptrString("test@example.com"),
        CountryCode:  "ZA",
        Timezone:     ptrString("Africa/Johannesburg"),
        Locale:       ptrString("en-ZA"),
    }
    providedKey := "any-key"
    expectedKey := "any-key" // we're not checking platform key in this demo

    resp, err := svc.CreateTenant(context.Background(), req, providedKey, expectedKey)
    require.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "Test Tenant", resp.Name)
    assert.NotEmpty(t, resp.TenantKey)
    assert.NotEmpty(t, resp.TenantID)
}

func TestMemoryService_GetTenant(t *testing.T) {
    repo := NewMemoryRepository()
    svc := NewMemoryService(repo)

    // Create a tenant first
    req := CreateTenantRequest{
        Name:        "Get Test",
        CountryCode: "ZA",
    }
    createResp, err := svc.CreateTenant(context.Background(), req, "key", "key")
    require.NoError(t, err)

    // Retrieve by ID
    getResp, err := svc.GetTenant(context.Background(), createResp.TenantID)
    require.NoError(t, err)
    assert.Equal(t, createResp.TenantID, getResp.ID)
    assert.Equal(t, "Get Test", getResp.Name)
}

func TestMemoryService_GetTenant_NotFound(t *testing.T) {
    repo := NewMemoryRepository()
    svc := NewMemoryService(repo)

    resp, err := svc.GetTenant(context.Background(), "00000000-0000-0000-0000-000000000000")
    assert.NoError(t, err)
    assert.Nil(t, resp)
}

func TestMemoryService_ListTenants(t *testing.T) {
    repo := NewMemoryRepository()
    svc := NewMemoryService(repo)

    // Create two tenants
    req1 := CreateTenantRequest{Name: "Tenant A", CountryCode: "ZA"}
    req2 := CreateTenantRequest{Name: "Tenant B", CountryCode: "ZA"}
    _, err := svc.CreateTenant(context.Background(), req1, "key", "key")
    require.NoError(t, err)
    _, err = svc.CreateTenant(context.Background(), req2, "key", "key")
    require.NoError(t, err)

    tenants, err := svc.ListTenants(context.Background(), 10, 0)
    assert.NoError(t, err)
    assert.Len(t, tenants, 2)
}

func TestMemoryService_ListTenants_Empty(t *testing.T) {
    repo := NewMemoryRepository()
    svc := NewMemoryService(repo)

    tenants, err := svc.ListTenants(context.Background(), 10, 0)
    assert.NoError(t, err)
    assert.Empty(t, tenants)
}

// Helper (already exists in types.go but we'll duplicate for safety)
func ptrString(s string) *string {
    return &s
}
