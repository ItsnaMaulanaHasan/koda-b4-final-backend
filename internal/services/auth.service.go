package services

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/utils"
	"context"
	"errors"

	"github.com/matthewhartstonge/argon2"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (string, error)
	Logout(ctx context.Context, req *models.LogoutRequest) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	err = s.userRepo.UpdateCreatedByAndUpdatedBy(ctx, user.Id)
	if err != nil {
		return nil, errors.New("failed to update user metadata")
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	isPasswordValid, err := argon2.VerifyEncoded(
		[]byte(req.Password),
		[]byte(user.Password),
	)
	if err != nil || !isPasswordValid {
		return nil, errors.New("wrong email aa or password")
	}

	refreshToken, expiresAt, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	session := &models.Session{
		UserId:       user.Id,
		RefreshToken: refreshToken,
		ExpiredAt:    expiresAt,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	sessionId, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, errors.New("failed to create session")
	}

	err = s.sessionRepo.UpdateCreatedByAndUpdatedBy(ctx, user.Id)
	if err != nil {
		return nil, errors.New("failed to update user metadata")
	}

	accessToken, err := utils.GenerateAccessToken(user.Id, sessionId)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (string, error) {
	claims, err := utils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	session, err := s.sessionRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	accessToken, err := utils.GenerateAccessToken(claims.Id, session.Id)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func (s *authService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	return s.sessionRepo.Invalidate(ctx, req.RefreshToken)
}
