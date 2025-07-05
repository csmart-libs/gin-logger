package logger

import (
	"os"
	"strconv"
	"strings"
)

// Environment constants
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)

// Log level constants
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"
)

// Encoding constants
const (
	EncodingJSON    = "json"
	EncodingConsole = "console"
)

// ConfigFromEnv creates logger configuration from environment variables
func ConfigFromEnv() Config {
	config := DefaultConfig()

	// Get environment
	if env := os.Getenv("APP_ENV"); env != "" {
		config.Environment = strings.ToLower(env)
	}
	if env := os.Getenv("GIN_MODE"); env != "" {
		switch env {
		case "release":
			config.Environment = EnvProduction
		case "test":
			config.Environment = EnvTest
		default:
			config.Environment = EnvDevelopment
		}
	}

	// Get log level
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Level = strings.ToLower(level)
	}

	// Get encoding
	if encoding := os.Getenv("LOG_ENCODING"); encoding != "" {
		config.Encoding = strings.ToLower(encoding)
	}

	// Get output paths
	if outputs := os.Getenv("LOG_OUTPUT_PATHS"); outputs != "" {
		config.OutputPaths = strings.Split(outputs, ",")
	}

	// Get file options from environment
	if filename := os.Getenv("LOG_FILE"); filename != "" {
		config.FileOptions.Filename = filename
	}
	if maxSize := os.Getenv("LOG_FILE_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.Atoi(maxSize); err == nil {
			config.FileOptions.MaxSize = size
		}
	}
	if maxAge := os.Getenv("LOG_FILE_MAX_AGE"); maxAge != "" {
		if age, err := strconv.Atoi(maxAge); err == nil {
			config.FileOptions.MaxAge = age
		}
	}
	if maxBackups := os.Getenv("LOG_FILE_MAX_BACKUPS"); maxBackups != "" {
		if backups, err := strconv.Atoi(maxBackups); err == nil {
			config.FileOptions.MaxBackups = backups
		}
	}
	if localTime := os.Getenv("LOG_FILE_LOCAL_TIME"); localTime != "" {
		config.FileOptions.LocalTime = strings.ToLower(localTime) == "true"
	}
	if compress := os.Getenv("LOG_FILE_COMPRESS"); compress != "" {
		config.FileOptions.Compress = strings.ToLower(compress) == "true"
	}
	if createDir := os.Getenv("LOG_FILE_CREATE_DIR"); createDir != "" {
		config.FileOptions.CreateDir = strings.ToLower(createDir) == "true"
	}
	if rotationMode := os.Getenv("LOG_FILE_ROTATION_MODE"); rotationMode != "" {
		config.FileOptions.RotationMode = RotationMode(strings.ToLower(rotationMode))
	}
	if timeInterval := os.Getenv("LOG_FILE_TIME_INTERVAL"); timeInterval != "" {
		config.FileOptions.TimeRotationInterval = TimeRotationInterval(strings.ToLower(timeInterval))
	}
	if timeFormat := os.Getenv("LOG_FILE_TIME_FORMAT"); timeFormat != "" {
		config.FileOptions.TimeRotationFormat = timeFormat
	}

	// Adjust config based on environment
	switch config.Environment {
	case EnvProduction:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvStaging:
		if config.Level == "" {
			config.Level = LevelInfo
		}
		config.Encoding = EncodingJSON
	case EnvTest:
		if config.Level == "" {
			config.Level = LevelError
		}
		config.Encoding = EncodingConsole
	default: // development
		if config.Level == "" {
			config.Level = LevelDebug
		}
		config.Encoding = EncodingConsole
	}

	return config
}

// DevelopmentConfig returns configuration optimized for development
func DevelopmentConfig() Config {
	config := DefaultConfig()
	config.Level = LevelDebug
	config.Environment = EnvDevelopment
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingConsole
	return config
}

// ProductionConfig returns configuration optimized for production
func ProductionConfig() Config {
	config := DefaultConfig()
	config.Level = LevelInfo
	config.Environment = EnvProduction
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingJSON
	return config
}

// ProductionConfigWithFile returns production configuration with file output
func ProductionConfigWithFile(filename string) Config {
	config := ProductionConfig()
	config.FileOptions.Filename = filename
	config.OutputPaths = []string{"file"} // Only file output for production
	return config
}

// TestConfig returns configuration optimized for testing
func TestConfig() Config {
	config := DefaultConfig()
	config.Level = LevelError
	config.Environment = EnvTest
	config.OutputPaths = []string{"stdout"}
	config.Encoding = EncodingConsole
	config.FileOptions.Filename = "" // No file output for tests
	return config
}

// IsProduction checks if the environment is production
func (c Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsDevelopment checks if the environment is development
func (c Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

// IsTest checks if the environment is test
func (c Config) IsTest() bool {
	return c.Environment == EnvTest
}

// Validate validates the configuration
func (c Config) Validate() error {
	// Validate log level
	validLevels := map[string]bool{
		LevelDebug: true,
		LevelInfo:  true,
		LevelWarn:  true,
		LevelError: true,
		LevelFatal: true,
		LevelPanic: true,
	}
	if !validLevels[c.Level] {
		c.Level = LevelInfo
	}

	// Validate environment
	validEnvs := map[string]bool{
		EnvDevelopment: true,
		EnvStaging:     true,
		EnvProduction:  true,
		EnvTest:        true,
	}
	if !validEnvs[c.Environment] {
		c.Environment = EnvDevelopment
	}

	// Validate encoding
	validEncodings := map[string]bool{
		EncodingJSON:    true,
		EncodingConsole: true,
	}
	if !validEncodings[c.Encoding] {
		if c.IsProduction() {
			c.Encoding = EncodingJSON
		} else {
			c.Encoding = EncodingConsole
		}
	}

	// Validate output paths
	if len(c.OutputPaths) == 0 {
		c.OutputPaths = []string{"stdout"}
	}

	return nil
}

// WithLevel sets the log level
func (c Config) WithLevel(level string) Config {
	c.Level = strings.ToLower(level)
	return c
}

// WithEnvironment sets the environment
func (c Config) WithEnvironment(env string) Config {
	c.Environment = strings.ToLower(env)
	return c
}

// WithEncoding sets the encoding
func (c Config) WithEncoding(encoding string) Config {
	c.Encoding = strings.ToLower(encoding)
	return c
}

// WithOutputPaths sets the output paths
func (c Config) WithOutputPaths(paths ...string) Config {
	c.OutputPaths = paths
	return c
}

// AddOutputPath adds an output path
func (c Config) AddOutputPath(path string) Config {
	c.OutputPaths = append(c.OutputPaths, path)
	return c
}

// WithFileOutput sets file output options
func (c Config) WithFileOutput(filename string) Config {
	c.FileOptions.Filename = filename
	return c
}

// WithFileRotation sets file rotation options
func (c Config) WithFileRotation(maxSize, maxAge, maxBackups int) Config {
	c.FileOptions.MaxSize = maxSize
	c.FileOptions.MaxAge = maxAge
	c.FileOptions.MaxBackups = maxBackups
	return c
}

// WithFileCompression enables or disables file compression
func (c Config) WithFileCompression(compress bool) Config {
	c.FileOptions.Compress = compress
	return c
}

// WithLocalTime enables or disables local time for file timestamps
func (c Config) WithLocalTime(localTime bool) Config {
	c.FileOptions.LocalTime = localTime
	return c
}

// WithCreateDir enables or disables directory creation
func (c Config) WithCreateDir(createDir bool) Config {
	c.FileOptions.CreateDir = createDir
	return c
}

// WithFileMode sets the file mode for log files
func (c Config) WithFileMode(mode os.FileMode) Config {
	c.FileOptions.FileMode = mode
	return c
}

// WithRotationMode sets the rotation mode (size, time, or both)
func (c Config) WithRotationMode(mode RotationMode) Config {
	c.FileOptions.RotationMode = mode
	return c
}

// WithTimeRotation sets time-based rotation options
func (c Config) WithTimeRotation(interval TimeRotationInterval) Config {
	c.FileOptions.RotationMode = RotationModeTime
	c.FileOptions.TimeRotationInterval = interval
	return c
}

// WithTimeRotationFormat sets custom time format for rotation
func (c Config) WithTimeRotationFormat(format string) Config {
	c.FileOptions.TimeRotationFormat = format
	return c
}

// WithBothRotation enables both size and time-based rotation
func (c Config) WithBothRotation(maxSize, maxAge, maxBackups int, interval TimeRotationInterval) Config {
	c.FileOptions.RotationMode = RotationModeBoth
	c.FileOptions.MaxSize = maxSize
	c.FileOptions.MaxAge = maxAge
	c.FileOptions.MaxBackups = maxBackups
	c.FileOptions.TimeRotationInterval = interval
	return c
}

// WithHourlyRotation enables hourly rotation
func (c Config) WithHourlyRotation() Config {
	return c.WithTimeRotation(RotationHourly)
}

// WithDailyRotation enables daily rotation
func (c Config) WithDailyRotation() Config {
	return c.WithTimeRotation(RotationDaily)
}

// WithWeeklyRotation enables weekly rotation
func (c Config) WithWeeklyRotation() Config {
	return c.WithTimeRotation(RotationWeekly)
}

// WithMonthlyRotation enables monthly rotation
func (c Config) WithMonthlyRotation() Config {
	return c.WithTimeRotation(RotationMonthly)
}

// GetEffectiveConfig returns the effective configuration after applying defaults and validation
func GetEffectiveConfig() Config {
	config := ConfigFromEnv()
	config.Validate()
	return config
}
