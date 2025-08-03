package token_generator

import (
	"aidanwoods.dev/go-paseto"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/domain"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"time"
)

type PasetoTokenGenerator struct {
	privateKey paseto.V4AsymmetricSecretKey
}

func NewPasetoTokenGenerator(key string) (*PasetoTokenGenerator, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid key size: %d, expected: %d", len(keyBytes), ed25519.PrivateKeySize)
	}
	v4PrivKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(keyBytes)

	if err != nil {
		return nil, fmt.Errorf("failed to create V4AsymmetricSecretKey: %w", err)
	}
	return &PasetoTokenGenerator{
		privateKey: v4PrivKey,
	}, nil
}

func (g *PasetoTokenGenerator) GenerateToken(user *domain.User, duration time.Duration) (string, *auth.TokenClaims, error) {
	if user == nil {
		return "", nil, fmt.Errorf("user is nil")
	}
	claims := auth.NewTokenClaims(
		auth.AUTH_ISSUER,
		auth.AUTH_AUDIENCE,
		auth.AUTH_SUBJECT,
		user.ID.String(),
		user.Email.String(),
		time.Now(),
		time.Now().Add(duration))
	token := paseto.NewToken()
	tokenWithClaims, err := auth.SetTokenClaims(claims, &token)
	if err != nil {
		return "", nil, err
	}

	signedToken := tokenWithClaims.V4Sign(g.privateKey, nil)

	return signedToken, &claims, nil
}

var _ services.AuthTokenGenerator = (*PasetoTokenGenerator)(nil)
