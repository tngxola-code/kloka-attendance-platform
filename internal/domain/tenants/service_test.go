package tenants

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// mockRepository implements Repository interface for testing.
type mockRepository struct {
    createFunc func(ctx context.Context, tenant *Tenant, jwtSecret, keyHash string) error
    getByIDFunc func(ctx context.Context, id string) (*Tenant, error)
    listFunc func(ctx context.Context, limit, offset int) ([]Tenant, error)
}

func (m *mockRepository) Create(ctx context.Context, tenant *Tenant, jwtSecret, keyHash string) error {
    if m.createFunc != nil {
        return m.createFunc(ctx, tenant, jwtSecret, keyHash)
    }
    return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Tenant, error) {
    if m.getByIDFunc != nil {
        return m.getByIDFunc(ctx, id)
    }
    return nil, nil
}

func (m *mockRepository) List(ctx context.Context, limit, offset int) ([]Tenant, error) {
    if m.listFunc != nil {
        return m.listFunc(ctx, limit, offset)
    }
    return nil, nil
}

func TestCreateTenant_Success(t *testing.T) {
    repo := &mockRepository{
        createFunc: func(ctx context.Context, tenant *Tenant, jwtSecret, keyHash string) error {
            // Simulate successful creation
            return nil
        },
    }
    svc := NewService(repo)

    req := CreateTenantRequest{
        Name:        "Test Tenant",
        BillingEmail: ptrString("test@example.com"),
        CountryCode: "ZA",
    }
    providedKey := "correct-key"
    expectedKey := "correct-key"

    resp, err := svc.CreateTenant(context.Background(), req, providedKey, expectedKey)
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, "Test Tenant", resp.Name)
    assert.NotEmpty(t, resp.TenantKey)
    assert.NotEmpty(t, resp.TenantID)
}

func TestCreateTenant_InvalidPlatformKey(t *testing.T) {
    repo := &mockRepository{}
    svc := NewService(repo)

    req := CreateTenantRequest{Name: "Test"}
    _, err := svc.CreateTenant(context.Background(), req, "wrong", "correct")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid platform key")
}

func TestGetTenant_NotFound(t *testing.T) {
    repo := &mockRepository{
        getByIDFunc: func(ctx context.Context, id string) (*Tenant, error) {
            return nil, nil
        },
    }
    svc := NewService(repo)

    resp, err := svc.GetTenant(context.Background(), "00000000-0000-0000-0000-000000000000")
    assert.NoError(t, err)
    assert.Nil(t, resp)
}

func TestListTenants_Empty(t *testing.T) {
    repo := &mockRepository{
        listFunc: func(ctx context.Context, limit, offset int) ([]Tenant, error) {
            return []Tenant{}, nil
        },
    }
    svc := NewService(repo)

    tenants, err := svc.ListTenants(context.Background(), 10, 0)
    assert.NoError(t, err)
    assert.Empty(t, tenants)
}

// Helper
func ptrString(s string) *string {
    return &s
}
