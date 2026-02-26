package dto

import (
	"errors"
	"time"
)

type CreateURLRequest struct {
	URL           string  `json:"url"`
	UserID        string  `json:"user_id,omitempty"`
	ExpiresInDays *uint32 `json:"expires_in_days,omitempty"`
}

func (req *CreateURLRequest) Validate() error {
	if req.URL == "" {
		return errors.New("url is required")
	}

	return nil
}

type CreateURLResponse struct {
	ShortCode string     `json:"short_code"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
