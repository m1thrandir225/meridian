package services

import "github.com/m1thrandir225/meridian/internal/integration/domain"

type TokenGenerator interface {
	Generate() (rawToken string, hashedToken domain.APIToken, err error)
	Compare(rawToken string, hashedToken domain.APIToken) error
}
