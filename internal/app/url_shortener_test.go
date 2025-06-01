package app

import (
	"context"
	"testing"

	"tiny-url/internal/repo"
)

// in-memory ID generator for tests
type memoryIDGen struct {
	next int64
}

func (g *memoryIDGen) NextID(ctx context.Context) (int64, error) {
	g.next++
	return g.next, nil
}

func TestURLShortenerService_Success(t *testing.T) {
	repository := repo.NewMemoryRepo()
	idGen := &memoryIDGen{}
	service := NewURLShortenerService(repository, nil, idGen)
	ctx := context.Background()
	original := "https://golang.org"
	code, err := service.ShortenURL(ctx, original, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved, err := service.ResolveURL(ctx, code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != original {
		t.Errorf("got %q, want %q", resolved, original)
	}
	stats, err := service.GetURLStats(ctx, code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.OriginalURL != original {
		t.Errorf("got %q, want %q", stats.OriginalURL, original)
	}
}

func TestURLShortenerService_NotFound(t *testing.T) {
	repository := repo.NewMemoryRepo()
	idGen := &memoryIDGen{}
	service := NewURLShortenerService(repository, nil, idGen)
	ctx := context.Background()
	_, err := service.ResolveURL(ctx, "notfound")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
