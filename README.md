# Logger Package

Package logger cung cấp một wrapper mạnh mẽ cho Zap logger với hỗ trợ đầy đủ cho file logging, rotation, và nhiều tính năng khác.

## Tính năng

- ✅ **Structured Logging**: Sử dụng Zap để có hiệu suất cao
- ✅ **Multiple Environments**: Development, Staging, Production, Test
- ✅ **File Output**: Ghi logs ra file với đường dẫn tùy chỉnh
- ✅ **Size-based Rotation**: Tự động rotate files dựa trên kích thước
- ✅ **Time-based Rotation**: Tự động rotate files theo thời gian (hourly, daily, weekly, monthly)
- ✅ **Combined Rotation**: Kết hợp cả size và time rotation
- ✅ **Compression**: Nén các file logs cũ để tiết kiệm dung lượng
- ✅ **Environment Configuration**: Cấu hình qua biến môi trường
- ✅ **Flexible Output**: Console, JSON, hoặc cả hai
- ✅ **Global Logger**: Sử dụng global logger hoặc tạo instance riêng

## Cài đặt

```bash
go get github.com/nguyendkn/go-libs/logger
```

## Sử dụng cơ bản

### Quick Start

```go
package main

import "github.com/nguyendkn/go-libs/logger"

func main() {
    // Sử dụng config mặc định
    logger.Initialize(logger.DefaultConfig())
    
    // Log messages
    logger.Info("Application started")
    logger.Error("Something went wrong", logger.String("error", "connection failed"))
}
```

### File Logging

```go
// Ghi logs ra file
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log")

logger.Initialize(config)
logger.Info("This will be written to file")
```

### Size-based Rotation

```go
// File logging với size rotation
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithFileRotation(100, 30, 10). // 100MB, 30 ngày, 10 backup files
    WithFileCompression(true).
    WithLocalTime(true)

logger.Initialize(config)
```

### Time-based Rotation

```go
// Hourly rotation - rotate mỗi giờ
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithHourlyRotation()

// Daily rotation - rotate mỗi ngày
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithDailyRotation()

// Weekly rotation - rotate mỗi tuần
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithWeeklyRotation()

// Monthly rotation - rotate mỗi tháng
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithMonthlyRotation()

logger.Initialize(config)
```

### Combined Size and Time Rotation

```go
// Rotate khi file đạt 50MB HOẶC mỗi giờ (cái nào đến trước)
config := logger.DefaultConfig().
    WithFileOutput("logs/app.log").
    WithBothRotation(50, 14, 7, logger.RotationHourly). // 50MB, 14 ngày, 7 backups, hourly
    WithFileCompression(true)

logger.Initialize(config)
```

### Production Setup

```go
// Production config với file output
config := logger.ProductionConfigWithFile("logs/production.log")
logger.Initialize(config)

logger.Info("Production service started",
    logger.String("version", "1.0.0"),
    logger.String("environment", "production"),
)
```

## Environment Configuration

Cấu hình qua biến môi trường:

```bash
# Basic config
export LOG_LEVEL="info"
export LOG_ENCODING="json"

# File output
export LOG_FILE="logs/app.log"
export LOG_FILE_MAX_SIZE="100"      # MB
export LOG_FILE_MAX_AGE="30"        # days
export LOG_FILE_MAX_BACKUPS="10"    # number of files
export LOG_FILE_COMPRESS="true"     # true/false
export LOG_FILE_LOCAL_TIME="false"  # true/false
export LOG_FILE_CREATE_DIR="true"   # true/false

# Time rotation options
export LOG_FILE_ROTATION_MODE="size"        # size, time, both
export LOG_FILE_TIME_INTERVAL="daily"       # hourly, daily, weekly, monthly
export LOG_FILE_TIME_FORMAT="2006-01-02"    # Custom time format (optional)
```

Sau đó:

```go
config := logger.ConfigFromEnv()
logger.Initialize(config)
```

## Preset Configurations

### Development

```go
config := logger.DevelopmentConfig()
// Level: Debug, Encoding: Console, Output: stdout
```

### Production

```go
config := logger.ProductionConfig()
// Level: Info, Encoding: JSON, Output: stdout

// Hoặc với file output
config := logger.ProductionConfigWithFile("logs/production.log")
// Level: Info, Encoding: JSON, Output: file only
```

### Test

```go
config := logger.TestConfig()
// Level: Error, Encoding: Console, Output: stdout, No file output
```

## Field Types

Package hỗ trợ nhiều loại field:

```go
logger.Info("User action",
    logger.String("user_id", "123"),
    logger.Int("age", 25),
    logger.Bool("active", true),
    logger.Float64("score", 95.5),
    logger.Any("metadata", map[string]interface{}{"key": "value"}),
    logger.Err(err),
)
```

## Custom Logger Instance

```go
// Tạo logger instance riêng
config := logger.DefaultConfig().WithFileOutput("logs/custom.log")
customLogger, err := logger.NewLogger(config)
if err != nil {
    panic(err)
}

customLogger.Info("Custom logger message")
```

## File Options

Xem [FILE_LOGGING.md](FILE_LOGGING.md) để biết chi tiết về các options cho file logging.

### FileOptions Struct

```go
type FileOptions struct {
    Filename   string      // Đường dẫn file
    MaxSize    int         // Kích thước tối đa (MB)
    MaxAge     int         // Thời gian giữ file (ngày)
    MaxBackups int         // Số lượng backup files
    LocalTime  bool        // Sử dụng local time
    Compress   bool        // Nén file cũ
    FileMode   os.FileMode // Quyền truy cập file
    CreateDir  bool        // Tự động tạo thư mục
}
```

## Examples

Xem thư mục `examples/` để có các ví dụ chi tiết:

- `file_logging_example.go`: Ví dụ về file logging với size rotation
- `time_rotation_example.go`: Ví dụ về time-based rotation (hourly, daily, weekly, monthly)

## Testing

```bash
go test -v
```

## Dependencies

- [go.uber.org/zap](https://github.com/uber-go/zap): High-performance logging
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack): Log rotation

## License

MIT License
