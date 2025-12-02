package models

import "time"

type Session struct {
	Id           int        `json:"id" db:"id"`
	UserId       int        `json:"userId" db:"user_id"`
	RefreshToken string     `json:"-" db:"refresh_token"`
	LoginTime    *time.Time `json:"loginTime,omitempty" db:"-"`
	LogoutTime   *time.Time `json:"logoutTime,omitempty" db:"-"`
	ExpiredAt    time.Time  `json:"expiredAt" db:"expired_at"`
	IpAddress    string     `json:"ipAddress,omitempty" db:"-"`
	UserAgent    string     `json:"userAgent,omitempty" db:"-"`
	IsActive     bool       `json:"isActive" db:"is_active"`
}
