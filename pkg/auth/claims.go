package auth

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenClaims struct {
	claims map[string]interface{}
	Custom CustomClaims `json:"custom_claims"`
}

func NewTokenClaims() TokenClaims {
	return TokenClaims{}
}
