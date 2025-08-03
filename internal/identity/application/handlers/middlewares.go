package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"log"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(verifier auth.TokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		claims, err := verifier.Verify(tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims == nil {
			log.Printf("AuthenticationMiddleware: Token verification returned nil claims")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
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
