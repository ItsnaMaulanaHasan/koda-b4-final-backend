package models

import (
	"backend-koda-shortlink/internal/database"
	"context"
	"errors"
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
	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		err = errors.New("failed to start database transaction")
		return err
	}
	defer tx.Rollback(ctx)

	var sessionId int

	query := `
		INSERT INTO sessions (user_id, refresh_token, expired_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	err = tx.QueryRow(ctx, query, userId, refreshToken, expiredAt, ipAddress, userAgent).Scan(&sessionId)
	if err != nil {
		err = errors.New("internal server error while inserting new user")
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE sessions SET created_by = $1, updated_by = $1 WHERE id = $2`, userId, sessionId)
	if err != nil {
		err = errors.New("internal server error while update created_by and updated_by")
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = errors.New("failed to commit transaction")
		return err
	}

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

func InvalidateSession(userId int, refreshToken string) error {
	query := `
		UPDATE sessions 
		SET is_active = false, logout_time = NOW(), updated_at = NOW(), updated_by = $2
		WHERE refresh_token = $1
	`

	_, err := database.DB.Exec(context.Background(), query, refreshToken, userId)
	return err
}

func InvalidateAllUserSessions(userId int) error {
	query := `
		UPDATE sessions 
		SET is_active = false, logout_time = NOW(), updated_at = NOW(), updated_by = $1
		WHERE user_id = $1 AND is_active = true
	`

	_, err := database.DB.Exec(context.Background(), query, userId)
	return err
}
