# Gin Logger Example

This example demonstrates how to use the gin-logger package with various middleware configurations.

## Features Demonstrated

- **RequestIDMiddleware**: Automatic request ID generation and tracking
- **SecurityLogger**: Detection of suspicious request patterns
- **PerformanceLogger**: Monitoring and logging of slow requests
- **StructuredLogger**: Comprehensive request logging with custom fields
- **ErrorLogger**: Dedicated error logging
- **RecoveryLogger**: Panic recovery with detailed logging

## Running the Example

1. Navigate to the example directory:
```bash
cd gin-logger/example
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the example:
```bash
go run main.go
```

The server will start on port 8080.

## Test Endpoints

### Basic Endpoints
- `GET /` - Home page with basic logging
- `GET /users/123` - User retrieval with parameter logging
- `POST /users` - User creation with request body logging

### Special Endpoints
- `GET /health` - Health check (skipped from logging)
- `GET /metrics` - Metrics endpoint (skipped from logging)
- `GET /error` - Triggers error logging
- `GET /panic` - Triggers panic recovery logging
- `GET /slow` - Triggers performance warning (takes 2 seconds)
- `GET /admin` - May trigger security logging

## Example Requests

### Create a user:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### Test error handling:
```bash
curl http://localhost:8080/error
```

### Test panic recovery:
```bash
curl http://localhost:8080/panic
```

### Test performance monitoring:
```bash
curl http://localhost:8080/slow
```

## Log Output

The example will create logs in the `logs/gin-app.log` file with structured JSON format. You'll see different types of logs:

- Request/response logs with timing information
- Error logs with context
- Security alerts for suspicious patterns
- Performance warnings for slow requests
- Panic recovery logs

## Configuration

The example uses development configuration with:
- File output to `logs/gin-app.log`
- File rotation (10MB, 7 days, 5 backups)
- Compression enabled
- Request body logging (up to 1MB)
- Custom fields for service identification
