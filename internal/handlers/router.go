package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-todo-service/pkg/logger"
)

// NewRouter wires HTTP routes to handlers with production-grade middleware.
func NewRouter(authHandler *AuthHandler, taskHandler *TaskHandler, authMiddleware *AuthMiddleware, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(NewRecoveryMiddleware(log).Wrap)
	r.Use(NewRequestIDMiddleware().Wrap)
	r.Use(NewRequestLoggerMiddleware(log).Wrap)
	r.Use(middleware.Timeout(60 * time.Second))

	docsHandler := NewDocsHandler()
	r.Get("/docs", docsHandler.UI)
	r.Get("/docs/openapi.yaml", docsHandler.SpecYAML)
	r.Get("/docs/openapi.json", docsHandler.SpecJSON)

	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)

	r.Route("/tasks", func(sub chi.Router) {
		sub.Use(authMiddleware.Wrap)

		sub.Get("/", taskHandler.List)
		sub.Post("/", taskHandler.Create)
		sub.Get("/{id}", taskHandler.Get)
		sub.Put("/{id}", taskHandler.Update)
		sub.Delete("/{id}", taskHandler.Delete)
	})

	return r
}
