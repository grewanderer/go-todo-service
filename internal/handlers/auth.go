package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go-todo-service/internal/domain"
	"go-todo-service/internal/service/auth"
	"go-todo-service/pkg/logger"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	service *auth.Service
	log     *logger.Logger
}

// NewAuthHandler constructs the handler.
func NewAuthHandler(service *auth.Service, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		log:     log,
	}
}

// Signup handles POST /auth/signup.
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid json payload")
		return
	}

	user, err := h.service.Signup(r.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidEmail), errors.Is(err, auth.ErrWeakPassword):
			respondError(w, r, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrConflict):
			respondError(w, r, http.StatusConflict, "email already registered")
		default:
			h.log.Error("signup failed", map[string]any{"error": err.Error()})
			respondError(w, r, http.StatusInternalServerError, "could not create user")
		}
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"message": "signup successful",
		"user": map[string]any{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, r, http.StatusBadRequest, "invalid json payload")
		return
	}

	token, err := h.service.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			respondError(w, r, http.StatusUnauthorized, "invalid credentials")
		} else {
			h.log.Error("login failed", map[string]any{"error": err.Error()})
			respondError(w, r, http.StatusInternalServerError, "could not login user")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
