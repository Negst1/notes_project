package repository

import (
	"context"
	"database/sql"
	"idm/database/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, id int, user *models.User) error
	Delete(ctx context.Context, id int) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, last_name, email) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		user.Username, user.LastName, user.Email,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, last_name, email, created_at, updated_at 
              FROM users 
              WHERE id = $1 AND deleted_at IS NULL`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, last_name, email, created_at, updated_at 
              FROM users 
              ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.LastName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id int, user *models.User) error {
	query := `UPDATE users 
              SET username = $1,
                  email = $2,
                  first_name = $3,
                  last_name = $4,
                  password_hash = $5,
                  is_active = $6,
                  is_verified = $7,
                  updated_at = CURRENT_TIMESTAMP
              WHERE id = $8`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash,
		user.IsActive,
		user.IsVerified,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, last_name, created_at, updated_at
			FROM users 
			WHERE email=$1`
	var user models.User
	err := r.db.QueryRowContext(ctx, query).Scan(
		&user.ID,
		&user.Username,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
