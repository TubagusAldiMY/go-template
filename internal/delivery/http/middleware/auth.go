package middleware

import (
	"strings"

	"github.com/TubagusAldiMY/go-template/internal/shared/constants"
	"github.com/TubagusAldiMY/go-template/pkg/jwt"
	"github.com/TubagusAldiMY/go-template/pkg/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.HeaderAuthorization)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user context
		c.Set(constants.ContextKeyUserID, claims.UserID)
		c.Set(constants.ContextKeyUserEmail, claims.Email)
		c.Set(constants.ContextKeyUserRole, claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString(constants.ContextKeyUserRole)
		if userRole == "" {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
