package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-todo-service/internal/domain"
	"go-todo-service/internal/repository"
	"go-todo-service/pkg/jwt"
	"go-todo-service/pkg/password"
	"go-todo-service/pkg/uuid"
)

var (
	// ErrInvalidEmail indicates a malformed or empty email address.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrWeakPassword indicates insufficient password length.
	ErrWeakPassword = errors.New("password must be at least 6 characters")
)

// Service provides authentication use-cases.
type Service struct {
	users     repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
	now       func() time.Time
}

// New constructs a Service instance.
func New(users repository.UserRepository, jwtSecret string, tokenTTL time.Duration) *Service {
	return &Service{
		users:     users,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
		now:       time.Now,
	}
}

// WithNow overrides the time source (primarily for testing).
func (s *Service) WithNow(fn func() time.Time) {
	if fn != nil {
		s.now = fn
	}
}

// Signup registers a new user.
func (s *Service) Signup(ctx context.Context, email, plainPassword string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	if len(plainPassword) < 6 {
		return nil, ErrWeakPassword
	}

	if existing, err := s.users.GetByEmail(ctx, email); err == nil && existing != nil {
		return nil, domain.ErrConflict
	} else if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	hashed, err := password.Hash(plainPassword)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewString()
	if err != nil {
		return nil, err
	}

	now := s.now().UTC()
	user := &domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: hashed,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Login verifies credentials and returns a signed JWT.
func (s *Service) Login(ctx context.Context, email, plainPassword string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || plainPassword == "" {
		return "", domain.ErrInvalidCredentials
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	ok, err := password.Compare(user.PasswordHash, plainPassword)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", domain.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, s.tokenTTL, s.now())
	if err != nil {
		return "", err
	}

	return token, nil
}
