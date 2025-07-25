package auth

import (
	"net/http"
	"strings"
	"task_manager/domain"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSvc domain.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.SplitN(tokenHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		userID, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}