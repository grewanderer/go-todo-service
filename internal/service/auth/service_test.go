package auth_test

import (
	"context"
	"testing"
	"time"

	"go-todo-service/internal/domain"
	authsvc "go-todo-service/internal/service/auth"
	"go-todo-service/pkg/password"
)

type fakeUserRepo struct {
	users map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *fakeUserRepo) Create(ctx context.Context, user *domain.User) error {
	if _, exists := r.users[user.Email]; exists {
		return domain.ErrConflict
	}
	u := *user
	r.users[user.Email] = &u
	return nil
}

func (r *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if user, ok := r.users[email]; ok {
		u := *user
		return &u, nil
	}
	return nil, domain.ErrNotFound
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			u := *user
			return &u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func TestSignupSuccess(t *testing.T) {
	repo := newFakeUserRepo()
	service := authsvc.New(repo, "secret", 15*time.Minute)
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	service.WithNow(func() time.Time { return fixed })

	user, err := service.Signup(context.Background(), "user@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("expected email stored, got %s", user.Email)
	}
	if user.CreatedAt.IsZero() || !user.CreatedAt.Equal(fixed) {
		t.Fatalf("expected timestamp to equal fixed time, got %v", user.CreatedAt)
	}
}

func TestSignupDuplicateEmail(t *testing.T) {
	repo := newFakeUserRepo()
	service := authsvc.New(repo, "secret", 15*time.Minute)
	service.WithNow(func() time.Time { return time.Now() })

	_, err := service.Signup(context.Background(), "dup@example.com", "password")
	if err != nil {
		t.Fatalf("unexpected error on first signup: %v", err)
	}
	_, err = service.Signup(context.Background(), "dup@example.com", "password")
	if err == nil {
		t.Fatal("expected error on duplicate signup")
	}
	if err != domain.ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestLogin(t *testing.T) {
	repo := newFakeUserRepo()
	hashed, _ := password.Hash("password")
	repo.users["user@example.com"] = &domain.User{
		ID:           "abc",
		Email:        "user@example.com",
		PasswordHash: hashed,
	}

	service := authsvc.New(repo, "secret", 15*time.Minute)
	service.WithNow(func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) })

	token, err := service.Login(context.Background(), "user@example.com", "password")
	if err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	repo := newFakeUserRepo()
	hashed, _ := password.Hash("password")
	repo.users["user@example.com"] = &domain.User{
		ID:           "abc",
		Email:        "user@example.com",
		PasswordHash: hashed,
	}

	service := authsvc.New(repo, "secret", 15*time.Minute)
	_, err := service.Login(context.Background(), "user@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
