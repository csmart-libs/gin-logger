# Gin Logger

A comprehensive logging middleware package for the Gin web framework, built on top of the generic `go-logger` package. This package provides Gin-specific middleware and features while leveraging the powerful logging capabilities of the underlying `go-logger` foundation.

## Features

### 🚀 **Gin-Specific Middleware**

- **GinLogger**: Basic HTTP request logging middleware with status-based log levels
- **StructuredLogger**: Advanced structured logging with highly customizable fields and filtering
- **RequestIDMiddleware**: Automatic request ID generation, tracking, and header injection
- **ErrorLogger**: Dedicated middleware for logging Gin errors with context
- **RecoveryLogger**: Panic recovery with detailed logging and graceful error handling
- **RequestBodyLogger**: Request body logging with configurable size limits and path filtering
- **PerformanceLogger**: Performance monitoring with automatic slow request detection (>1s)
- **SecurityLogger**: Security event logging with attack pattern detection (SQL injection, XSS, path traversal)
- **LoggerFromContext**: Extract logger instance with request-specific context fields

### 🔧 **Advanced Configuration Options**

- **Path Filtering**: Skip specific paths or use regex patterns for flexible filtering
- **Body Logging**: Log request/response bodies with configurable size limits (default 1MB)
- **Custom Fields**: Inject custom fields via callback functions for each request
- **Time Handling**: UTC or local time support with configurable timestamps
- **Header Logging**: Selectively log specific headers (Authorization, Content-Type, etc.)
- **Client Information**: Log client IP, User Agent, Referer, and custom headers
- **Performance Tuning**: Configurable thresholds for slow request detection
- **Security Monitoring**: Built-in detection for common web attack patterns

## Installation

```bash
go get github.com/csmart-libs/gin-logger
```

## Quick Start

```go
package main

import (
    "github.com/gin-gonic/gin"
    logger "github.com/csmart-libs/gin-logger"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger with development configuration
    config := logger.DevelopmentConfig().
        WithFileOutput("logs/app.log").
        WithFileRotation(100, 30, 10) // 100MB, 30 days, 10 backups

    logger.Initialize(config)

    // Create Gin router without default middleware
    r := gin.New()

    // Add comprehensive logging middleware chain
    r.Use(logger.RequestIDMiddleware())
    r.Use(logger.SecurityLogger())
    r.Use(logger.PerformanceLogger())
    r.Use(logger.GinLogger())
    r.Use(logger.ErrorLogger())
    r.Use(logger.RecoveryLogger())

    // Your routes
    r.GET("/", func(c *gin.Context) {
        // Get logger with request context
        contextLogger := logger.LoggerFromContext(c)
        contextLogger.Info("Processing request")

        c.JSON(200, gin.H{"message": "Hello World"})
    })

    r.Run(":8080")
}
```

## Advanced Usage

### Structured Logging with Custom Configuration

```go
import (
    "regexp"
    logger "github.com/csmart-libs/gin-logger"
    "go.uber.org/zap"
)

// Advanced structured logging configuration
r.Use(logger.StructuredLogger(logger.StructuredLoggerConfig{
    Logger:          nil, // Uses global logger if nil
    LogClientIP:     true,
    LogUserAgent:    true,
    LogReferer:      true,
    LogRequestBody:  true,
    LogResponseBody: false, // Be careful with large responses
    MaxBodySize:     1024 * 1024, // 1MB limit
    LogHeaders:      []string{"Authorization", "Content-Type", "X-API-Key"},
    SkipPaths:       []string{"/health", "/metrics", "/favicon.ico"},
    SkipPathRegexps: []*regexp.Regexp{
        regexp.MustCompile(`^/static/.*`),
        regexp.MustCompile(`^/assets/.*`),
    },
    UTC: true, // Use UTC timestamps
    CustomFields: func(c *gin.Context) []zap.Field {
        return []zap.Field{
            zap.String("service", "my-api"),
            zap.String("version", "1.0.0"),
            zap.String("environment", "production"),
        }
    },
}))
```

### Performance Monitoring

```go
// Automatically log slow requests (> 1 second)
// This middleware should be placed early in the chain
r.Use(logger.PerformanceLogger())
```

### Security Monitoring

```go
// Detect and log suspicious request patterns
// Monitors for SQL injection, XSS, path traversal attempts
r.Use(logger.SecurityLogger())
```

### Request Body Logging

```go
// Log request bodies with size and path filtering
r.Use(logger.RequestBodyLogger(logger.RequestBodyLoggerConfig{
    MaxBodySize: 1024 * 1024, // 1MB limit
    SkipPaths:   []string{"/upload", "/binary", "/files"},
}))
```

### Individual Middleware Usage

```go
// Basic request logging
r.Use(logger.GinLogger())

// Or with custom configuration
r.Use(logger.GinLoggerWithConfig(logger.GinLoggerConfig{
    Logger:    nil, // Uses global logger
    UTC:       true,
    SkipPaths: []string{"/health"},
}))

// Request ID middleware (should be first)
r.Use(logger.RequestIDMiddleware())

// Error logging middleware
r.Use(logger.ErrorLogger())

// Panic recovery middleware (should be last)
r.Use(logger.RecoveryLogger())
```

## Logger Configuration

The package re-exports all configuration functions from `go-logger` for seamless integration:

```go
// Development configuration (debug level, console output)
config := logger.DevelopmentConfig()

// Production configuration (info level, JSON output)
config := logger.ProductionConfig()

// Production with file output
config := logger.ProductionConfigWithFile("logs/app.log")

// Environment-based configuration
config := logger.ConfigFromEnv()

// Custom configuration with builder pattern
config := logger.DefaultConfig().
    WithLevel("info").
    WithEnvironment("production").
    WithEncoding("json").
    WithFileOutput("logs/app.log").
    WithFileRotation(100, 30, 10). // 100MB, 30 days, 10 backups
    WithFileCompression(true).
    WithDailyRotation()

// Initialize the logger
logger.Initialize(config)
```

### Environment Variables

Configure via environment variables:

```bash
# Basic settings
LOG_LEVEL=info
LOG_ENCODING=json
LOG_OUTPUT_PATHS=stdout,file

# File settings
LOG_FILE=logs/app.log
LOG_FILE_MAX_SIZE=100
LOG_FILE_MAX_AGE=30
LOG_FILE_MAX_BACKUPS=10
LOG_FILE_COMPRESS=true
LOG_FILE_LOCAL_TIME=false

# Rotation settings
LOG_FILE_ROTATION_MODE=both
LOG_FILE_TIME_INTERVAL=daily
```

## Recommended Middleware Chain

```go
r := gin.New()

// 1. Request ID (should be first to track requests)
r.Use(logger.RequestIDMiddleware())

// 2. Security monitoring (early detection)
r.Use(logger.SecurityLogger())

// 3. Performance monitoring (timing)
r.Use(logger.PerformanceLogger())

// 4. Main structured logging
r.Use(logger.StructuredLogger(logger.StructuredLoggerConfig{
    LogClientIP:    true,
    LogUserAgent:   true,
    LogRequestBody: true,
    MaxBodySize:    1024 * 1024,
    SkipPaths:      []string{"/health", "/metrics"},
    LogHeaders:     []string{"Authorization", "Content-Type"},
}))

// 5. Error logging (captures Gin errors)
r.Use(logger.ErrorLogger())

// 6. Recovery (should be last to catch panics)
r.Use(logger.RecoveryLogger())
```

## Working with Request Context

```go
r.GET("/users/:id", func(c *gin.Context) {
    // Get logger with request context (includes request_id, user_id, etc.)
    contextLogger := logger.LoggerFromContext(c)

    userID := c.Param("id")

    // Log with automatic context fields
    contextLogger.Info("Fetching user", zap.String("user_id", userID))

    // Set user context for subsequent logs
    c.Set("user_id", userID)

    // All subsequent logs in this request will include user_id
    contextLogger.Info("User data retrieved")

    c.JSON(200, gin.H{"user_id": userID})
})
```

## Security Features

The SecurityLogger middleware automatically detects and logs:

- **SQL Injection attempts**: Patterns like `UNION SELECT`, `DROP TABLE`, etc.
- **XSS attempts**: Script tags, javascript protocols, event handlers
- **Path Traversal**: `../`, `..\\`, directory traversal patterns
- **Suspicious User Agents**: Known attack tools and scanners

```go
// Security events are logged with details
{
  "level": "warn",
  "timestamp": "2024-01-15T10:30:00Z",
  "message": "Security threat detected",
  "request_id": "req-123",
  "client_ip": "192.168.1.100",
  "user_agent": "sqlmap/1.0",
  "threat_type": "sql_injection",
  "pattern_matched": "UNION SELECT",
  "url": "/api/users?id=1' UNION SELECT * FROM users--"
}
```

## Performance Monitoring

Automatic performance monitoring logs slow requests:

```go
// Slow request log (> 1 second)
{
  "level": "warn",
  "timestamp": "2024-01-15T10:30:00Z",
  "message": "Slow request detected",
  "request_id": "req-124",
  "method": "GET",
  "path": "/api/heavy-operation",
  "duration": "2.5s",
  "status": 200
}
```

## Field Helpers

```go
// Use structured logging fields from go-logger
logger.Info("User action",
    zap.String("user_id", "123"),
    zap.String("action", "login"),
    zap.Duration("duration", time.Since(start)),
    zap.Int("attempts", 3),
)

// Or use the re-exported helpers
logger.Info("User action",
    logger.String("user_id", "123"),
    logger.String("action", "login"),
    logger.Duration("duration", time.Since(start)),
    logger.Int("attempts", 3),
)
```

## Example Application

See the complete example in `example/main.go` which demonstrates:

- Full middleware chain setup
- All logging features
- Test endpoints for each feature
- Production-ready configuration

Run the example:

```bash
cd example
go mod tidy
go run main.go
```

Test endpoints:
- `GET /` - Basic logging
- `GET /users/123` - Parameter logging
- `POST /users` - Request body logging
- `GET /error` - Error logging
- `GET /panic` - Panic recovery
- `GET /slow` - Performance monitoring
- `GET /admin` - Security monitoring

## Testing

```bash
# Run tests
go test -v

# Run tests with coverage
go test -v -cover

# Run example application
cd example && go run main.go
```

## Dependencies

- [Gin Web Framework](https://github.com/gin-gonic/gin) v1.9.1+
- [go-logger](https://github.com/csmart-libs/go-logger) - Generic logging foundation
- [Zap](https://github.com/uber-go/zap) v1.27.0+ - Fast, structured logging

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions:
- Create an issue on GitHub
- Check the example application for usage patterns
- Review the go-logger documentation for configuration options
