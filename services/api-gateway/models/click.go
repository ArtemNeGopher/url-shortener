// Package models
// Models for services
package models

import "time"

type Click struct {
	ShortCode string
	IPAdress  string
	UserAgent string
	Referer   string
	ClickedAt time.Time
}
