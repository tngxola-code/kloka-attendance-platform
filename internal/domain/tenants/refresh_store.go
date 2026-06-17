package tenants

import (
    "context"
    "errors"
    "sync"
    "time"
)

type RefreshStore interface {
    Store(ctx context.Context, tenantID, rawToken string, exp time.Time) error
    Consume(ctx context.Context, rawToken string) (tenantID string, err error)
    RevokeAllForTenant(ctx context.Context, tenantID string) error
}

var ErrRefreshNotFound = errors.New("refresh token not found or expired")

type memoryRefreshRecord struct {
    tenantID string
    hash     string
    expires  time.Time
}

type MemoryRefreshStore struct {
    mu      sync.Mutex
    records []memoryRefreshRecord
}

func NewMemoryRefreshStore() *MemoryRefreshStore {
    return &MemoryRefreshStore{}
}

func (s *MemoryRefreshStore) Store(_ context.Context, tenantID, rawToken string, exp time.Time) error {
    hash, err := hashOpaqueToken(rawToken)
    if err != nil {
        return err
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    s.gcLocked()
    s.records = append(s.records, memoryRefreshRecord{tenantID: tenantID, hash: string(hash), expires: exp})
    return nil
}

func (s *MemoryRefreshStore) Consume(_ context.Context, rawToken string) (string, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.gcLocked()
    now := time.Now()
    for i := range s.records {
        rec := s.records[i]
        if now.After(rec.expires) {
            continue
        }
        if compareOpaqueToken(rec.hash, rawToken) {
            s.records = append(s.records[:i], s.records[i+1:]...)
            return rec.tenantID, nil
        }
    }
    return "", ErrRefreshNotFound
}

func (s *MemoryRefreshStore) RevokeAllForTenant(_ context.Context, tenantID string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    kept := s.records[:0]
    for _, rec := range s.records {
        if rec.tenantID != tenantID {
            kept = append(kept, rec)
        }
    }
    s.records = kept
    return nil
}

func (s *MemoryRefreshStore) gcLocked() {
    now := time.Now()
    kept := s.records[:0]
    for _, rec := range s.records {
        if now.Before(rec.expires) {
            kept = append(kept, rec)
        }
    }
    s.records = kept
}
