package auth

import "time"

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenClaims struct {
	Issuer         string       `json:"issuer"`
	IssuedAt       time.Time    `json:"issued_at"`
	ExpirationDate time.Time    `json:"exp"`
	Audience       string       `json:"audience"`
	ID             string       `json:"id"`
	Subject        string       `json:"sub"`
	Custom         CustomClaims `json:"custom_claims"`
}

func NewTokenClaims(isuer, audience, id, sub, userId, email string, issuedAt, expirationDate time.Time) TokenClaims {
	return TokenClaims{
		ID:             id,
		Issuer:         isuer,
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
