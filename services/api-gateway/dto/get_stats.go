package dto

import "time"

type GetStatsResponse struct {
	ShortCode      string     `json:"short_code"`
	URL            string     `json:"url"`
	TotalClicks    int64      `json:"total_clicks"`
	UniqueVisitors int64      `json:"unique_visitors"`
	Referers       []string   `json:"referers,omitempty"`
	LastClickedAt  time.Time  `json:"last_clicked_at"`
	IsActive       bool       `json:"is_active"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}
