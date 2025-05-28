package auth

import "context"

type contextKey string

const (
	UserIDKey = contextKey("userID")
	EmailKey  = contextKey("email")
	ClaimsKey = contextKey("tokenClaims")
)

func ContextWithUserID(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIDKey, userId)
}

func ContextWithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, EmailKey, email)
}

func ContextWithTokenClaims(ctx context.Context, tokenClaims *TokenClaims) context.Context {
	return context.WithValue(ctx, ClaimsKey, tokenClaims)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(UserIDKey).(string)
	return userId, ok
}

func EmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

func TokenClaimsFromContext(ctx context.Context) (*TokenClaims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*TokenClaims)
	return claims, ok
}
