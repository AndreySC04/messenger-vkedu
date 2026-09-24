package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (repo *Repository) Create(ctx context.Context, username string, passwordHash string) (*User, error) {
	const query = `
	INSERT INTO users (username, password_hash) VALUES ($1, $2)
	RETURNING id, username, password_hash, created_at`
	var user User
	// Возвращает Row, а Scan освобождает соединение и пишет в поля
	err := repo.pool.QueryRow(ctx, query, username, passwordHash).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (repo *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	const query = `
	SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	var user User
	err := repo.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &user, nil
}

func (repo *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `
	SELECT id, username, password_hash, created_at FROM users WHERE id = $1`
	var user User
	err := repo.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}
