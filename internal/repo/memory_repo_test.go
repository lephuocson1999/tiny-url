package repo

import (
	"context"
	"testing"
	"time"

	"tiny-url/internal/domain"
)

func TestMemoryRepo_SaveAndFindByCode(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()
	url := &domain.URL{
		ID:          1,
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	if err := repo.Save(ctx, url); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := repo.FindByCode(ctx, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OriginalURL != url.OriginalURL {
		t.Errorf("got %q, want %q", got.OriginalURL, url.OriginalURL)
	}
}

func TestMemoryRepo_NotFound(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()
	_, err := repo.FindByCode(ctx, "notfound")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryRepo_Expired(t *testing.T) {
	repo := NewMemoryRepo()
	ctx := context.Background()
	url := &domain.URL{
		ID:          2,
		ShortCode:   "expired",
		OriginalURL: "https://expired.com",
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		ExpiresAt:   time.Now().Add(-24 * time.Hour),
	}
	repo.Save(ctx, url)
	_, err := repo.FindByCode(ctx, "expired")
	if err != ErrExpired {
		t.Errorf("expected ErrExpired, got %v", err)
	}
}
