// internal/api/middleware/auth_middleware.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/veezify/interview-ai-backend/internal/service"
	"net/http"
	"strings"
)

// AuthMiddleware validates JWT tokens and authorizes users
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userID", user.ID)

		c.Next()
	}
}
