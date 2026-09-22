package repository

import (
	"context"
	"database/sql"
	"testsmth/internal/models"
)

type AuthRepository struct {
	DB *sql.DB
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*models.Player, error) {
	query := `
		SELECT id, name, balance, email, password_hash, created_at
		FROM players
		WHERE email = $1
`
	var player models.Player
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&player.ID,
		&player.Name,
		&player.Balance,
		&player.Email,
		&player.PasswordHash,
		&player.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &player, nil
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, params models.CreateUserParams) error {

	query := `
		INSERT INTO players(
		    name,
		    email,
		    password_hash
		)
		VALUES ($1, $2, $3)
`

	_, err := r.DB.ExecContext(ctx, query, params.Name, params.Email, params.PasswordHash)

	return err
}
