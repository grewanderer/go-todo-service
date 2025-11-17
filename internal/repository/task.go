package repository

import (
	"context"

	"go-todo-service/internal/domain"
)

// TaskRepository defines persistence operations for Task entities.
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	ListByUser(ctx context.Context, userID string) ([]domain.Task, error)
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
}
