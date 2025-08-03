package auth

import (
	"aidanwoods.dev/go-paseto"
	"time"
)

const (
	AUTH_ISSUER   = "meridian-identity-service"
	AUTH_AUDIENCE = "meridian-services"
	AUTH_SUBJECT  = "meridian-auth-token"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenClaims struct {
	Issuer         string       `json:"issuer"`
	IssuedAt       time.Time    `json:"issued_at"`
	ExpirationDate time.Time    `json:"exp"`
	Audience       string       `json:"audience"`
	Subject        string       `json:"sub"`
	Custom         CustomClaims `json:"custom_claims"`
}

func NewTokenClaims(
	issuer,
	audience,
	sub,
	userId,
	email string,
	issuedAt,
	expirationDate time.Time) TokenClaims {
	return TokenClaims{
		Issuer:         issuer,
		IssuedAt:       issuedAt,
		ExpirationDate: expirationDate,
		Audience:       audience,
		Subject:        sub,
		Custom: CustomClaims{
			UserID: userId,
			Email:  email,
		},
	}
}

func SetTokenClaims(claims TokenClaims, token *paseto.Token) (*paseto.Token, error) {
	token.SetExpiration(claims.ExpirationDate)
	token.SetIssuedAt(claims.IssuedAt)
	token.SetIssuer(claims.Issuer)
	token.SetAudience(claims.Audience)
	token.SetSubject(claims.Subject)
	if err := token.Set("custom_claims", claims.Custom); err != nil {
		return nil, err
	}
	return token, nil
}
