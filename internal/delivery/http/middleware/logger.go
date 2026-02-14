package middleware

import (
	"time"

	"github.com/TubagusAldiMY/go-template/internal/shared/constants"
	"github.com/TubagusAldiMY/go-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID
		requestID := c.GetHeader(constants.HeaderRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(constants.ContextKeyRequestID, requestID)
		c.Header(constants.HeaderRequestID, requestID)

		// Process request
		c.Next()

		// Log request
		duration := time.Since(start)
		logger.Info("http request",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.GetHeader(constants.HeaderUserAgent)),
		)
	}
}
