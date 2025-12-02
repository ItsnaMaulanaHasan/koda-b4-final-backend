package repository

import (
	"backend-koda-shortlink/internal/database"
	"backend-koda-shortlink/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) (int, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	CheckActive(ctx context.Context, sessionId int) (bool, error)
	Invalidate(ctx context.Context, refreshToken string) error
	InvalidateById(ctx context.Context, sessionId int) error
	InvalidateAllByUserId(ctx context.Context, userId int) error
	UpdateCreatedByAndUpdatedBy(ctx context.Context, userId int) error
}

type sessionRepository struct{}

func NewSessionRepository() SessionRepository {
	return &sessionRepository{}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) (int, error) {
	var sessionId int
	query := `
		INSERT INTO sessions (user_id, refresh_token, expired_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := database.DB.QueryRow(
		ctx,
		query,
		session.UserId,
		session.RefreshToken,
		session.ExpiredAt,
		session.IpAddress,
		session.UserAgent,
	).Scan(&sessionId)

	return sessionId, err
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, expired_at, is_active
		FROM sessions
		WHERE refresh_token = $1 AND is_active = true AND expired_at > NOW()
	`

	rows, err := database.DB.Query(ctx, query, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	session, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("session not found or expired")
		}
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) CheckActive(ctx context.Context, sessionId int) (bool, error) {
	var isActive bool
	query := `SELECT is_active FROM sessions WHERE id = $1 AND expired_at > NOW()`

	err := database.DB.QueryRow(ctx, query, sessionId).Scan(&isActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return isActive, nil
}

func (r *sessionRepository) Invalidate(ctx context.Context, refreshToken string) error {
	query := `
		UPDATE sessions
		SET is_active = false, logout_time = NOW(), updated_at = NOW()
		WHERE refresh_token = $1
	`

	_, err := database.DB.Exec(ctx, query, refreshToken)
	return err
}

func (r *sessionRepository) InvalidateById(ctx context.Context, sessionId int) error {
	query := `
		UPDATE sessions
		SET is_active = false, logout_time = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.DB.Exec(ctx, query, sessionId)
	return err
}

func (r *sessionRepository) InvalidateAllByUserId(ctx context.Context, userId int) error {
	query := `
		UPDATE sessions
		SET is_active = false, logout_time = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND is_active = true
	`

	_, err := database.DB.Exec(ctx, query, userId)
	return err
}

func (r *sessionRepository) UpdateCreatedByAndUpdatedBy(ctx context.Context, userId int) error {
	query := `UPDATE sessions SET created_by = $1, updated_by = $1 WHERE id = $1`
	_, err := database.DB.Exec(ctx, query, userId)
	return err
}
