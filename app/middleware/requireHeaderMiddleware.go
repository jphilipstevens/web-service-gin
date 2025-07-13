package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireHeader ensures a header with a specific value is present.
// The middleware responds with 400 if the header is missing or mismatched.
// @Middleware
func RequireHeader(name, value string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(name) != value {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": name + " header required"})
			return
		}
		c.Next()
	}
}
