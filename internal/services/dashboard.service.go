package services

import (
	"backend-koda-shortlink/internal/repository"
	"context"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{
		repo: repo,
	}
}

type DashboardStats struct {
	TotalLinks    int
	TotalVisits   int
	AvgClickRate  float64
	Last7DaysStat any
}

func (s *DashboardService) Stats(ctx context.Context, userId int) (*DashboardStats, error) {
	totalLinks, _ := s.repo.TotalLinks(ctx, userId)
	totalVisits, _ := s.repo.TotalVisits(ctx, userId)
	last7, _ := s.repo.Last7DaysChart(ctx, userId)

	avgClickRate := 0.0
	if totalLinks > 0 {
		avgClickRate = float64(totalVisits) / float64(totalLinks)
	}

	return &DashboardStats{
		TotalLinks:    totalLinks,
		TotalVisits:   totalVisits,
		AvgClickRate:  avgClickRate,
		Last7DaysStat: last7,
	}, nil
}
