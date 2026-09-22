package service

import (
	"context"
	"testsmth/internal/apperrors"
	"testsmth/internal/auth"
	"testsmth/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, params models.CreateUserParams) error
	GetByEmail(ctx context.Context, email string) (*models.Player, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) error {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return apperrors.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	params := models.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	return s.repo.CreateUser(ctx, params)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	if req.Email == "" || req.Password == "" {
		return "", apperrors.ErrInvalidInput
	}

	player, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(req.Password)); err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(player.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
