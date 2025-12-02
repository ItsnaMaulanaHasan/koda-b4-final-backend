package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

type DailyVisit struct {
	Day   time.Time `json:"day"`
	Count int       `json:"count"`
}

func (r *DashboardRepository) TotalLinks(ctx context.Context, userId int) (int, error) {
	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM short_links WHERE user_id = $1`, userId)
	var total int
	return total, row.Scan(&total)
}

func (r *DashboardRepository) TotalVisits(ctx context.Context, userId int) (int, error) {
	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*) 
		 FROM clicks c
		 JOIN short_links sl ON sl.id = c.short_link_id
		 WHERE sl.user_id = $1`, userId)
	var total int
	return total, row.Scan(&total)
}

func (r *DashboardRepository) Last7DaysChart(ctx context.Context, userId int) ([]DailyVisit, error) {
	query := `
		SELECT DATE(clicked_at) AS day, COUNT(*) 
		FROM clicks c
		JOIN short_links sl ON sl.id = c.short_link_id
		WHERE sl.user_id = $1
		AND clicked_at >= NOW() - INTERVAL '7 days'
		GROUP BY day
		ORDER BY day ASC`

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DailyVisit
	for rows.Next() {
		var d DailyVisit
		err := rows.Scan(&d.Day, &d.Count)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, nil
}
