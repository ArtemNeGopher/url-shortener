// Package models
package models

import "time"

type ClickEvent struct {
	ShortCode string
	IPAddress string
	UserAgent string
	Referer   string
	Timestamp time.Time
}
