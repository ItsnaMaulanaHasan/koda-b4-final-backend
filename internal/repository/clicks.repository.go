package repository

import (
	"backend-koda-shortlink/internal/config"
	"backend-koda-shortlink/internal/models"
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClickRepository struct {
	db *pgxpool.Pool
}

func NewClickRepository(db *pgxpool.Pool) *ClickRepository {
	return &ClickRepository{db: db}
}

func (r *ClickRepository) Insert(ctx context.Context, data *models.Click) error {
	query := `
	INSERT INTO clicks
	(short_link_id, ip_address, referer, user_agent, country, city, device_type, browser, os)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	RETURNING short_link_id`

	_, err := r.db.Exec(ctx, query,
		data.ShortLinkID,
		data.IPAddress,
		data.Referer,
		data.UserAgent,
		data.Country,
		data.City,
		data.DeviceType,
		data.Browser,
		data.OS,
	)

	if err != nil {
		return err
	}

	var userId int
	err = r.db.QueryRow(ctx,
		`SELECT user_id FROM short_links WHERE id = $1`,
		data.ShortLinkID,
	).Scan(&userId)
	if err != nil {
		return err
	}

	config.Rdb.Del(ctx,
		"user:"+strconv.Itoa(userId)+":stats:links",
		"user:"+strconv.Itoa(userId)+":stats:visits",
		"analytics:"+strconv.Itoa(userId)+":7d",
	)

	return nil
}
