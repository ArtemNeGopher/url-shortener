package models

type DayStats struct {
	ShortCode      string
	Date           string
	TotalClicks    int64
	UniqueVisitors int64
	Referers       []string
}
