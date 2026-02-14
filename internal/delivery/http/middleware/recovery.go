package middleware

import (
	"fmt"
	"net/http"

	"github.com/TubagusAldiMY/go-template/pkg/logger"
	"github.com/TubagusAldiMY/go-template/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Return internal server error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
					"error":   fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
