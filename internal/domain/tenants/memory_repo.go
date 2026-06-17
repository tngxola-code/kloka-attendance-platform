package tenants

import (
    "context"
    "sync"

    "github.com/google/uuid"
)

type tenantSecrets struct {
    signingSecret string
    keyHash       string
}

type MemoryRepository struct {
    mu      sync.RWMutex
    tenants map[uuid.UUID]*Tenant
    secrets map[uuid.UUID]tenantSecrets
}

func NewMemoryRepository() *MemoryRepository {
    return &MemoryRepository{
        tenants: make(map[uuid.UUID]*Tenant),
        secrets: make(map[uuid.UUID]tenantSecrets),
    }
}

func (r *MemoryRepository) Create(ctx context.Context, tenant *Tenant, signingSecret, keyHash string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.tenants[tenant.ID] = tenant
    r.secrets[tenant.ID] = tenantSecrets{signingSecret: signingSecret, keyHash: keyHash}
    return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    t, ok := r.tenants[id]
    if !ok {
        return nil, nil
    }
    copy := *t
    return &copy, nil
}

func (r *MemoryRepository) List(ctx context.Context, limit, offset int) ([]Tenant, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var list []Tenant
    for _, t := range r.tenants {
        list = append(list, *t)
    }
    if offset >= len(list) {
        return []Tenant{}, nil
    }
    end := offset + limit
    if end > len(list) {
        end = len(list)
    }
    return list[offset:end], nil
}

func (r *MemoryRepository) ValidateKey(ctx context.Context, id uuid.UUID, key string) ([]byte, bool, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, exists := r.secrets[id]
    if !exists {
        _ = compareOpaqueToken("$2a$10$........................................................", key)
        return nil, false, nil
    }
    if !compareOpaqueToken(s.keyHash, key) {
        return nil, false, nil
    }
    return []byte(s.signingSecret), true, nil
}

func (r *MemoryRepository) SigningSecret(ctx context.Context, id uuid.UUID) ([]byte, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, ok := r.secrets[id]
    if !ok {
        return nil, nil
    }
    return []byte(s.signingSecret), nil
}
