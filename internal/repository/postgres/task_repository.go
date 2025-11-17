package postgres

import (
	"context"
	"database/sql"
	"errors"

	"go-todo-service/internal/domain"
)

// TaskRepository persists tasks in PostgreSQL.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository constructs the repository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts a task row.
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	const query = `
		INSERT INTO tasks (id, user_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		task.ID,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)
	return err
}

// ListByUser returns tasks for a given user ordered by creation time.
func (r *TaskRepository) ListByUser(ctx context.Context, userID string) ([]domain.Task, error) {
	const query = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetByID fetches a task by identifier.
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	const query = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	task := &domain.Task{}
	if err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

// Update mutates an existing task row.
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	const query = `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete removes a task row.
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	const query = `
		DELETE FROM tasks
		WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
