package models

import "time"

type URL struct {
	ShortCode string
	URL       string
	ExpiresAt *time.Time
	IsActive  bool
}
