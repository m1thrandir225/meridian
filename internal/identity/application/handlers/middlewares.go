package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/pkg/auth"
)

func AuthenticationMiddleware(verifier auth.TokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrUnauthorized))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthorizationFormat))
			return
		}
		tokenString := parts[1]

		claims, err := verifier.Verify(tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidToken))
			return
		}

		if claims == nil {
			log.Printf("AuthenticationMiddleware: Token verification returned nil claims")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidToken))
			return
		}

		//Store user information in context
		customContext := auth.ContextWithTokenClaims(ctx.Request.Context(), claims)
		customContext = auth.ContextWithUserID(customContext, claims.Custom.UserID)
		customContext = auth.ContextWithEmail(customContext, claims.Custom.Email)
		ctx.Request = ctx.Request.WithContext(customContext)

		log.Printf("AuthenticationMiddleware: Successfully set context values")
		ctx.Next()

	}
}
