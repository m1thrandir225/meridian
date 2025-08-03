package auth

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
)

type TokenVerifier interface {
	Verify(token string) (*TokenClaims, error)
}

type PasetoTokenVerifier struct {
	publicKey paseto.V4AsymmetricPublicKey
}

func NewPasetoTokenVerifier() (*PasetoTokenVerifier, error) {
	pubKeyBytes, err := GetPublicKey()
	if err != nil {
		return nil, err
	}

	v4PubKey, err := paseto.NewV4AsymmetricPublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return &PasetoTokenVerifier{
		publicKey: v4PubKey,
	}, nil
}

func (v *PasetoTokenVerifier) Verify(tokenString string) (*TokenClaims, error) {
	parser := paseto.NewParser()

	parser.AddRule(paseto.IssuedBy(AUTH_ISSUER))
	parser.AddRule(paseto.ForAudience(AUTH_AUDIENCE))
	parser.AddRule(paseto.Subject(AUTH_SUBJECT))

	token, err := parser.ParseV4Public(v.publicKey, tokenString, nil)
	if err != nil {
		return nil, err
	}

	aud, err := token.GetAudience()
	if err != nil {
		return nil, err
	}
	exp, err := token.GetExpiration()
	if err != nil {
		return nil, err
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	issuer, err := token.GetIssuer()
	if err != nil {
		return nil, err
	}
	sub, err := token.GetSubject()
	if err != nil {
		return nil, err
	}
	var customClaims CustomClaims
	err = token.Get("custom_claims", &customClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom claims: %w", err)
	}

	tokenClaims := NewTokenClaims(issuer, issuedAt, exp, aud, sub, customClaims)

	return &tokenClaims, nil
}

var _ TokenVerifier = (*PasetoTokenVerifier)(nil)
