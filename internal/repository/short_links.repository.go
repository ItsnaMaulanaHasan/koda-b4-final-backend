package repository

import (
	"backend-koda-shortlink/internal/config"
	"backend-koda-shortlink/internal/models"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortLinkRepository struct {
	db *pgxpool.Pool
}

func NewShortLinkRepository(db *pgxpool.Pool) *ShortLinkRepository {
	return &ShortLinkRepository{db: db}
}

func (r *ShortLinkRepository) Create(ctx context.Context, link *models.ShortLink) error {
	query := `
		INSERT INTO short_links 
		(user_id, short_code, original_url, created_by, updated_by) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at, is_active, click_count
	`

	config.Rdb.Del(ctx, "link:"+link.ShortCode+":destination")

	return r.db.QueryRow(
		ctx,
		query,
		link.UserID,
		link.ShortCode,
		link.OriginalURL,
		link.CreatedBy,
		link.UpdatedBy,
	).Scan(&link.ID, &link.CreatedAt, &link.UpdatedAt, &link.IsActive, &link.ClickCount)
}

func (r *ShortLinkRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.ShortLink, error) {
	cacheKey := "link:" + shortCode + ":destination"

	if cached, err := config.Rdb.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var link models.ShortLink
		if json.Unmarshal([]byte(cached), &link) == nil {
			return &link, nil
		}
	}

	query := `
		SELECT id, user_id, short_code, original_url, is_active, 
			   click_count, last_clicked_at, created_at, updated_at,
			   created_by, updated_by
		FROM short_links 
		WHERE short_code = $1
	`
	link := &models.ShortLink{}
	err := r.db.QueryRow(ctx, query, shortCode).Scan(
		&link.ID, &link.UserID, &link.ShortCode, &link.OriginalURL,
		&link.IsActive, &link.ClickCount, &link.LastClickedAt,
		&link.CreatedAt, &link.UpdatedAt, &link.CreatedBy, &link.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("short link not found")
		}
		return nil, err
	}

	jsonData, _ := json.Marshal(link)
	config.Rdb.Set(ctx, cacheKey, jsonData, 15*time.Minute)

	return link, nil
}

func (r *ShortLinkRepository) GetAllByUserIDWithFilter(ctx context.Context, userID, limit, offset int, search, status string) ([]models.ShortLink, int, error) {
	// Build query with filters
	baseQuery := `FROM short_links WHERE user_id = $1`
	countQuery := `SELECT COUNT(*) ` + baseQuery
	selectQuery := `
		SELECT id, user_id, short_code, original_url, is_active, 
			   click_count, last_clicked_at, created_at, updated_at,
			   created_by, updated_by
	` + baseQuery

	args := []interface{}{userID}
	argCount := 1

	if search != "" {
		argCount++
		baseQuery += ` AND (short_code ILIKE $` + strconv.Itoa(argCount) + ` OR original_url ILIKE $` + strconv.Itoa(argCount) + `)`
		args = append(args, "%"+search+"%")
		countQuery = `SELECT COUNT(*) ` + baseQuery
		selectQuery = `
			SELECT id, user_id, short_code, original_url, is_active, 
				   click_count, last_clicked_at, created_at, updated_at,
				   created_by, updated_by
		` + baseQuery
	}

	if status == "active" || status == "inactive" {
		argCount++
		isActive := status == "active"
		baseQuery += ` AND is_active = $` + strconv.Itoa(argCount)
		args = append(args, isActive)
		countQuery = `SELECT COUNT(*) ` + baseQuery
		selectQuery = `
			SELECT id, user_id, short_code, original_url, is_active, 
				   click_count, last_clicked_at, created_at, updated_at,
				   created_by, updated_by
		` + baseQuery
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	links := []models.ShortLink{}
	for rows.Next() {
		var link models.ShortLink
		err := rows.Scan(
			&link.ID, &link.UserID, &link.ShortCode, &link.OriginalURL,
			&link.IsActive, &link.ClickCount, &link.LastClickedAt,
			&link.CreatedAt, &link.UpdatedAt, &link.CreatedBy, &link.UpdatedBy,
		)
		if err != nil {
			return nil, 0, err
		}
		links = append(links, link)
	}

	return links, total, nil
}

func (r *ShortLinkRepository) Update(ctx context.Context, shortCode string, userID int, link *models.ShortLink) error {
	query := `
		UPDATE short_links 
		SET original_url = COALESCE($1, original_url),
			is_active = COALESCE($2, is_active),
			updated_by = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE short_code = $4 AND user_id = $5
		RETURNING id, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		link.OriginalURL,
		link.IsActive,
		link.UpdatedBy,
		shortCode,
		userID,
	).Scan(&link.ID, &link.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("short link not found or unauthorized")
		}
		return err
	}

	config.Rdb.Del(ctx, "link:"+link.ShortCode+":destination")

	return nil
}

func (r *ShortLinkRepository) Delete(ctx context.Context, shortCode string, userID int) error {
	query := `DELETE FROM short_links WHERE short_code = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, shortCode, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("short link not found or unauthorized")
	}

	config.Rdb.Del(ctx, "link:"+shortCode+":destination")

	return nil
}

func (r *ShortLinkRepository) CheckShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM short_links WHERE short_code = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, shortCode).Scan(&exists)
	return exists, err
}

func (r *ShortLinkRepository) IncrementClick(ctx context.Context, code string) error {
	query := `
	UPDATE short_links 
	SET click_count = click_count + 1,
		last_clicked_at = NOW()
	WHERE short_code = $1`

	_, err := r.db.Exec(ctx, query, code)
	if err != nil {
		return err
	}

	config.Rdb.Incr(ctx, "link:"+code+":clicks")

	return nil
}
