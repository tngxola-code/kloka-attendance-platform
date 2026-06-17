package workers

import (
    "context"
    "sync"
    "github.com/google/uuid"
)

type MemoryRepository struct {
    mu      sync.RWMutex
    workers map[uuid.UUID]*Worker
    phoneIndex map[string]uuid.UUID
}

func NewMemoryRepository() *MemoryRepository {
    return &MemoryRepository{
        workers:   make(map[uuid.UUID]*Worker),
        phoneIndex: make(map[string]uuid.UUID),
    }
}

func (r *MemoryRepository) Create(ctx context.Context, worker *Worker) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.workers[worker.ID] = worker
    r.phoneIndex[worker.Phone] = worker.ID
    return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil, nil
    }
    copy := *w
    return &copy, nil
}

func (r *MemoryRepository) GetByPhone(ctx context.Context, phone string, tenantID uuid.UUID) (*Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    id, ok := r.phoneIndex[phone]
    if !ok {
        return nil, nil
    }
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil, nil
    }
    copy := *w
    return &copy, nil
}

func (r *MemoryRepository) Update(ctx context.Context, worker *Worker) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    // If phone changed, update index
    if old, ok := r.workers[worker.ID]; ok && old.Phone != worker.Phone {
        delete(r.phoneIndex, old.Phone)
        r.phoneIndex[worker.Phone] = worker.ID
    }
    r.workers[worker.ID] = worker
    return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil
    }
    delete(r.phoneIndex, w.Phone)
    delete(r.workers, id)
    return nil
}

func (r *MemoryRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var list []Worker
    for _, w := range r.workers {
        if w.TenantID == tenantID {
            list = append(list, *w)
        }
    }
    if offset >= len(list) {
        return []Worker{}, nil
    }
    end := offset + limit
    if end > len(list) {
        end = len(list)
    }
    return list[offset:end], nil
}
EOFcat > internal/domain/workers/memory_repo.go << 'EOF'
package workers

import (
    "context"
    "sync"
    "github.com/google/uuid"
)

type MemoryRepository struct {
    mu      sync.RWMutex
    workers map[uuid.UUID]*Worker
    phoneIndex map[string]uuid.UUID
}

func NewMemoryRepository() *MemoryRepository {
    return &MemoryRepository{
        workers:   make(map[uuid.UUID]*Worker),
        phoneIndex: make(map[string]uuid.UUID),
    }
}

func (r *MemoryRepository) Create(ctx context.Context, worker *Worker) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.workers[worker.ID] = worker
    r.phoneIndex[worker.Phone] = worker.ID
    return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil, nil
    }
    copy := *w
    return &copy, nil
}

func (r *MemoryRepository) GetByPhone(ctx context.Context, phone string, tenantID uuid.UUID) (*Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    id, ok := r.phoneIndex[phone]
    if !ok {
        return nil, nil
    }
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil, nil
    }
    copy := *w
    return &copy, nil
}

func (r *MemoryRepository) Update(ctx context.Context, worker *Worker) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    // If phone changed, update index
    if old, ok := r.workers[worker.ID]; ok && old.Phone != worker.Phone {
        delete(r.phoneIndex, old.Phone)
        r.phoneIndex[worker.Phone] = worker.ID
    }
    r.workers[worker.ID] = worker
    return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    w, ok := r.workers[id]
    if !ok || w.TenantID != tenantID {
        return nil
    }
    delete(r.phoneIndex, w.Phone)
    delete(r.workers, id)
    return nil
}

func (r *MemoryRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]Worker, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var list []Worker
    for _, w := range r.workers {
        if w.TenantID == tenantID {
            list = append(list, *w)
        }
    }
    if offset >= len(list) {
        return []Worker{}, nil
    }
    end := offset + limit
    if end > len(list) {
        end = len(list)
    }
    return list[offset:end], nil
}
