package models

import "time"

type ShortLink struct {
	ID            int        `json:"id" db:"id"`
	UserID        *int       `json:"userId" db:"user_id"`
	ShortCode     string     `json:"shortCode" db:"short_code"`
	OriginalURL   string     `json:"originalUrl" db:"original_url"`
	IsActive      bool       `json:"isActive" db:"is_active"`
	ClickCount    int        `json:"clickCount" db:"click_count"`
	LastClickedAt *time.Time `json:"lastClicked_at,omitempty" db:"last_clicked_at"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	CreatedBy     *int       `json:"createdBy,omitempty" db:"created_by"`
	UpdatedBy     *int       `json:"updatedBy,omitempty" db:"updated_by"`
}

type ShortLinkResponse struct {
	ShortCode   string `json:"shortCode"`
	OriginalUrl string `json:"originalUrl"`
	ShortUrl    string `json:"shortUrl"`
}

type CreateShortLinkRequest struct {
	OriginalURL string `json:"originalUrl" validate:"required,url"`
}

type UpdateShortLinkRequest struct {
	OriginalURL *string `json:"originalUrl,omitempty" validate:"omitempty,url"`
	IsActive    *bool   `json:"isActive,omitempty"`
}
