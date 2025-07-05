package ginlogger

import (
	"github.com/csmart-libs/go-logger"
)

// Re-export types and functions from go-logger for convenience
type Logger = logger.Logger
type Config = logger.Config
type FileOptions = logger.FileOptions
type RotationMode = logger.RotationMode
type TimeRotationInterval = logger.TimeRotationInterval

// Re-export constants
const (
	RotationModeSize = logger.RotationModeSize
	RotationModeTime = logger.RotationModeTime
	RotationModeBoth = logger.RotationModeBoth

	RotationHourly  = logger.RotationHourly
	RotationDaily   = logger.RotationDaily
	RotationWeekly  = logger.RotationWeekly
	RotationMonthly = logger.RotationMonthly

	EnvDevelopment = logger.EnvDevelopment
	EnvStaging     = logger.EnvStaging
	EnvProduction  = logger.EnvProduction
	EnvTest        = logger.EnvTest

	LevelDebug = logger.LevelDebug
	LevelInfo  = logger.LevelInfo
	LevelWarn  = logger.LevelWarn
	LevelError = logger.LevelError
	LevelFatal = logger.LevelFatal
	LevelPanic = logger.LevelPanic

	EncodingJSON    = logger.EncodingJSON
	EncodingConsole = logger.EncodingConsole
)

// Re-export functions from go-logger
var (
	Initialize               = logger.Initialize
	NewLogger                = logger.NewLogger
	GetLogger                = logger.GetLogger
	DefaultConfig            = logger.DefaultConfig
	DefaultFileOptions       = logger.DefaultFileOptions
	DevelopmentConfig        = logger.DevelopmentConfig
	ProductionConfig         = logger.ProductionConfig
	ProductionConfigWithFile = logger.ProductionConfigWithFile
	TestConfig               = logger.TestConfig
	ConfigFromEnv            = logger.ConfigFromEnv
	GetEffectiveConfig       = logger.GetEffectiveConfig

	Debug = logger.Debug
	Info  = logger.Info
	Warn  = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Panic = logger.Panic
	With  = logger.With
	Sync  = logger.Sync

	String   = logger.String
	Int      = logger.Int
	Int64    = logger.Int64
	Uint     = logger.Uint
	Uint32   = logger.Uint32
	Uint64   = logger.Uint64
	Float64  = logger.Float64
	Bool     = logger.Bool
	Any      = logger.Any
	Err      = logger.Err
	Duration = logger.Duration
)

// Note: Gin-specific middleware and handlers are implemented in gin.go
// This includes: GinLogger, GinLoggerWithConfig, RequestIDMiddleware,
// ErrorLogger, RecoveryLogger, RequestBodyLogger, and LoggerFromContext
