package handlers

import (
	"net/http"
	"runtime/debug"
	"time"

	"go-todo-service/pkg/logger"
	"go-todo-service/pkg/uuid"
)

// RequestIDMiddleware injects a request ID into the context and response headers.
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware constructs the middleware.
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// Wrap applies the request ID behaviour.
func (m *RequestIDMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.NewString()
		if err != nil {
			http.Error(w, "unable to generate request id", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Request-ID", id)
		ctx := WithRequestID(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestLoggerMiddleware logs request metadata and status codes.
type RequestLoggerMiddleware struct {
	log *logger.Logger
}

// NewRequestLoggerMiddleware constructs the middleware.
func NewRequestLoggerMiddleware(log *logger.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{log: log}
}

// Wrap applies logging around the request lifecycle.
func (m *RequestLoggerMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)

		fields := map[string]any{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     recorder.status,
			"durationMs": time.Since(start).Milliseconds(),
		}
		if requestID, ok := RequestIDFromContext(r.Context()); ok {
			fields["request_id"] = requestID
		}
		m.log.Info("http_request", fields)
	})
}

// RecoveryMiddleware captures panics and returns a 500 response.
type RecoveryMiddleware struct {
	log *logger.Logger
}

// NewRecoveryMiddleware constructs the middleware.
func NewRecoveryMiddleware(log *logger.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{log: log}
}

// Wrap recovers from panics in downstream handlers.
func (m *RecoveryMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				m.log.Error("panic recovered", map[string]any{
					"error": rec,
					"stack": string(debug.Stack()),
				})
				respondError(w, r, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
