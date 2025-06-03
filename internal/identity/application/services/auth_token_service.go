package services

import (
	"time"

	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

type AuthTokenGenerator interface {
	GenerateToken(user *domain.User, duration time.Duration) (tokenString string, claims *auth.TokenClaims, err error)
}
