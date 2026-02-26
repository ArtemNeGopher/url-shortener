package dto

type GetDayStatsResponse struct {
	ShortCode      string   `json:"short_code"`
	URL            string   `json:"url"`
	Date           string   `json:"date"`
	TotalClicks    int64    `json:"total_clicks"`
	UniqueVisitors int64    `json:"unique_visitors"`
	Referers       []string `json:"referers,omitempty"`
}
