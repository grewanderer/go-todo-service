package repository

import (
	"context"

	"go-todo-service/internal/domain"
)

// UserRepository defines persistence operations for User entities.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
