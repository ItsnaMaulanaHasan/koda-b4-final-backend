package services

import (
	"backend-koda-shortlink/internal/config"
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/utils"
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/mssola/user_agent"
)

type ShortLinkService struct {
	shortLinkRepo *repository.ShortLinkRepository
	clickRepo     *repository.ClickRepository
}

func NewShortLinkService(shortLinkRepo *repository.ShortLinkRepository, clickRepo *repository.ClickRepository) *ShortLinkService {
	return &ShortLinkService{
		shortLinkRepo: shortLinkRepo,
		clickRepo:     clickRepo,
	}
}

func (s *ShortLinkService) CreateShortLink(ctx context.Context, userID int, req *models.CreateShortLinkRequest) (*models.ShortLink, error) {
	shortCode, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		return nil, err
	}

	var createdBy *int
	if userID > 0 {
		createdBy = &userID
	} else {
		createdBy = nil
	}

	link := &models.ShortLink{
		UserID:      createdBy,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}

	err = s.shortLinkRepo.Create(ctx, link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *ShortLinkService) GetUserLinks(ctx context.Context, userID int) ([]models.ShortLink, error) {
	return s.shortLinkRepo.GetAllByUserID(ctx, userID)
}

func (s *ShortLinkService) GetLinkByShortCode(ctx context.Context, shortCode string, userID int) (*models.ShortLink, error) {
	link, err := s.shortLinkRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if link.UserID != &userID {
		return nil, errors.New("unauthorized access")
	}

	return link, nil
}

func (s *ShortLinkService) UpdateShortLink(ctx context.Context, shortCode string, userID int, req *models.UpdateShortLinkRequest) (*models.ShortLink, error) {
	existing, err := s.shortLinkRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if existing.UserID != &userID {
		return nil, errors.New("unauthorized access")
	}

	link := &models.ShortLink{
		OriginalURL: "",
		UpdatedBy:   &userID,
	}

	if req.OriginalURL != nil {
		link.OriginalURL = *req.OriginalURL
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}

	err = s.shortLinkRepo.Update(ctx, shortCode, userID, link)
	if err != nil {
		return nil, err
	}

	return s.shortLinkRepo.GetByShortCode(ctx, shortCode)
}

func (s *ShortLinkService) DeleteShortLink(ctx context.Context, shortCode string, userID int) error {
	existing, err := s.shortLinkRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return err
	}
	if existing.UserID != &userID {
		return errors.New("unauthorized access")
	}

	return s.shortLinkRepo.Delete(ctx, shortCode, userID)
}

func (s *ShortLinkService) generateUniqueShortCode(ctx context.Context) (string, error) {
	maxAttempts := 5
	for range maxAttempts {
		code := utils.GenerateRandomCode(6)
		exists, err := s.shortLinkRepo.CheckShortCodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique short code")
}

func (s *ShortLinkService) ResolveShortCode(ctx context.Context, code string) (*models.ShortLink, error) {
	url, err := config.Rdb.Get(ctx, "sl:"+code).Result()
	if err == nil && url != "" {
		return &models.ShortLink{
			ShortCode:   code,
			OriginalURL: url,
			IsActive:    true,
		}, nil
	}

	link, err := s.shortLinkRepo.GetByShortCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if !link.IsActive {
		return nil, errors.New("short link inactive")
	}

	config.Rdb.Set(ctx, "sl:"+code, link.OriginalURL, 15*time.Minute)

	return link, nil
}

func (s *ShortLinkService) LogClick(code string) {
	ctx := context.Background()

	_ = s.shortLinkRepo.IncrementClick(ctx, code)
}

func (s *ShortLinkService) SaveClickAnalytics(req *http.Request, link *models.ShortLink) {
	go func() {
		ctx := context.Background()

		ua := user_agent.New(req.UserAgent())
		browser, _ := ua.Browser()

		deviceType := "desktop"
		if ua.Mobile() {
			deviceType = "mobile"
		}

		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(req.RemoteAddr)
		}

		click := &models.Click{
			ShortLinkID: link.ID,
			IPAddress:   ip,
			Referer:     req.Referer(),
			UserAgent:   req.UserAgent(),
			Country:     "",
			City:        "",
			DeviceType:  deviceType,
			Browser:     browser,
			OS:          ua.OS(),
		}

		_ = s.clickRepo.Insert(ctx, click)
	}()
}
