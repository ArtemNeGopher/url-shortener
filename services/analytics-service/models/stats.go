package models

import "time"

type Stats struct {
	ShortCode      string
	TotalClicks    int64
	UniqueVisitors int64
	LastClickedAt  *time.Time
	Referers       []string
}
