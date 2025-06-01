package domain

import "errors"

var (
	ErrNotFound = errors.New("URL not found")
	ErrExpired  = errors.New("URL expired")
)
