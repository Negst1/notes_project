package repository

import (
	"context"
	"database/sql"
	"time"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (int, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllRefreshTokens(ctx context.Context, userID int) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (token, expiresAt)
			  VALUES ($1, $2, $3) 
			  ON CONFLICT (token) DO UPDATE
			  SET expires_at=EXCLUDED.expires_at`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (int, error) {
	query := `SELECT user_id 
			FROM refresh_tokens 
			WHERE token=$1 AND expires_at > NOW()`

	var user_id int
	err := r.db.QueryRowContext(ctx, query).Scan(user_id)

	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokem WHERE token=$1`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) DeleteAllRefreshTokens(ctx context.Context, userID int) error {
	query := `DELETE FROM refresh_token WHERE user_id=$1`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
