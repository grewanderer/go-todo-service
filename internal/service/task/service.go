package task

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-todo-service/internal/domain"
	"go-todo-service/internal/repository"
	"go-todo-service/pkg/uuid"
)

var (
	// ErrUserRequired indicates a missing user identifier.
	ErrUserRequired = errors.New("user id required")
	// ErrTitleRequired indicates a missing task title.
	ErrTitleRequired = errors.New("title is required")
	// ErrInvalidStatus indicates status is outside supported values.
	ErrInvalidStatus = errors.New("invalid status")
)

// Service encapsulates task management use cases.
type Service struct {
	tasks repository.TaskRepository
	now   func() time.Time
}

// New constructs a task service.
func New(tasks repository.TaskRepository) *Service {
	return &Service{
		tasks: tasks,
		now:   time.Now,
	}
}

// WithNow overrides the time source (testing).
func (s *Service) WithNow(fn func() time.Time) {
	if fn != nil {
		s.now = fn
	}
}

// CreateTask stores a new task for the provided user.
func (s *Service) CreateTask(ctx context.Context, userID, title, description string) (*domain.Task, error) {
	title = strings.TrimSpace(title)
	if userID == "" {
		return nil, ErrUserRequired
	}
	if title == "" {
		return nil, ErrTitleRequired
	}

	id, err := uuid.NewString()
	if err != nil {
		return nil, err
	}

	now := s.now().UTC()
	task := &domain.Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: strings.TrimSpace(description),
		Status:      domain.TaskStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.tasks.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// ListTasks returns all tasks for the provided user.
func (s *Service) ListTasks(ctx context.Context, userID string) ([]domain.Task, error) {
	if userID == "" {
		return nil, ErrUserRequired
	}
	return s.tasks.ListByUser(ctx, userID)
}

// GetTask fetches a single task ensuring the owner matches.
func (s *Service) GetTask(ctx context.Context, userID, id string) (*domain.Task, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.UserID != userID {
		return nil, domain.ErrNotFound
	}
	return task, nil
}

// UpdateTask updates mutable fields of a task.
func (s *Service) UpdateTask(ctx context.Context, userID, id, title, description, status string) (*domain.Task, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task.UserID != userID {
		return nil, domain.ErrNotFound
	}

	if title = strings.TrimSpace(title); title != "" {
		task.Title = title
	}
	task.Description = strings.TrimSpace(description)

	if status != "" {
		switch domain.TaskStatus(status) {
		case domain.TaskStatusPending, domain.TaskStatusDone:
			task.Status = domain.TaskStatus(status)
		default:
			return nil, ErrInvalidStatus
		}
	}

	task.UpdatedAt = s.now().UTC()

	if err := s.tasks.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// DeleteTask removes a task owned by the user.
func (s *Service) DeleteTask(ctx context.Context, userID, id string) error {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return domain.ErrNotFound
	}
	return s.tasks.Delete(ctx, id)
}
