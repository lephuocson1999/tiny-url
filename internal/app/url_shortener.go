package app

import (
	"context"
	"time"

	"tiny-url/internal/domain"
)

type URLShortenerService interface {
	ShortenURL(ctx context.Context, originalURL string, expiresAt *time.Time) (string, error)
	ResolveURL(ctx context.Context, shortCode string) (string, error)
	GetURLStats(ctx context.Context, shortCode string) (*domain.URL, error)
}

type urlShortenerService struct {
	repo  Repository
	cache Cache
	idGen IDGenerator
}

type IDGenerator interface {
	NextID(ctx context.Context) (int64, error)
}

func NewURLShortenerService(repo Repository, cache Cache, idGen IDGenerator) URLShortenerService {
	return &urlShortenerService{repo: repo, cache: cache, idGen: idGen}
}

func (s *urlShortenerService) ShortenURL(ctx context.Context, originalURL string, expiresAt *time.Time) (string, error) {
	id, err := s.idGen.NextID(ctx)
	if err != nil {
		return "", err
	}
	code := domain.EncodeBase62(id)
	if expiresAt == nil {
		t := time.Now().Add(5 * 365 * 24 * time.Hour)
		expiresAt = &t
	}
	url := &domain.URL{
		ID:          id,
		ShortCode:   code,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   *expiresAt,
	}
	if err := s.repo.Save(ctx, url); err != nil {
		return "", err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, code, url, 24*time.Hour)
	}
	return code, nil
}

func (s *urlShortenerService) ResolveURL(ctx context.Context, shortCode string) (string, error) {
	if s.cache != nil {
		if url, err := s.cache.Get(ctx, shortCode); err == nil {
			return url.OriginalURL, nil
		}
	}

	url, err := s.repo.FindByCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, shortCode, url, 24*time.Hour)
	}
	return url.OriginalURL, nil
}

func (s *urlShortenerService) GetURLStats(ctx context.Context, shortCode string) (*domain.URL, error) {
	if s.cache != nil {
		if url, err := s.cache.Get(ctx, shortCode); err == nil {
			return url, nil
		}
	}
	url, err := s.repo.FindByCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, shortCode, url, 24*time.Hour)
	}
	return url, nil
}

// Repository interface for dependency inversion
//go:generate mockgen -source=url_shortener.go -destination=../../test/mock_repo.go -package=test
