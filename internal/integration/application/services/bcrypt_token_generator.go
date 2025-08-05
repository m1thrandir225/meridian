package services

import (
	"errors"
	"fmt"

	"github.com/m1thrandir225/meridian/internal/integration/domain"
	"github.com/m1thrandir225/meridian/pkg/common"
	"golang.org/x/crypto/bcrypt"
)

type BcryptTokenGenerator struct{}

func NewBcryptTokenGenerator() *BcryptTokenGenerator { return &BcryptTokenGenerator{} }

func (g *BcryptTokenGenerator) Generate() (rawToken string, hashedToken *domain.APIToken, err error) {
	rawToken, err = common.GenerateRandomString(domain.API_TOKEN_BYTES)
	if err != nil {
		return "", nil, err
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(rawToken), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, fmt.Errorf("bcrypt hashing failed: %w", err)
	}
	hashedToken, err = domain.NewAPIToken(string(hashedBytes))
	return rawToken, hashedToken, err
}

func (g *BcryptTokenGenerator) Compare(rawToken string, hashedToken *domain.APIToken) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedToken.Hash()), []byte(rawToken)); err != nil {
		return errors.New("token mismatch")
	}
	return nil
}
