package models

import "time"

type ShortLink struct {
	ID            int        `json:"id" db:"id"`
	UserID        int        `json:"user_id" db:"user_id"`
	ShortCode     string     `json:"short_code" db:"short_code"`
	OriginalURL   string     `json:"original_url" db:"original_url"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	ClickCount    int        `json:"click_count" db:"click_count"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty" db:"last_clicked_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy     *int       `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy     *int       `json:"updated_by,omitempty" db:"updated_by"`
}

type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
}

type UpdateShortLinkRequest struct {
	OriginalURL *string `json:"original_url,omitempty" validate:"omitempty,url"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
