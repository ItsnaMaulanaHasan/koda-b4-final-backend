package repository

import (
	"backend-koda-shortlink/internal/config"
	"context"
	"encoding/json"
	"strconv"
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
	key := "user:" + strconv.Itoa(userId) + ":stats:links"

	if cached, err := config.Rdb.Get(ctx, key).Result(); err == nil {
		val, _ := strconv.Atoi(cached)
		return val, nil
	}

	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM short_links WHERE user_id = $1`, userId)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	config.Rdb.Set(ctx, key, total, 5*time.Minute)

	return total, nil
}

func (r *DashboardRepository) TotalVisits(ctx context.Context, userId int) (int, error) {
	key := "user:" + strconv.Itoa(userId) + ":stats:visits"

	if cached, err := config.Rdb.Get(ctx, key).Result(); err == nil {
		val, _ := strconv.Atoi(cached)
		return val, nil
	}

	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*) 
         FROM clicks c
         JOIN short_links sl ON sl.id = c.short_link_id
         WHERE sl.user_id = $1`, userId)

	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	config.Rdb.Set(ctx, key, total, 5*time.Minute)
	return total, nil
}

func (r *DashboardRepository) Last7DaysChart(ctx context.Context, userId int) ([]DailyVisit, error) {
	key := "analytics:" + strconv.Itoa(userId) + ":7d"

	if cached, err := config.Rdb.Get(ctx, key).Result(); err == nil && cached != "" {
		var result []DailyVisit
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result, nil
		}
	}

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
		if err := rows.Scan(&d.Day, &d.Count); err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	jsonData, _ := json.Marshal(result)
	config.Rdb.Set(ctx, key, jsonData, 1*time.Minute)

	return result, nil
}
