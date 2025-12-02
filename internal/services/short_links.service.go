package services

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/repository"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

type ShortLinkService struct {
	repo *repository.ShortLinkRepository
}

func NewShortLinkService(repo *repository.ShortLinkRepository) *ShortLinkService {
	return &ShortLinkService{repo: repo}
}

func (s *ShortLinkService) CreateShortLink(ctx context.Context, userID int, req *models.CreateShortLinkRequest) (*models.ShortLink, error) {
	shortCode, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		return nil, err
	}

	link := &models.ShortLink{
		UserID:      userID,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		CreatedBy:   &userID,
		UpdatedBy:   &userID,
	}

	err = s.repo.Create(ctx, link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *ShortLinkService) GetUserLinks(ctx context.Context, userID int) ([]models.ShortLink, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *ShortLinkService) GetLinkByShortCode(ctx context.Context, shortCode string, userID int) (*models.ShortLink, error) {
	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if link.UserID != userID {
		return nil, errors.New("unauthorized access")
	}

	return link, nil
}

func (s *ShortLinkService) UpdateShortLink(ctx context.Context, shortCode string, userID int, req *models.UpdateShortLinkRequest) (*models.ShortLink, error) {
	existing, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
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

	err = s.repo.Update(ctx, shortCode, userID, link)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByShortCode(ctx, shortCode)
}

func (s *ShortLinkService) DeleteShortLink(ctx context.Context, shortCode string, userID int) error {
	existing, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return err
	}
	if existing.UserID != userID {
		return errors.New("unauthorized access")
	}

	return s.repo.Delete(ctx, shortCode, userID)
}

func (s *ShortLinkService) generateUniqueShortCode(ctx context.Context) (string, error) {
	maxAttempts := 5
	for range maxAttempts {
		code := generateRandomCode(6)
		exists, err := s.repo.CheckShortCodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique short code")
}

func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	code := base64.URLEncoding.EncodeToString(bytes)
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")
	if len(code) > length {
		code = code[:length]
	}
	return code
}
