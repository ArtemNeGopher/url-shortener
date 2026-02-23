package models

import "time"

type Stats struct {
	TotalClicks    int64
	UniqueVisitors int64
	LastClickedAt  *time.Time
}
