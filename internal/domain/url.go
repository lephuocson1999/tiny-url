package domain

import "time"

type URL struct {
	ID           int64
	ShortCode    string
	OriginalURL  string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	AccessCount  int64
	LastAccessed *time.Time
}
