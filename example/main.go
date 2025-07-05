package main

import (
	"net/http"
	"time"

	logger "github.com/csmart-libs/gin-logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger with development configuration
	config := logger.DevelopmentConfig().
		WithFileOutput("logs/gin-app.log").
		WithFileRotation(10, 7, 5). // 10MB, 7 days, 5 backups
		WithFileCompression(true)

	logger.Initialize(config)

	// Create Gin router without default middleware
	r := gin.New()

	// Add comprehensive logging middleware chain
	r.Use(logger.RequestIDMiddleware())
	r.Use(logger.SecurityLogger())
	r.Use(logger.PerformanceLogger())

	// Main structured logging middleware
	r.Use(logger.StructuredLogger(logger.StructuredLoggerConfig{
		LogClientIP:    true,
		LogUserAgent:   true,
		LogReferer:     true,
		LogRequestBody: true,
		MaxBodySize:    1024 * 1024, // 1MB
		LogHeaders:     []string{"Authorization", "Content-Type", "X-API-Key"},
		SkipPaths:      []string{"/health", "/metrics"},
		CustomFields: func(c *gin.Context) []zap.Field {
			return []zap.Field{
				zap.String("service", "gin-example"),
				zap.String("version", "1.0.0"),
			}
		},
	}))

	// Error handling middleware
	r.Use(logger.ErrorLogger())
	r.Use(logger.RecoveryLogger())

	// Health check endpoint (skipped from logging)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics endpoint (skipped from logging)
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"requests": 1000, "uptime": "24h"})
	})

	// Example endpoints
	r.GET("/", func(c *gin.Context) {
		// Get logger with request context
		logger := logger.LoggerFromContext(c)
		logger.Info("Processing home request")

		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Logger Example",
			"time":    time.Now(),
		})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		logger := logger.LoggerFromContext(c)

		logger.Info("Fetching user",
			zap.String("user_id", userID),
		)

		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"name":    "John Doe",
			"email":   "john@example.com",
		})
	})

	r.POST("/users", func(c *gin.Context) {
		logger := logger.LoggerFromContext(c)

		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			logger.Warn("Invalid user data",
				zap.Error(err),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		logger.Info("Creating new user",
			zap.String("name", user.Name),
			zap.String("email", user.Email),
		)

		c.JSON(http.StatusCreated, gin.H{
			"id":    "123",
			"name":  user.Name,
			"email": user.Email,
		})
	})

	// Error endpoint for testing error logging
	r.GET("/error", func(c *gin.Context) {
		logger := logger.LoggerFromContext(c)
		logger.Error("Simulated error occurred",
			zap.String("error_type", "simulation"),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	})

	// Panic endpoint for testing recovery logging
	r.GET("/panic", func(c *gin.Context) {
		panic("This is a test panic!")
	})

	// Slow endpoint for testing performance logging
	r.GET("/slow", func(c *gin.Context) {
		logger := logger.LoggerFromContext(c)
		logger.Info("Processing slow request")

		// Simulate slow processing (> 1 second to trigger performance warning)
		time.Sleep(2 * time.Second)

		c.JSON(http.StatusOK, gin.H{"message": "Slow operation completed"})
	})

	// Suspicious endpoint for testing security logging
	r.GET("/admin", func(c *gin.Context) {
		// This will trigger security logging due to suspicious patterns in URL
		c.JSON(http.StatusOK, gin.H{"message": "Admin panel"})
	})

	logger.Info("Starting Gin server",
		zap.String("port", "8080"),
		zap.String("environment", "development"),
	)

	// Start server
	if err := r.Run(":8080"); err != nil {
		logger.Error("Failed to start server",
			zap.Error(err),
		)
	}
}
