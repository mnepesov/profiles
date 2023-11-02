package hasher

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Get(data []byte) (string, error) {
	hashedData, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hasing: %w", err)
	}

	return string(hashedData), nil
}

func (h *Hasher) Check(hashedData, data []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedData, data)
	if err != nil && !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, fmt.Errorf("errer checking: %w", err)
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return true, nil
}
