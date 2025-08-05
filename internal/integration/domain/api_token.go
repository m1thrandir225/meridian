package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type APIToken struct {
	hash string
}

func NewAPIToken(hash string) (*APIToken, error) {
	if hash == "" {
		return nil, fmt.Errorf("API Token hash cannot be empty")
	}
	return &APIToken{hash: hash}, nil
}

func GenerateLookupHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (t *APIToken) String() string {
	return t.hash
}

func (t *APIToken) Hash() string {
	return t.hash
}

func (t *APIToken) SetHash(hash string) {
	t.hash = hash
}
