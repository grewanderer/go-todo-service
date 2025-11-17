package password

import "testing"

func TestHashAndCompare(t *testing.T) {
	hash, err := Hash("supersafe")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	ok, err := Compare(hash, "supersafe")
	if err != nil {
		t.Fatalf("compare error: %v", err)
	}
	if !ok {
		t.Fatal("expected password match")
	}
}

func TestCompareMismatch(t *testing.T) {
	hash, err := Hash("supersafe")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	ok, err := Compare(hash, "wrong")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected mismatch")
	}
}

func TestHashEmptyPassword(t *testing.T) {
	if _, err := Hash(""); err == nil {
		t.Fatal("expected error for empty password")
	}
}
