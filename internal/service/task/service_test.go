package task_test

import (
	"context"
	"testing"
	"time"

	"go-todo-service/internal/domain"
	tasksvc "go-todo-service/internal/service/task"
)

type fakeTaskRepo struct {
	tasks map[string]domain.Task
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{
		tasks: make(map[string]domain.Task),
	}
}

func (r *fakeTaskRepo) Create(ctx context.Context, task *domain.Task) error {
	r.tasks[task.ID] = *task
	return nil
}

func (r *fakeTaskRepo) ListByUser(ctx context.Context, userID string) ([]domain.Task, error) {
	var out []domain.Task
	for _, task := range r.tasks {
		if task.UserID == userID {
			out = append(out, task)
		}
	}
	return out, nil
}

func (r *fakeTaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copy := task
	return &copy, nil
}

func (r *fakeTaskRepo) Update(ctx context.Context, task *domain.Task) error {
	if _, ok := r.tasks[task.ID]; !ok {
		return domain.ErrNotFound
	}
	r.tasks[task.ID] = *task
	return nil
}

func (r *fakeTaskRepo) Delete(ctx context.Context, id string) error {
	if _, ok := r.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func TestCreateTask(t *testing.T) {
	repo := newFakeTaskRepo()
	service := tasksvc.New(repo)
	fixed := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
	service.WithNow(func() time.Time { return fixed })

	task, err := service.CreateTask(context.Background(), "user-1", "Title", "Desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Status != domain.TaskStatusPending {
		t.Fatalf("expected pending status, got %s", task.Status)
	}
	if !task.CreatedAt.Equal(fixed) || !task.UpdatedAt.Equal(fixed) {
		t.Fatalf("timestamps not set to fixed time")
	}
}

func TestCreateTaskMissingTitle(t *testing.T) {
	service := tasksvc.New(newFakeTaskRepo())
	if _, err := service.CreateTask(context.Background(), "user-1", "", "desc"); err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestUpdateTaskInvalidStatus(t *testing.T) {
	repo := newFakeTaskRepo()
	task := domain.Task{
		ID:          "task-1",
		UserID:      "user-1",
		Title:       "Title",
		Description: "Desc",
		Status:      domain.TaskStatusPending,
	}
	repo.tasks[task.ID] = task

	service := tasksvc.New(repo)
	if _, err := service.UpdateTask(context.Background(), "user-1", task.ID, "", "", "invalid"); err != tasksvc.ErrInvalidStatus {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestGetTaskOwnership(t *testing.T) {
	repo := newFakeTaskRepo()
	task := domain.Task{
		ID:     "task-1",
		UserID: "user-1",
		Title:  "Title",
		Status: domain.TaskStatusPending,
	}
	repo.tasks[task.ID] = task

	service := tasksvc.New(repo)
	if _, err := service.GetTask(context.Background(), "user-2", task.ID); err != domain.ErrNotFound {
		t.Fatalf("expected ErrNotFound when accessing another user's task, got %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	repo := newFakeTaskRepo()
	task := domain.Task{
		ID:     "task-1",
		UserID: "user-1",
		Title:  "Title",
		Status: domain.TaskStatusPending,
	}
	repo.tasks[task.ID] = task

	service := tasksvc.New(repo)
	if err := service.DeleteTask(context.Background(), "user-1", task.ID); err != nil {
		t.Fatalf("expected no error deleting own task: %v", err)
	}
	if _, exists := repo.tasks[task.ID]; exists {
		t.Fatal("task should be deleted")
	}
}
