package postgres

import (
	"context"
	"fmt"

	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/users"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
	"github.com/jackc/pgx/v5"
)

type UsersRepository struct {
	db *DB
}

func NewUsersRepository(db *DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Create(ctx context.Context, user *users.User) error {
	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (*users.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user users.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user users.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *UsersRepository) Update(ctx context.Context, user *users.User) error {
	query := `
		UPDATE users
		SET email = $2, password = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.Email, user.Password, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UsersRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

