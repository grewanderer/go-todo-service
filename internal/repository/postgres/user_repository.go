package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"go-todo-service/internal/domain"
)

// UserRepository stores and retrieves users from PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository constructs the repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create persists a new user row.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return err
	}
	return nil
}

// GetByEmail fetches a user via email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, strings.ToLower(email))
	user := &domain.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetByID fetches a user via id.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &domain.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}
