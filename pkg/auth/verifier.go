package auth

import "aidanwoods.dev/go-paseto"

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

	parser.AddRule(paseto.IssuedBy("meridian-identity-service"))
	parser.AddRule(paseto.ForAudience("meridian-services"))

	token, err := parser.ParseV4Public(v.publicKey, tokenString, nil)
	if err != nil {
		return nil, nil
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
	jti, err := token.GetJti()
	if err != nil {
		return nil, err
	}
	sub, err := token.GetSubject()
	if err != nil {
		return nil, err
	}
	email, err := token.GetString("email")
	if err != nil {
		return nil, err
	}
	userId, err := token.GetString("user_id")
	if err != nil {
		return nil, err
	}

	tokenClaims := NewTokenClaims(issuer, aud, jti, sub, userId, email, issuedAt, exp)

	return &tokenClaims, nil
}
