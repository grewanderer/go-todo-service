package domain

import "time"

// TaskStatus is an enum representing the state of a task.
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusDone    TaskStatus = "done"
)

// Task represents a todo entry owned by a user.
type Task struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
