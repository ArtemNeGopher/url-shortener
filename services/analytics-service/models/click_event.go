package models

import "time"

type ClickEvent struct {
	ShortCode string
	IPAddress string
	UserAgent string
	Timestamp time.Time
}
