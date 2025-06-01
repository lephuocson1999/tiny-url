package repo

import (
	"context"
	"sync"
	"time"

	"tiny-url/internal/domain"
)

type MemoryRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.URL
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: make(map[string]*domain.URL),
	}
}

func (r *MemoryRepo) Save(ctx context.Context, url *domain.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[url.ShortCode] = url
	return nil
}

func (r *MemoryRepo) FindByCode(ctx context.Context, code string) (*domain.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	url, ok := r.data[code]
	if !ok {
		return nil, ErrNotFound
	}
	if url.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpired
	}
	return url, nil
}

var (
	ErrNotFound = domain.ErrNotFound
	ErrExpired  = domain.ErrExpired
)
