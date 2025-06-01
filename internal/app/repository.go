package app

import (
	"context"
	"tiny-url/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, url *domain.URL) error
	FindByCode(ctx context.Context, code string) (*domain.URL, error)
}
