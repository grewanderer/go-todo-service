package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"go-todo-service/internal/domain"
	tasksvc "go-todo-service/internal/service/task"
	"go-todo-service/pkg/logger"
)

// TaskHandler exposes task management endpoints.
type TaskHandler struct {
	service *tasksvc.Service
	log     *logger.Logger
}

// NewTaskHandler constructs the handler.
func NewTaskHandler(service *tasksvc.Service, log *logger.Logger) *TaskHandler {
	return &TaskHandler{service: service, log: log}
}

// List handles GET /tasks.
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	tasks, err := h.service.ListTasks(r.Context(), userID)
	if err != nil {
		h.log.Error("list tasks failed", map[string]any{"error": err.Error()})
		respondError(w, r, http.StatusInternalServerError, "could not list tasks")
		return
	}

	response := make([]map[string]any, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, presentTask(task))
	}
	respondJSON(w, http.StatusOK, response)
}

// Create handles POST /tasks.
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid json payload")
		return
	}

	task, err := h.service.CreateTask(r.Context(), userID, payload.Title, payload.Description)
	if err != nil {
		switch {
		case errors.Is(err, tasksvc.ErrTitleRequired), errors.Is(err, tasksvc.ErrUserRequired):
			respondError(w, r, http.StatusBadRequest, err.Error())
		default:
			h.log.Error("create task failed", map[string]any{"error": err.Error()})
			respondError(w, r, http.StatusInternalServerError, "could not create task")
		}
		return
	}
	respondJSON(w, http.StatusCreated, presentTask(*task))
}

// Get handles GET /tasks/{id}.
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		respondError(w, r, http.StatusNotFound, "task not found")
		return
	}

	task, err := h.service.GetTask(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, http.StatusNotFound, "task not found")
		} else {
			h.log.Error("get task failed", map[string]any{"error": err.Error(), "task_id": id})
			respondError(w, r, http.StatusInternalServerError, "could not fetch task")
		}
		return
	}
	respondJSON(w, http.StatusOK, presentTask(*task))
}

// Update handles PUT /tasks/{id}.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		respondError(w, r, http.StatusNotFound, "task not found")
		return
	}

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid json payload")
		return
	}

	task, err := h.service.UpdateTask(r.Context(), userID, id, payload.Title, payload.Description, payload.Status)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			respondError(w, r, http.StatusNotFound, "task not found")
		case errors.Is(err, tasksvc.ErrInvalidStatus):
			respondError(w, r, http.StatusBadRequest, err.Error())
		default:
			h.log.Error("update task failed", map[string]any{"error": err.Error(), "task_id": id})
			respondError(w, r, http.StatusInternalServerError, "could not update task")
		}
		return
	}
	respondJSON(w, http.StatusOK, presentTask(*task))
}

// Delete handles DELETE /tasks/{id}.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		respondError(w, r, http.StatusNotFound, "task not found")
		return
	}

	if err := h.service.DeleteTask(r.Context(), userID, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, http.StatusNotFound, "task not found")
		} else {
			h.log.Error("delete task failed", map[string]any{"error": err.Error(), "task_id": id})
			respondError(w, r, http.StatusInternalServerError, "could not delete task")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func presentTask(task domain.Task) map[string]any {
	return map[string]any{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"user_id":     task.UserID,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	}
}
