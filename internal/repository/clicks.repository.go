package repository

import (
	"backend-koda-shortlink/internal/models"
	"context"

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
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

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

	return err
}
