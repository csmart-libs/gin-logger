package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface for dependency injection
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// ZapLogger wraps zap.Logger to implement our Logger interface
type ZapLogger struct {
	logger *zap.Logger
}

// TimeRotatingWriter wraps lumberjack.Logger to support time-based rotation
type TimeRotatingWriter struct {
	*lumberjack.Logger
	options           FileOptions
	currentTimeFormat string
	lastRotationTime  time.Time
	mu                sync.Mutex
	baseFilename      string
}

// NewTimeRotatingWriter creates a new time-based rotating writer
func NewTimeRotatingWriter(options FileOptions) *TimeRotatingWriter {
	// Extract base filename and extension
	baseFilename := options.Filename

	// Set time format based on interval
	timeFormat := options.TimeRotationFormat
	if timeFormat == "" {
		switch options.TimeRotationInterval {
		case RotationHourly:
			timeFormat = "2006-01-02-15"
		case RotationDaily:
			timeFormat = "2006-01-02"
		case RotationWeekly:
			timeFormat = "2006-W01"
		case RotationMonthly:
			timeFormat = "2006-01"
		default:
			timeFormat = "2006-01-02"
		}
	}

	// Create initial filename with timestamp
	now := time.Now()
	if options.LocalTime {
		now = now.Local()
	} else {
		now = now.UTC()
	}

	timestampedFilename := generateTimestampedFilename(baseFilename, now, timeFormat)

	lj := &lumberjack.Logger{
		Filename:   timestampedFilename,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		LocalTime:  options.LocalTime,
		Compress:   options.Compress,
	}

	return &TimeRotatingWriter{
		Logger:            lj,
		options:           options,
		currentTimeFormat: timeFormat,
		lastRotationTime:  now,
		baseFilename:      baseFilename,
	}
}

// Write implements io.Writer interface with time-based rotation check
func (w *TimeRotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if w.options.LocalTime {
		now = now.Local()
	} else {
		now = now.UTC()
	}

	// Check if we need to rotate based on time
	if w.shouldRotateByTime(now) {
		if err := w.rotateByTime(now); err != nil {
			return 0, err
		}
	}

	return w.Logger.Write(p)
}

// shouldRotateByTime checks if rotation is needed based on time interval
func (w *TimeRotatingWriter) shouldRotateByTime(now time.Time) bool {
	switch w.options.TimeRotationInterval {
	case RotationHourly:
		return now.Hour() != w.lastRotationTime.Hour() ||
			now.Day() != w.lastRotationTime.Day() ||
			now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	case RotationDaily:
		return now.Day() != w.lastRotationTime.Day() ||
			now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	case RotationWeekly:
		_, thisWeek := now.ISOWeek()
		_, lastWeek := w.lastRotationTime.ISOWeek()
		return thisWeek != lastWeek || now.Year() != w.lastRotationTime.Year()
	case RotationMonthly:
		return now.Month() != w.lastRotationTime.Month() ||
			now.Year() != w.lastRotationTime.Year()
	}
	return false
}

// rotateByTime performs time-based rotation
func (w *TimeRotatingWriter) rotateByTime(now time.Time) error {
	// Close current file
	if err := w.Logger.Close(); err != nil {
		return err
	}

	// Generate new filename with current timestamp
	newFilename := generateTimestampedFilename(w.baseFilename, now, w.currentTimeFormat)

	// Update lumberjack logger with new filename
	w.Logger.Filename = newFilename
	w.lastRotationTime = now

	return nil
}

// generateTimestampedFilename creates a filename with timestamp
func generateTimestampedFilename(baseFilename string, t time.Time, timeFormat string) string {
	dir := filepath.Dir(baseFilename)
	filename := filepath.Base(baseFilename)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	timestamp := t.Format(timeFormat)
	timestampedName := fmt.Sprintf("%s-%s%s", nameWithoutExt, timestamp, ext)

	return filepath.Join(dir, timestampedName)
}

// Global logger instance
var globalLogger Logger

// RotationMode defines how log files should be rotated
type RotationMode string

const (
	// RotationModeSize rotates based on file size only
	RotationModeSize RotationMode = "size"
	// RotationModeTime rotates based on time only
	RotationModeTime RotationMode = "time"
	// RotationModeBoth rotates based on both size and time (whichever comes first)
	RotationModeBoth RotationMode = "both"
)

// TimeRotationInterval defines the time interval for rotation
type TimeRotationInterval string

const (
	// RotationHourly rotates every hour
	RotationHourly TimeRotationInterval = "hourly"
	// RotationDaily rotates every day
	RotationDaily TimeRotationInterval = "daily"
	// RotationWeekly rotates every week
	RotationWeekly TimeRotationInterval = "weekly"
	// RotationMonthly rotates every month
	RotationMonthly TimeRotationInterval = "monthly"
)

// FileOptions holds file-specific logging options
type FileOptions struct {
	// Filename is the file to write logs to. If empty, logs will only go to stdout
	Filename string `json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int `json:"max_size" yaml:"max_size"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `json:"max_age" yaml:"max_age"`

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time. Default is UTC time.
	LocalTime bool `json:"local_time" yaml:"local_time"`

	// Compress determines if the rotated log files should be compressed using gzip
	Compress bool `json:"compress" yaml:"compress"`

	// FileMode is the file mode to use when creating log files
	FileMode os.FileMode `json:"file_mode" yaml:"file_mode"`

	// CreateDir determines if the directory should be created if it doesn't exist
	CreateDir bool `json:"create_dir" yaml:"create_dir"`

	// RotationMode determines how files should be rotated (size, time, or both)
	RotationMode RotationMode `json:"rotation_mode" yaml:"rotation_mode"`

	// TimeRotationInterval determines the time interval for rotation when using time-based rotation
	TimeRotationInterval TimeRotationInterval `json:"time_rotation_interval" yaml:"time_rotation_interval"`

	// TimeRotationFormat is the format string for time-based file naming
	// Default formats:
	// - Hourly: "2006-01-02-15"
	// - Daily: "2006-01-02"
	// - Weekly: "2006-W01"
	// - Monthly: "2006-01"
	TimeRotationFormat string `json:"time_rotation_format" yaml:"time_rotation_format"`
}

// Config holds logger configuration
type Config struct {
	Level       string      `json:"level" yaml:"level"`
	Environment string      `json:"environment" yaml:"environment"`
	OutputPaths []string    `json:"output_paths" yaml:"output_paths"`
	Encoding    string      `json:"encoding" yaml:"encoding"`
	FileOptions FileOptions `json:"file_options" yaml:"file_options"`
}

// DefaultFileOptions returns default file options
func DefaultFileOptions() FileOptions {
	return FileOptions{
		Filename:             "",               // Empty means no file output
		MaxSize:              100,              // 100 MB
		MaxAge:               30,               // 30 days
		MaxBackups:           10,               // Keep 10 backup files
		LocalTime:            false,            // Use UTC time
		Compress:             true,             // Compress old files
		FileMode:             0644,             // Standard file permissions
		CreateDir:            true,             // Create directory if not exists
		RotationMode:         RotationModeSize, // Default to size-based rotation
		TimeRotationInterval: RotationDaily,    // Default to daily rotation
		TimeRotationFormat:   "",               // Will be set based on interval
	}
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Environment: "development",
		OutputPaths: []string{"stdout"},
		Encoding:    "console",
		FileOptions: DefaultFileOptions(),
	}
}

// Initialize initializes the global logger with the given configuration
func Initialize(config Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config Config) (Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Create encoder config based on environment
	var encoderConfig zapcore.EncoderConfig
	if config.Environment == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		config.Encoding = "json"
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configure time encoding
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create encoder
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create writer syncer
	var writeSyncer zapcore.WriteSyncer

	// Check if we need file output
	if config.FileOptions.Filename != "" {
		// Create directory if needed
		if config.FileOptions.CreateDir {
			dir := filepath.Dir(config.FileOptions.Filename)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}

		var fileWriter io.Writer

		// Choose writer based on rotation mode
		switch config.FileOptions.RotationMode {
		case RotationModeTime, RotationModeBoth:
			// Use time-based rotating writer
			fileWriter = NewTimeRotatingWriter(config.FileOptions)
		default:
			// Use size-based rotating writer (lumberjack)
			fileWriter = &lumberjack.Logger{
				Filename:   config.FileOptions.Filename,
				MaxSize:    config.FileOptions.MaxSize,
				MaxAge:     config.FileOptions.MaxAge,
				MaxBackups: config.FileOptions.MaxBackups,
				LocalTime:  config.FileOptions.LocalTime,
				Compress:   config.FileOptions.Compress,
			}
		}

		// Combine stdout and file output if needed
		if len(config.OutputPaths) > 0 && config.OutputPaths[0] != "stdout" {
			// Only file output
			writeSyncer = zapcore.AddSync(fileWriter)
		} else {
			// Both stdout and file output
			writeSyncer = zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(fileWriter),
			)
		}
	} else {
		// Only stdout output
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{logger: zapLogger}, nil
}

// Implementation of Logger interface

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *ZapLogger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

// Global logger functions

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		config := DefaultConfig()
		if env := os.Getenv("APP_ENV"); env != "" {
			config.Environment = env
		}
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			config.Level = strings.ToLower(level)
		}
		Initialize(config)
	}
	return globalLogger
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	GetLogger().Panic(msg, fields...)
}

// With creates a child logger with additional fields
func With(fields ...zap.Field) Logger {
	return GetLogger().With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return GetLogger().Sync()
}

// Common field helpers

// String creates a string field
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int creates an int field
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 creates an int64 field
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint creates a uint field
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Uint32 creates a uint32 field
func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}

// Uint64 creates a uint64 field
func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

// Float64 creates a float64 field
func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

// Bool creates a bool field
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Any creates a field with any value
func Any(key string, val any) zap.Field {
	return zap.Any(key, val)
}

// Error creates an error field
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration creates a duration field
func Duration(key string, val any) zap.Field {
	return zap.Any(key, val)
}
