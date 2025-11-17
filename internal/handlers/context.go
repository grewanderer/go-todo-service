package handlers

import "context"

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

// WithUserID stores the authenticated user id in the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts the authenticated user id if present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(userIDKey).(string)
	return value, ok && value != ""
}

// WithRequestID stores the request ID on the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext retrieves the request ID if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(requestIDKey).(string)
	return value, ok && value != ""
}
