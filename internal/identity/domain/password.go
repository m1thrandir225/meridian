package domain

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const passwordBcryptCost = bcrypt.DefaultCost

type passwordHash struct {
	hash string
}

func NewPasswordHash(rawPassword string) (passwordHash, error) {
	if err := validatePasswordPolicy(rawPassword); err != nil {
		return passwordHash{}, err
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), passwordBcryptCost)
	if err != nil {
		return passwordHash{}, fmt.Errorf("failed to hash password: %w", err)
	}
	return passwordHash{hash: string(hashedBytes)}, nil
}

func (ph passwordHash) Match(rawPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ph.hash), []byte(rawPassword))
	return err == nil
}

func FromHashedString(hashedPassword string) (passwordHash, error) {
	if hashedPassword == "" {
		return passwordHash{}, errors.New("hashed password string cannot be empty")
	}
	if len(hashedPassword) < 10 {
		return passwordHash{}, errors.New("hashed password string too short")
	}
	return passwordHash{hash: hashedPassword}, nil
}

func (ph passwordHash) String() string {
	return ph.hash
}

func validatePasswordPolicy(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return errors.New("password must include uppercase, lowercase, digit, and special character")
	}
	return nil
}
