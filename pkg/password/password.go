package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hash returns a bcrypt hash for the provided plaintext password.
func Hash(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("empty password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Compare verifies that the provided password matches the stored bcrypt hash.
func Compare(hashed, plain string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
