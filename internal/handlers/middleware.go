package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go-todo-service/pkg/jwt"
	"go-todo-service/pkg/logger"
)

// AuthMiddleware verifies JWT tokens on protected routes.
type AuthMiddleware struct {
	secret string
	log    *logger.Logger
}

// NewAuthMiddleware constructs the middleware.
func NewAuthMiddleware(secret string, log *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
		log:    log,
	}
}

// Wrap applies authentication checks to the provided handler.
func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			m.respondUnauthorized(w, r, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			m.respondUnauthorized(w, r, "invalid authorization header")
			return
		}
		token := strings.TrimSpace(parts[1])
		claims, err := jwt.ParseAndValidate(token, m.secret, time.Now())
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrExpiredToken):
				m.respondUnauthorized(w, r, "token expired")
			case errors.Is(err, jwt.ErrInvalidToken):
				m.respondUnauthorized(w, r, "invalid token")
			default:
				m.log.Error("token validation failed", map[string]any{
					"error":      err.Error(),
					"request_id": requestIDFromContextOrEmpty(r.Context()),
				})
				respondError(w, r, http.StatusUnauthorized, "invalid token")
			}
			return
		}

		ctx := WithUserID(r.Context(), claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) respondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	m.log.Info("auth failure", map[string]any{
		"reason":     message,
		"path":       r.URL.Path,
		"request_id": requestIDFromContextOrEmpty(r.Context()),
	})
	respondError(w, r, http.StatusUnauthorized, message)
}

func requestIDFromContextOrEmpty(ctx context.Context) string {
	if id, ok := RequestIDFromContext(ctx); ok {
		return id
	}
	return ""
}
