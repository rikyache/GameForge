package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"testsmth/internal/apperrors"
	"testsmth/internal/models"
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
