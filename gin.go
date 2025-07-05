package logger

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger returns a gin.HandlerFunc for logging HTTP requests
func GinLogger() gin.HandlerFunc {
	return GinLoggerWithConfig(GinLoggerConfig{})
}

// GinLoggerConfig defines the config for GinLogger middleware
type GinLoggerConfig struct {
	Logger    Logger
	UTC       bool
	SkipPaths []string
}

// GinLoggerWithConfig returns a gin.HandlerFunc using configs
func GinLoggerWithConfig(config GinLoggerConfig) gin.HandlerFunc {
	logger := config.Logger
	if logger == nil {
		logger = GetLogger()
	}

	skipPaths := make(map[string]bool, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		// Skip logging for specified paths
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		if config.UTC {
			start = start.UTC()
		}

		// Build fields
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.Int("body_size", c.Writer.Size()),
		}

		if raw != "" {
			fields = append(fields, zap.String("query", raw))
		}

		// Add request ID if available
		if requestID := c.GetString("request_id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// Add user ID if available
		if userID := c.GetString("user_id"); userID != "" {
			fields = append(fields, zap.String("user_id", userID))
		}

		// Log based on status code
		switch {
		case c.Writer.Status() >= 500:
			logger.Error("Server error", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("Client error", fields...)
		case c.Writer.Status() >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// ErrorLogger middleware logs errors that occur during request processing
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred
		for _, err := range c.Errors {
			fields := []zap.Field{
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.Error(err.Err),
			}

			if requestID := c.GetString("request_id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			switch err.Type {
			case gin.ErrorTypeBind:
				GetLogger().Warn("Binding error", fields...)
			case gin.ErrorTypeRender:
				GetLogger().Error("Rendering error", fields...)
			case gin.ErrorTypePublic:
				GetLogger().Info("Public error", fields...)
			default:
				GetLogger().Error("Internal error", fields...)
			}
		}
	}
}

// RecoveryLogger middleware recovers from panics and logs them
func RecoveryLogger() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Any("panic", recovered),
		}

		if requestID := c.GetString("request_id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		GetLogger().Error("Panic recovered", fields...)
		c.AbortWithStatus(500)
	})
}

// RequestBodyLogger middleware logs request body (use with caution for large payloads)
type RequestBodyLoggerConfig struct {
	MaxBodySize int64
	SkipPaths   []string
}

func RequestBodyLogger(config RequestBodyLoggerConfig) gin.HandlerFunc {
	if config.MaxBodySize == 0 {
		config.MaxBodySize = 1024 * 1024 // 1MB default
	}

	skipPaths := make(map[string]bool, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		if c.Request.Body != nil && c.Request.ContentLength <= config.MaxBodySize {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// Restore the body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				fields := []zap.Field{
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("body", string(bodyBytes)),
				}

				if requestID := c.GetString("request_id"); requestID != "" {
					fields = append(fields, zap.String("request_id", requestID))
				}

				GetLogger().Debug("Request body", fields...)
			}
		}

		c.Next()
	}
}

// Helper function to generate request ID
func generateRequestID() string {
	// Simple implementation - in production, consider using UUID
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// LoggerFromContext extracts logger with request context from gin.Context
func LoggerFromContext(c *gin.Context) Logger {
	logger := GetLogger()

	// Add request context fields
	fields := []zap.Field{}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	if len(fields) > 0 {
		return logger.With(fields...)
	}

	return logger
}
