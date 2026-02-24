package models

import "time"

type DayStats struct {
	ShortCode      string
	Day            time.Time
	TotalClicks    int64
	UniqueVisitors int64
	Referers       []string
}
