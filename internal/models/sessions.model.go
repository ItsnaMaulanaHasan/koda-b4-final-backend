package models

import (
	"backend-koda-shortlink/internal/database"
	"context"
	"time"
)

type Session struct {
	Id           int        `json:"id"`
	UserId       int        `json:"userId"`
	RefreshToken string     `json:"refreshToken"`
	LoginTime    time.Time  `json:"loginTime"`
	LogoutTime   *time.Time `json:"logoutTime"`
	ExpiredAt    time.Time  `json:"expiredAt"`
	IpAddress    string     `json:"ipAddress"`
	UserAgent    string     `json:"userAgent"`
	IsActive     bool       `json:"isActive"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func CreateSession(userId int, refreshToken string, expiredAt time.Time, ipAddress, userAgent string) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, expired_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := database.DB.Exec(context.Background(), query, userId, refreshToken, expiredAt, ipAddress, userAgent)
	return err
}

func GetSessionByRefreshToken(refreshToken string) (*Session, error) {
	var session Session
	query := `
		SELECT id, user_id, refresh_token, expired_at, is_active 
		FROM sessions 
		WHERE refresh_token = $1 AND is_active = true AND expired_at > NOW()
	`

	err := database.DB.QueryRow(context.Background(), query, refreshToken).Scan(
		&session.Id,
		&session.UserId,
		&session.RefreshToken,
		&session.ExpiredAt,
		&session.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func InvalidateSession(refreshToken string) error {
	query := `
		UPDATE sessions 
		SET is_active = false, logout_time = NOW(), updated_at = NOW()
		WHERE refresh_token = $1
	`

	_, err := database.DB.Exec(context.Background(), query, refreshToken)
	return err
}

func InvalidateAllUserSessions(userId int) error {
	query := `
		UPDATE sessions 
		SET is_active = false, logout_time = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND is_active = true
	`

	_, err := database.DB.Exec(context.Background(), query, userId)
	return err
}
