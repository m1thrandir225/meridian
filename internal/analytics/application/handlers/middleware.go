package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminHeader := c.Request.Header.Get("X-User-Is-Admin")
		if adminHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		admin, err := strconv.ParseBool(adminHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin header"})
			c.Abort()
			return
		}

		if !admin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		}

		c.Next()
	}
}
